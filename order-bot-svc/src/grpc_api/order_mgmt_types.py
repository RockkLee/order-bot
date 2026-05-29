from dataclasses import dataclass


@dataclass
class CheckoutOrderItem:
    menu_item_id: str
    name: str
    quantity: int
    unit_price_scaled: int
    total_price_scaled: int


@dataclass
class CheckoutOrderRequest:
    order_id: str
    bot_id: str
    cart_id: str
    session_id: str
    total_scaled: int
    items: list[CheckoutOrderItem]
