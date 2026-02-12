import json

from mcp.server.fastmcp import FastMCP
from pydantic import BaseModel

from src.schemas import CartItemIntent


MCP_SERVER_NAME = "order-bot-intent"


mcp = FastMCP(MCP_SERVER_NAME)


def _json_payload(payload: dict) -> str:
    # Add this because the field in the json payload could be a pydantic object (e.g., items: list[IntentItem])
    def _default(value: object) -> object:
        if isinstance(value, BaseModel):
            return value.model_dump()
        raise TypeError(f"Unsupported type: {type(value)!r}")

    return json.dumps(payload, default=_default)


@mcp.tool()
def search_menu() -> str:
    """Search menu"""
    return _json_payload({"intent_type": "search_menu"})


@mcp.tool()
def mutate_cart_items(items: list[CartItemIntent]) -> str:
    """Update menu items in the cart by menu item id and quantity and return the whole cart items including the ones in the cart"""
    return _json_payload({"intent_type": "mutate_cart_items", "items": items})


@mcp.tool()
def show_cart() -> str:
    """Show the current cart contents."""
    return _json_payload({"intent_type": "show_cart"})


@mcp.tool()
def checkout(confirmed: bool = False) -> str:
    """Checkout the cart; set confirmed once the user explicitly agrees."""
    return _json_payload({"intent_type": "checkout", "confirmed": confirmed})


@mcp.tool()
def unknown(reason: str = "unknown") -> str:
    """If you can't choose any other MCP tool, choose this one"""
    return _json_payload({"intent_type": "unknown", "reason": reason})


def main() -> None:
    mcp.run()


if __name__ == "__main__":
    main()
