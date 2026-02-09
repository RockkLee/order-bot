from datetime import datetime
from sqlalchemy.ext.asyncio import AsyncSession
from src import repositories
from src.entities import Cart
from src.enums import CartStatus
from src.schemas import CartSummary, CartItemOut, IntentResult, ChatResponse
from fastapi import HTTPException
from src.services import response_builder
from src.utils import money_util


def build_cart_summary(cart: Cart) -> CartSummary:
    items = [
        CartItemOut(
            menu_item_id=item.menu_item_id,
            name=item.name,
            quantity=item.quantity,
            unit_price_scaled=item.unit_price_scaled,
            total_price_scaled=item.total_price_scaled,
            unit_price=money_util.to_float(item.unit_price_scaled),
            total_price=money_util.to_float(item.total_price_scaled)
        )
        for item in cart.items
    ]
    scaled_total = sum(item.total_price_scaled for item in items)
    total = sum(item.total_price for item in items)
    return CartSummary(
        session_id=cart.session_id,
        status=cart.status,
        items=items,
        total_price_scaled=scaled_total,
        total_price=total,
    )


async def ensure_cart(db: AsyncSession, session_id: str) -> Cart:
    cart = await repositories.get_cart_by_session(db, session_id)
    if not cart:
        cart = await repositories.create_cart(db, session_id)
    return cart


async def lock_cart(db: AsyncSession, session_id: str) -> Cart:
    cart = await repositories.get_cart_by_session(db, session_id, for_update=True)
    if not cart:
        cart = await repositories.create_cart(db, session_id)
    return cart


def touch_cart(cart: Cart) -> None:
    cart.updated_at = datetime.utcnow()


async def mutate_cart(db: AsyncSession, session_id: str, intent: IntentResult) -> ChatResponse:
    async with db.begin():
        cart = await lock_cart(db, session_id)
        if cart.status != CartStatus.OPEN:
            raise HTTPException(status_code=400, detail="Cart is closed")

        for item in intent.items:
            menu_item = await repositories.get_menu_item_by_menu_item_id(db, item.sku)
            if intent.intent_type == "remove_item":
                await repositories.remove_cart_item(db, cart, item.sku)
                continue

            if item.quantity <= 0:
                raise HTTPException(status_code=400, detail="Quantity must be positive")

            await repositories.upsert_cart_item(
                db,
                cart,
                menu_item_id=menu_item.id,
                name=menu_item.name,
                quantity=item.quantity,
                unit_price_cents=menu_item.price,
            )

        touch_cart(cart)

    cart_summary = build_cart_summary(cart)
    return ChatResponse(
        session_id=session_id,
        reply=response_builder.build_reply(intent, cart_summary),
        intent=intent,
        cart=cart_summary,
    )
