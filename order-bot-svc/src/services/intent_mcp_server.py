from mcp.server.fastmcp import FastMCP


MCP_SERVER_NAME = "order-bot-intent"


mcp = FastMCP(MCP_SERVER_NAME)


@mcp.tool()
def search_menu(query: str) -> dict:
    return {"intent_type": "search_menu", "query": query}


@mcp.tool()
def add_item(items: list[dict]) -> dict:
    return {"intent_type": "add_item", "items": items}


@mcp.tool()
def update_item(items: list[dict]) -> dict:
    return {"intent_type": "update_item", "items": items}


@mcp.tool()
def remove_item(items: list[dict]) -> dict:
    return {"intent_type": "remove_item", "items": items}


@mcp.tool()
def show_cart() -> dict:
    return {"intent_type": "show_cart"}


@mcp.tool()
def checkout(confirmed: bool = False) -> dict:
    return {"intent_type": "checkout", "confirmed": confirmed}


@mcp.tool()
def unknown(reason: str = "unknown") -> dict:
    return {"intent_type": "unknown", "reason": reason}


def main() -> None:
    mcp.run()


if __name__ == "__main__":
    main()
