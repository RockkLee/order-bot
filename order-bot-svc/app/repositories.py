from sqlalchemy import select
from sqlalchemy.orm import Session
from app.models import MenuItem, Cart, CartItem, Order, OrderItem


def get_menu_by_query(db: Session, query: str) -> list[MenuItem]:
    stmt = select(MenuItem).where(MenuItem.is_available.is_(True))
    if query:
        like_query = f"%{query.lower()}%"
        stmt = stmt.where(
            (MenuItem.name.ilike(like_query)) | (MenuItem.description.ilike(like_query))
        )
    return list(db.scalars(stmt).all())


def get_menu_item_by_sku(db: Session, sku: str) -> MenuItem | None:
    stmt = select(MenuItem).where(MenuItem.sku == sku)
    return db.scalars(stmt).first()


def get_cart_by_session(db: Session, session_id: str, for_update: bool = False) -> Cart | None:
    stmt = select(Cart).where(Cart.session_id == session_id)
    if for_update:
        stmt = stmt.with_for_update()
    return db.scalars(stmt).first()


def create_cart(db: Session, session_id: str) -> Cart:
    cart = Cart(session_id=session_id)
    db.add(cart)
    db.flush()
    return cart


def upsert_cart_item(
    db: Session,
    cart: Cart,
    sku: str,
    name: str,
    quantity: int,
    unit_price_cents: int,
) -> CartItem:
    stmt = select(CartItem).where(CartItem.cart_id == cart.id, CartItem.sku == sku)
    existing = db.scalars(stmt).first()
    line_total_cents = quantity * unit_price_cents
    if existing:
        existing.quantity = quantity
        existing.unit_price_cents = unit_price_cents
        existing.line_total_cents = line_total_cents
        return existing

    item = CartItem(
        cart_id=cart.id,
        sku=sku,
        name=name,
        quantity=quantity,
        unit_price_cents=unit_price_cents,
        line_total_cents=line_total_cents,
    )
    db.add(item)
    return item


def remove_cart_item(db: Session, cart: Cart, sku: str) -> bool:
    stmt = select(CartItem).where(CartItem.cart_id == cart.id, CartItem.sku == sku)
    existing = db.scalars(stmt).first()
    if not existing:
        return False
    db.delete(existing)
    return True


def insert_order(db: Session, cart: Cart, total_cents: int) -> Order:
    order = Order(cart_id=cart.id, session_id=cart.session_id, total_cents=total_cents)
    db.add(order)
    db.flush()
    return order


def insert_order_items(db: Session, order: Order, cart_items: list[CartItem]) -> None:
    for item in cart_items:
        order_item = OrderItem(
            order_id=order.id,
            sku=item.sku,
            name=item.name,
            quantity=item.quantity,
            unit_price_cents=item.unit_price_cents,
            line_total_cents=item.line_total_cents,
        )
        db.add(order_item)
