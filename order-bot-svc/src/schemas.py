from typing import Literal
from pydantic import BaseModel, Field

from src.enums import CartStatus


class ChatRequest(BaseModel):
    menu_id: str = Field(..., alias="menuId")
    bot_id: str = Field(..., alias="botId")
    message: str = Field(..., min_length=1)


class MenuItemOut(BaseModel):
    name: str
    price: float


class CartItemOut(BaseModel):
    menu_item_id: str
    name: str
    quantity: int
    unit_price_scaled: int
    total_price_scaled: int
    unit_price: float
    total_price: float


class CartSummary(BaseModel):
    session_id: str
    status: CartStatus
    items: list[CartItemOut] = Field(default_factory=list)
    total_price_scaled: int
    total_price: float = 0


class IntentItem(BaseModel):
    menu_item_id: str
    quantity: int


class IntentResult(BaseModel):
    valid: bool
    intent_type: Literal[
        "search_menu",
        "mutate_cart_items",
        "show_cart",
        "checkout",
        "unknown",
    ]
    items: list[IntentItem] = Field(default_factory=list)
    query: str | None = None
    confirmed: bool = False
    reason: str | None = None


class ChatResponse(BaseModel):
    session_id: str
    reply: str
    intent: IntentResult
    cart: CartSummary
    order_id: str | None = None
    menu_results: list[MenuItemOut] = Field(default_factory=list)
