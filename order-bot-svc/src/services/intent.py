import json
import sys
from typing import Any

from mcp.client.session import ClientSession
from mcp.client.stdio import StdioServerParameters, stdio_client

from src.schemas import IntentResult


class MCPIntentClient:
    def __init__(self) -> None:
        self._server_params = StdioServerParameters(
            command=sys.executable,
            args=["-m", "src.services.intent_mcp_server"],
        )

    async def infer_intent(self, message: str, has_cart_items: bool) -> IntentResult:
        async with stdio_client(self._server_params) as (read, write):
            async with ClientSession(read, write) as session:
                await session.initialize()
                result = await session.call_tool(
                    "infer_intent",
                    {"message": message, "has_cart_items": has_cart_items},
                )
        payload = self._parse_result_payload(result.content)
        return IntentResult.model_validate(payload)

    def _parse_result_payload(self, payload: Any) -> dict:
        if isinstance(payload, dict):
            return payload
        if isinstance(payload, str):
            return json.loads(payload)
        if isinstance(payload, list) and payload:
            first = payload[0]
            if isinstance(first, dict):
                return first
            if hasattr(first, "text"):
                return json.loads(first.text)
        return {}


class IntentParser:
    def __init__(self) -> None:
        self._client = MCPIntentClient()

    async def parse(self, message: str, has_cart_items: bool) -> IntentResult:
        text = message.strip()
        if not text:
            return IntentResult(valid=False, intent_type="unknown", reason="empty")
        try:
            return await self._client.infer_intent(text, has_cart_items)
        except Exception as exc:
            return IntentResult(
                valid=False,
                intent_type="unknown",
                reason=f"mcp_error:{exc.__class__.__name__}",
            )
