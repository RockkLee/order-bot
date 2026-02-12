import json
import logging
from pathlib import Path
from typing import Any

from langchain.agents import create_agent
from langchain_core.messages import ToolMessage
from langchain_core.output_parsers import PydanticOutputParser
from langchain_core.prompts import ChatPromptTemplate
from langchain_mcp_adapters.client import MultiServerMCPClient
from langchain_mistralai import ChatMistralAI

from src.config import settings
from src.schemas import IntentResult, MenuItemIntent, CartItemIntent


class MCPIntentClient:
    def __init__(self) -> None:
        project_root_path = str(Path(__file__).resolve().parents[2])
        self._client = MultiServerMCPClient(
            {
                "order_bot_intent": {
                    "transport": "stdio",
                    "command": "python",
                    "args": ["-m", "src.intent.mcp_server"],
                    "cwd": project_root_path,
                }
            }
        )

    async def get_tools(self) -> list[Any]:
        return await self._client.get_tools()


class IntentParser:
    def __init__(self) -> None:
        self._client = MCPIntentClient()
        self._llm = ChatMistralAI(
            model=settings.mistral_model,
            api_key=settings.mistral_api_key,
            temperature=0,
        )
        self._agent = None

    async def parse(
        self, message: str, menu_item_intents: list[MenuItemIntent], has_cart_items: bool, cart_item_intents: list[CartItemIntent]
    ) -> IntentResult:
        text = message.strip()
        if not text:
            return IntentResult(valid=False, intent_type="unknown", reason="empty")

        try:
            if self._agent is None:
                tools = await self._client.get_tools()
                for tool in tools:
                    # Ensure tool output is a plain string for the Mistral tool message schema.
                    tool.response_format = "content"
                self._agent = create_agent(
                    model=self._llm,
                    tools=tools,
                    system_prompt=self._system_prompt(),
                )

            response = await self._agent.ainvoke(
                {
                    "messages": [
                        {
                            "role": "user",
                            "content": self._build_user_prompt(text, menu_item_intents, has_cart_items, cart_item_intents),
                        }
                    ]
                }
            )

            return self._parse_agent_response(response)
        except Exception as exc:
            logging.exception("mcp error")
            return IntentResult(
                valid=False,
                intent_type="unknown",
                reason=f"mcp_error:{exc.__class__.__name__}",
            )

    def _system_prompt(self) -> str:
        parser = PydanticOutputParser(pydantic_object=IntentResult)
        return (
            "You are an intent classifier for an order-bot. "
            "You must call exactly one MCP tool that best represents the user request. "
            "Then return only valid JSON that follows this schema:\n"
            f"{parser.get_format_instructions()}"
        )

    def _build_user_prompt(
        self, message: str, menu_item_intents: list[MenuItemIntent], has_cart_items: bool, cart_item_intents: list[CartItemIntent]
    ) -> str:
        menu_payload = json.dumps([m.model_dump() for m in menu_item_intents], ensure_ascii=False)
        cart_payload = json.dumps([m.model_dump() for m in cart_item_intents], ensure_ascii=False)
        prompt = ChatPromptTemplate.from_template(
            "Message: {message}\n"
            "Has cart items: {has_cart_items}\n"
            "Cart: {cart}\n"
            "Menu: {menu}\n"
            "Steps:\n"
            "1. Select exactly one intent type from: search_menu, mutate_cart_items, show_cart, checkout, unknown.\n"
            "2. Use intent args derived from the message.\n"
            "3. Return the final response as IntentResult JSON only."
        )
        return prompt.format(message=message, has_cart_items=has_cart_items, menu=menu_payload, cart=cart_payload)

    def _parse_agent_response(self, response: dict[str, Any]) -> IntentResult:
        if isinstance(response, str):
            # Some agent configs return the tool payload directly as a string.
            payload = self._parse_text_payload(response) or {}
            if "valid" not in payload:
                payload["valid"] = payload.get("intent_type") != "unknown"
            return IntentResult.model_validate(payload)

        if isinstance(response, dict) and isinstance(response.get("output"), str):
            # Some agent configs return the final string payload under "output".
            payload = self._parse_text_payload(response["output"]) or {}
            if "valid" not in payload:
                payload["valid"] = payload.get("intent_type") != "unknown"
            return IntentResult.model_validate(payload)

        messages = response.get("messages", []) if isinstance(response, dict) else []
        tool_payload: dict[str, Any] | None = None
        final_payload: dict[str, Any] | None = None

        for msg in messages:
            if isinstance(msg, ToolMessage):
                tool_payload = self._parse_text_payload(getattr(msg, "content", ""))

            content = getattr(msg, "content", None)
            if isinstance(content, str):
                parsed = self._parse_text_payload(content)
                if parsed:
                    final_payload = parsed

        payload = final_payload or tool_payload or {}
        if "valid" not in payload:
            payload["valid"] = payload.get("intent_type") != "unknown"

        return IntentResult.model_validate(payload)

    def _parse_text_payload(self, payload: Any) -> dict[str, Any] | None:
        if isinstance(payload, dict):
            return payload
        if isinstance(payload, list):
            for part in payload:
                if isinstance(part, dict) and part.get("type") == "text":
                    text = part.get("text", "")
                    parsed = self._safe_json_loads(text)
                    if parsed:
                        return parsed
        if isinstance(payload, str):
            return self._safe_json_loads(payload)
        return None

    def _safe_json_loads(self, text: str) -> dict[str, Any] | None:
        text = text.strip()
        if not text:
            return None
        if text.startswith("```"):
            text = text.strip("`")
            if text.startswith("json"):
                text = text[4:].strip()
        try:
            loaded = json.loads(text)
            if isinstance(loaded, dict):
                return loaded
        except json.JSONDecodeError:
            return None
        return None
