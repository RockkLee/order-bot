from sqlalchemy.ext.asyncio import AsyncSession
from fastapi import HTTPException

from src import repositories
from src.entities import Cart, CartItem
from src.enums import CartStatus
from src.schemas import IntentResult, ChatResponse
from src.services import cart_service
from src.grpc.client import submit_order_background


async def checkout(db: AsyncSession, session_id: str, intent: IntentResult, cart: Cart) -> ChatResponse:
    cart_summary = await cart_service.build_cart_summary(cart)
    if not intent.confirmed:
        reply = "Please confirm checkout by replying with 'confirm' or 'yes'."
        return ChatResponse(
            session_id=session_id,
            reply=reply,
            intent=intent,
            cart=cart_summary,
        )

    cart = await cart_service.lock_cart(db, session_id)
    if cart.status != CartStatus.OPEN:
        raise HTTPException(status_code=400, detail="Cart is closed")
    items: list[CartItem] = await cart.awaitable_attrs.items
    if not items:
        raise HTTPException(status_code=400, detail="Cart is empty")

    total_scaled = sum(item.total_price_scaled for item in items)
    order = await repositories.insert_order(db, cart, total_scaled)
    await repositories.insert_order_items(db, order, items)
    order_items = list(await order.awaitable_attrs.order_items)
    order_id = order.id

    cart.status = CartStatus.CLOSED
    await db.commit()

    submit_order_background(order, order_items)

    cart_summary = await cart_service.build_cart_summary(cart)
    return ChatResponse(
        session_id=session_id,
        reply=f"Order placed! Your order id is {order_id}.",
        intent=intent,
        cart=cart_summary,
        order_id=order_id,
    )
