import uuid
from sqlalchemy.ext.asyncio import AsyncSession
from src import repositories
from src.entities import Cart, MenuItem, CartItem
from src.enums import CartStatus
from src.schemas import CartSummary, CartItemOut, IntentResult, ChatResponse, CartItemIntent
from fastapi import HTTPException
from src.services import response_builder
from src.utils import money_util


async def build_cart_summary(cart: Cart) -> CartSummary:
    items_src = await cart.awaitable_attrs.items
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
        for item in items_src
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


async def build_cart_item_intents(cart: Cart) -> list[CartItemIntent]:
    items_src = await cart.awaitable_attrs.items
    return [
        CartItemIntent(
            menu_item_id=item.menu_item_id,
            quantity=item.quantity,
        )
        for item in items_src
    ]


async def get_cart(db: AsyncSession, session_id: str) -> Cart:
    cart = await repositories.get_cart_by_session(db, session_id)
    if not cart:
        cart = await repositories.create_cart(db, session_id)
    return cart


async def lock_cart(db: AsyncSession, session_id: str) -> Cart:
    cart = await repositories.get_cart_by_session(db, session_id, for_update=True)
    if not cart:
        cart = await repositories.create_cart(db, session_id)
    return cart


async def mutate_cart(db: AsyncSession, session_id: str, intent: IntentResult) -> ChatResponse:
    # A tx will automatically start once the db session is created in SQLAlchemy 2.0,
    # so we don't have to manually create a tx
    # async with db.begin():
    cart = await lock_cart(db, session_id)
    if cart.status != CartStatus.OPEN:
        raise HTTPException(status_code=400, detail="Cart is closed")

    cart_items: list[CartItem] = []
    menu_item_ids: list[str] = list(map(lambda intent_item: intent_item.menu_item_id, intent.items))
    menu_items = await repositories.get_menu_item_by_menu_item_ids(db, menu_item_ids)
    menu_items_dic: dict[str, MenuItem] = {mi.id: mi for mi in menu_items}
    for item in intent.items:
        if item.quantity <= 0:
            raise HTTPException(status_code=400, detail="Quantity must be positive")
        unit_price_scaled = money_util.to_scaled_val(menu_items_dic[item.menu_item_id].price)
        cart_item = CartItem(
            id=str(uuid.uuid4()),
            cart_id=cart.id,
            menu_item_id=menu_items_dic[item.menu_item_id].id,
            name=menu_items_dic[item.menu_item_id].name,
            quantity=item.quantity,
            unit_price_scaled=unit_price_scaled,
            total_price_scaled=item.quantity * unit_price_scaled,
        )
        cart_items.append(cart_item)
    await repositories.upsert_cart_item(db, cart, cart_items)
    await db.commit()

    cart_summary = await build_cart_summary(cart)
    return ChatResponse(
        session_id=session_id,
        reply=response_builder.build_reply(intent, cart_summary),
        intent=intent,
        cart=cart_summary,
    )
