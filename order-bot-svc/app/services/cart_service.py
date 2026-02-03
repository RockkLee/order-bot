from datetime import datetime
from sqlalchemy.ext.asyncio import AsyncSession
from app import repositories
from app.models import Cart
from app.schemas import CartSummary, CartItemOut


def build_cart_summary(cart: Cart) -> CartSummary:
    items = [
        CartItemOut(
            sku=item.sku,
            name=item.name,
            quantity=item.quantity,
            unit_price_cents=item.unit_price_cents,
            line_total_cents=item.line_total_cents,
        )
        for item in cart.items
    ]
    total_cents = sum(item.line_total_cents for item in cart.items)
    return CartSummary(
        session_id=cart.session_id,
        status=cart.status,
        items=items,
        total_cents=total_cents,
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
