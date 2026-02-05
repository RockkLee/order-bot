import json

from mcp.server.fastmcp import FastMCP

from src.schemas import IntentResult


MCP_SERVER_NAME = "order-bot-intent"


mcp = FastMCP(MCP_SERVER_NAME)


@mcp.tool()
async def infer_intent(llm_response: str, has_cart_items: bool) -> dict:
    _ = has_cart_items
    payload = json.loads(llm_response)
    return IntentResult.model_validate(payload).model_dump()


def main() -> None:
    mcp.run()


if __name__ == "__main__":
    main()
