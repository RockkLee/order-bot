import json
import os

from langchain_core.output_parsers import PydanticOutputParser
from langchain_core.prompts import ChatPromptTemplate
from langchain_mistralai import ChatMistralAI
from mcp.server.fastmcp import FastMCP

from src.schemas import IntentResult


MCP_SERVER_NAME = "order-bot-intent"


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
        model=os.environ.get("MISTRAL_MODEL", "mistral-large-latest"),
        temperature=0,
    )


mcp = FastMCP(MCP_SERVER_NAME)


@mcp.tool()
async def infer_intent(message: str, has_cart_items: bool) -> dict:
    parser = PydanticOutputParser(pydantic_object=IntentResult)
    prompt = _build_prompt()
    llm = _build_model()
    chain = prompt | llm | parser
    result = await chain.ainvoke({"message": message, "has_cart_items": has_cart_items})
    return json.loads(result.model_dump_json())


def main() -> None:
    mcp.run()


if __name__ == "__main__":
    main()
