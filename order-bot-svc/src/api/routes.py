import uuid

from sqlalchemy.ext.asyncio import AsyncSession
from fastapi import APIRouter, Response, Header, Depends, HTTPException

from src import repositories
from src.db import get_db_session
from src.schemas import ChatRequest, ChatResponse, IntentResult
from src.services import menu_service
from src.services.cart_service import build_cart_summary, ensure_cart, lock_cart, touch_cart
from src.intent.intent import IntentParser
from src.services.response_builder import build_reply

router = APIRouter()
intent_parser = IntentParser()


@router.post("/chat", response_model=ChatResponse)
async def chat(
    payload: ChatRequest,
    response: Response,
    session_id: str | None = Header(default=None, alias="Session-Id"),
    db: AsyncSession = Depends(get_db_session),
):
    if not session_id:
        session_id = str(uuid.uuid4())
        cart = await ensure_cart(db, session_id)
        await db.commit()
        response.headers["Session-Id"] = session_id
    else:
        cart = await ensure_cart(db, session_id)

    cart_summary = build_cart_summary(cart)
    intent = await intent_parser.parse(payload.message, has_cart_items=bool(cart.items))

    if not intent.valid:
        return ChatResponse(
            session_id=session_id,
            reply=build_reply(intent, cart_summary),
            intent=intent,
            cart=cart_summary,
        )

    handler = _INTENT_HANDLERS.get(intent.intent_type)
    if not handler:
        return ChatResponse(
            session_id=session_id,
            reply=build_reply(IntentResult(valid=False, intent_type="unknown"), cart_summary),
            intent=intent,
            cart=cart_summary,
        )

    return await handler(db=db, session_id=session_id, intent=intent, cart=cart)


async def _handle_search_menu(
    *, db: AsyncSession, session_id: str, intent: IntentResult, cart
) -> ChatResponse:
    return await menu_service.search_menu(db, session_id, intent, cart)


async def _handle_cart_mutation(
    *, db: AsyncSession, session_id: str, intent: IntentResult, cart
) -> ChatResponse:
    async with db.begin():
        cart = await lock_cart(db, session_id)
        if cart.status != "OPEN":
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
                menu_item_id=menu_item.sku,
                name=menu_item.name,
                quantity=item.quantity,
                unit_price_cents=menu_item.price_cents,
            )

        touch_cart(cart)

    cart_summary = build_cart_summary(cart)
    return ChatResponse(
        session_id=session_id,
        reply=build_reply(intent, cart_summary),
        intent=intent,
        cart=cart_summary,
    )


async def _handle_show_cart(
    *, db: AsyncSession, session_id: str, intent: IntentResult, cart
) -> ChatResponse:
    cart_summary = build_cart_summary(cart)
    return ChatResponse(
        session_id=session_id,
        reply=build_reply(intent, cart_summary),
        intent=intent,
        cart=cart_summary,
    )


async def _handle_checkout(
    *, db: AsyncSession, session_id: str, intent: IntentResult, cart
) -> ChatResponse:
    cart_summary = build_cart_summary(cart)
    if not intent.confirmed:
        reply = "Please confirm checkout by replying with 'confirm' or 'yes'."
        return ChatResponse(
            session_id=session_id,
            reply=reply,
            intent=intent,
            cart=cart_summary,
        )

    async with db.begin():
        cart = await lock_cart(db, session_id)
        if cart.status != "OPEN":
            raise HTTPException(status_code=400, detail="Cart is closed")
        if not cart.items:
            raise HTTPException(status_code=400, detail="Cart is empty")

        total_cents = sum(item.line_total_cents for item in cart.items)
        order = await repositories.insert_order(db, cart, total_cents)
        await repositories.insert_order_items(db, order, cart.items)

        cart.status = "CLOSED"
        touch_cart(cart)
        cart.closed_at = cart.updated_at

    cart_summary = build_cart_summary(cart)
    return ChatResponse(
        session_id=session_id,
        reply=f"Order placed! Your order id is {order.id}.",
        intent=intent,
        cart=cart_summary,
        order_id=order.id,
    )


_INTENT_HANDLERS = {
    "search_menu": _handle_search_menu,
    "add_item": _handle_cart_mutation,
    "update_item": _handle_cart_mutation,
    "remove_item": _handle_cart_mutation,
    "show_cart": _handle_show_cart,
    "checkout": _handle_checkout,
}
