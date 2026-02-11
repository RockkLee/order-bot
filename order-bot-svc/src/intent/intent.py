import json
import logging
from pathlib import Path
from typing import Any

from httpx import HTTPStatusError
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
        self._tools_by_name: dict[str, Any] | None = None

    async def parse(self, message: str, has_cart_items: bool) -> IntentResult:
        text = message.strip()
        if not text:
            return IntentResult(valid=False, intent_type="unknown", reason="empty")

        try:
            if self._tools_by_name is None:
                tools = await self._client.get_tools()
                self._tools_by_name = {tool.name: tool for tool in tools}

            model_result = await self._classify_intent(text, has_cart_items)
            tool_result = await self._call_tool_for_intent(model_result)
            if tool_result is not None:
                return tool_result
            return model_result
        except HTTPStatusError as exc:
            response = exc.response
            logging.error("mcp_error:%s:%s", response.status_code, response.text)
            logging.exception("mcp http error")
            return IntentResult(
                valid=False,
                intent_type="unknown",
                reason=f"mcp_error:{response.status_code}:{response.text}",
            )
        except Exception as exc:
            logging.exception("mcp error")
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
            "Then return only valid JSON that follows this schema:\n"
            f"{parser.get_format_instructions()}"
        )

    def _build_user_prompt(self, message: str, has_cart_items: bool) -> str:
        prompt = ChatPromptTemplate.from_template(
            "Message: {message}\nHas cart items: {has_cart_items}\n"
            "Steps:\n"
            "1. Select exactly one intent type from: search_menu, mutate_cart_items, show_cart, checkout, unknown.\n"
            "2. Use intent args derived from the message.\n"
            "3. Return the final response as IntentResult JSON only."
        )
        return prompt.format(message=message, has_cart_items=has_cart_items)

    async def _classify_intent(self, message: str, has_cart_items: bool) -> IntentResult:
        parser = PydanticOutputParser(pydantic_object=IntentResult)
        prompt = ChatPromptTemplate.from_messages(
            [
                ("system", "{system}"),
                ("user", "{user}"),
            ]
        )
        chain = prompt | self._llm | parser
        return await chain.ainvoke(
            {
                "system": self._system_prompt(),
                "user": self._build_user_prompt(message, has_cart_items),
            }
        )

    async def _call_tool_for_intent(self, result: IntentResult) -> IntentResult | None:
        if self._tools_by_name is None:
            return None

        tool = self._tools_by_name.get(result.intent_type)
        if tool is None:
            return None

        args = self._tool_args_for_intent(result)
        tool_response = await tool.ainvoke(args)
        payload = self._parse_text_payload(tool_response) or {}
        if result.query and "query" not in payload:
            payload["query"] = result.query
        if result.reason and "reason" not in payload:
            payload["reason"] = result.reason
        if result.items and "items" not in payload:
            payload["items"] = [item.model_dump() for item in result.items]
        if "confirmed" not in payload:
            payload["confirmed"] = result.confirmed
        if "valid" not in payload:
            payload["valid"] = payload.get("intent_type") != "unknown"

        return IntentResult.model_validate(payload)

    def _tool_args_for_intent(self, result: IntentResult) -> dict[str, Any]:
        if result.intent_type == "mutate_cart_items":
            return {"items": [item.model_dump() for item in result.items]}
        if result.intent_type == "checkout":
            return {"confirmed": result.confirmed}
        if result.intent_type == "unknown":
            return {"reason": result.reason or "unknown"}
        return {}


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
