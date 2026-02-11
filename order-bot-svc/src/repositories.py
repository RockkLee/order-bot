from sqlalchemy import select, delete
from sqlalchemy.ext.asyncio import AsyncSession
from src.entities import MenuItem, Cart, CartItem, Order, OrderItem


async def get_menu_by_query(db: AsyncSession, menu_id: str) -> list[MenuItem]:
    stmt = select(MenuItem).where(
        MenuItem.menu_id == menu_id
    )
    result = await db.scalars(stmt)
    return list(result.all())


async def get_menu_item_by_menu_item_ids(db: AsyncSession, menu_item_ids: list[str]) -> list[MenuItem]:
    stmt = select(MenuItem).where(MenuItem.id.in_(menu_item_ids))
    result = await db.scalars(stmt)
    return list(result.all())


async def get_cart_by_session(
    db: AsyncSession, session_id: str, for_update: bool = False
) -> Cart | None:
    stmt = select(Cart).where(Cart.session_id == session_id)
    if for_update:
        stmt = stmt.with_for_update()
    result = await db.scalars(stmt)
    return result.first()


async def create_cart(db: AsyncSession, session_id: str) -> Cart:
    cart = Cart(session_id=session_id)
    db.add(cart)
    await db.flush()
    return cart


async def upsert_cart_item(
    db: AsyncSession,
    cart: Cart,
    cart_items: list[CartItem]
):
    delete(CartItem).where(CartItem.cart_id == cart.id)
    db.add_all(cart_items)
    await db.commit()


async def insert_order(db: AsyncSession, cart: Cart, total_cents: int) -> Order:
    order = Order(cart_id=cart.id, session_id=cart.session_id, total_cents=total_cents)
    db.add(order)
    await db.flush()
    return order


async def insert_order_items(
    db: AsyncSession, order: Order, cart_items: list[CartItem]
) -> None:
    for item in cart_items:
        order_item = OrderItem(
            order_id=order.id,
            menu_item_id=item.menu_item_id,
            name=item.name,
            quantity=item.quantity,
            unit_price_cents=item.unit_price_cents,
            line_total_cents=item.line_total_cents,
        )
        db.add(order_item)
