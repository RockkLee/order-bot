import json
import sys
from pathlib import Path
from typing import Any

from langchain_core.output_parsers import PydanticOutputParser
from langchain_core.prompts import ChatPromptTemplate
from langchain_mistralai import ChatMistralAI
from mcp.client.session import ClientSession
from mcp.client.stdio import StdioServerParameters, stdio_client

from src.config import settings
from src.schemas import IntentResult


def _build_prompt() -> ChatPromptTemplate:
    parser = PydanticOutputParser(pydantic_object=IntentResult)
    format_instructions = parser.get_format_instructions()
    return ChatPromptTemplate.from_messages(
        [
            (
                "system",
                "You are an intent classifier for an order-bot. "
                "Return only the IntentResult JSON schema and follow it strictly. "
                "Valid intent_type values: search_menu, add_item, update_item, remove_item, "
                "show_cart, checkout, unknown.",
            ),
            (
                "human",
                "Message: {message}\nHas cart items: {has_cart_items}\n{format_instructions}",
            ),
        ]
    ).partial(format_instructions=format_instructions)


def _build_model() -> ChatMistralAI:
    return ChatMistralAI(
        model=settings.mistral_model,
        api_key=settings.mistral_api_key,
        temperature=0,
    )


class LLMIntentClient:
    async def infer_intent(self, message: str, has_cart_items: bool) -> str:
        prompt = _build_prompt()
        llm = _build_model()
        chain = prompt | llm
        result = await chain.ainvoke({"message": message, "has_cart_items": has_cart_items})
        return result.content if isinstance(result.content, str) else str(result.content)


class MCPIntentClient:
    def __init__(self) -> None:
        project_root_path = (
            # Get the absolute path of the current module's file
            # __file__ is a built-in variable that holds the path to the current script
            Path(__file__).resolve()
            .parents[2]  # Get the project root path which is two levels up from the current file path
        )
        self._server_params = StdioServerParameters(
            command=sys.executable,
            args=["-m", "src.services.intent_mcp_server"],
            cwd=project_root_path,
        )

    async def infer_intent(self, llm_response: str, has_cart_items: bool) -> IntentResult:
        async with stdio_client(self._server_params) as (read, write):
            async with ClientSession(read, write) as session:
                await session.initialize()
                result = await session.call_tool(
                    "infer_intent",
                    {"llm_response": llm_response, "has_cart_items": has_cart_items},
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

    async def parse(self, llm_response: str, has_cart_items: bool) -> IntentResult:
        text = llm_response.strip()
        if not text:
            return IntentResult(valid=False, intent_type="unknown", reason="empty_llm_response")
        try:
            return await self._client.infer_intent(text, has_cart_items)
        except Exception as exc:
            return IntentResult(
                valid=False,
                intent_type="unknown",
                reason=f"intent_parse_error:{exc.__class__.__name__}",
            )
