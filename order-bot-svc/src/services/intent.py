import json
from pathlib import Path
from typing import Any

from langchain.agents import create_agent
from langchain_core.messages import ToolMessage
from langchain_core.output_parsers import PydanticOutputParser
from langchain_core.prompts import ChatPromptTemplate
from langchain_mcp_adapters.client import MultiServerMCPClient
from langchain_mistralai import ChatMistralAI

from src.config import settings
from src.schemas import IntentResult


class MCPIntentClient:
    def __init__(self) -> None:
        project_root_path = str(Path(__file__).resolve().parents[2])
        self._client = MultiServerMCPClient(
            {
                "order_bot_intent": {
                    "transport": "stdio",
                    "command": "python",
                    "args": ["-m", "src.services.intent_mcp_server"],
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

    async def parse(self, message: str, has_cart_items: bool) -> IntentResult:
        text = message.strip()
        if not text:
            return IntentResult(valid=False, intent_type="unknown", reason="empty")

        try:
            if self._agent is None:
                tools = await self._client.get_tools()
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
                            "content": self._build_user_prompt(text, has_cart_items),
                        }
                    ]
                }
            )

            return self._parse_agent_response(response)
        except Exception as exc:
            fallback = self._fallback_parse(text)
            if fallback is not None:
                return fallback
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

    def _build_user_prompt(self, message: str, has_cart_items: bool) -> str:
        prompt = ChatPromptTemplate.from_template(
            "Message: {message}\nHas cart items: {has_cart_items}\n"
            "Steps:\n"
            "1. Select and call exactly one tool from: search_menu, add_item, update_item, remove_item, show_cart, checkout, unknown.\n"
            "2. Use tool args derived from the message.\n"
            "3. Return the final response as IntentResult JSON only."
        )
        return prompt.format(message=message, has_cart_items=has_cart_items)


    def _fallback_parse(self, text: str) -> IntentResult | None:
        lowered = text.lower()
        if "show cart" in lowered:
            return IntentResult(valid=True, intent_type="show_cart")
        if lowered.startswith("checkout"):
            return IntentResult(valid=True, intent_type="checkout", confirmed=False)
        if lowered in {"confirm", "yes"}:
            return IntentResult(valid=True, intent_type="checkout", confirmed=True)
        if "add" in lowered and "sku" in lowered:
            parts = lowered.split()
            quantity = 1
            sku = None
            for idx, token in enumerate(parts):
                if token.isdigit():
                    quantity = int(token)
                if token == "sku" and idx + 1 < len(parts):
                    sku = parts[idx + 1]
            if sku:
                return IntentResult(
                    valid=True,
                    intent_type="add_item",
                    items=[{"sku": sku, "quantity": quantity}],
                )
        return None

    def _parse_agent_response(self, response: dict[str, Any]) -> IntentResult:
        messages = response.get("messages", [])
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
