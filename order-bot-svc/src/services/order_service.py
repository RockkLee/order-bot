from sqlalchemy.ext.asyncio import AsyncSession
from fastapi import HTTPException

from src import repositories
from src.entities import Cart, CartItem
from src.enums import CartStatus, OrderStatus
from src.schemas import IntentResult, ChatResponse
from src.services import cart_service
from src.services.order_mgmt_grpc_service import (
    OrderMgmtGrpcService,
)
from src.grpc_api.order_mgmt_types import (
    CheckoutOrderItem,
    CheckoutOrderRequest,
)


order_mgmt_grpc_service = OrderMgmtGrpcService()


async def checkout(db: AsyncSession, session_id: str, bot_id: str, intent: IntentResult, cart: Cart) -> ChatResponse:
    cart_summary = await cart_service.build_cart_summary(cart)
    if not intent.confirmed:
        reply = "Please confirm checkout by replying with 'confirm' or 'yes'."
        return ChatResponse(
            session_id=session_id,
            reply=reply,
            intent=intent,
            cart=cart_summary,
        )

    # A tx will automatically start once the db session is created in SQLAlchemy 2.0,
    # so we don't have to manually create a tx
    # async with db.begin():
    cart = await cart_service.lock_cart(db, session_id)
    if cart.status != CartStatus.OPEN:
        raise HTTPException(status_code=400, detail="Cart is closed")
    items: list[CartItem] = await cart.awaitable_attrs.items
    if not items:
        raise HTTPException(status_code=400, detail="Cart is empty")

    total_scaled = sum(item.total_price_scaled for item in items)
    order = await repositories.insert_order(db, cart, bot_id, total_scaled)
    await repositories.insert_order_items(db, order, items)

    grpc_request = CheckoutOrderRequest(
        order_id=order.id,
        bot_id=bot_id,
        cart_id=cart.id,
        session_id=cart.session_id,
        total_scaled=total_scaled,
        items=[
            CheckoutOrderItem(
                menu_item_id=item.menu_item_id,
                name=item.name,
                quantity=item.quantity,
                unit_price_scaled=item.unit_price_scaled,
                total_price_scaled=item.total_price_scaled,
            )
            for item in items
        ],
    )
    await order_mgmt_grpc_service.create_order_with_items(grpc_request)
    await order_mgmt_grpc_service.notify_order_completed(order.id)
    await repositories.update_order_status(db, order, OrderStatus.COMPLETED)
    order_id = order.id

    cart.status = CartStatus.CLOSED
    await db.commit()

    cart_summary = await cart_service.build_cart_summary(cart)
    return ChatResponse(
        session_id=session_id,
        reply=f"Order placed! Your order id is {order_id}.",
        intent=intent,
        cart=cart_summary,
        order_id=order_id,
    )
