import uuid
from importlib.metadata import always_iterable

from sqlalchemy.ext.asyncio import AsyncSession
from fastapi import APIRouter, Response, Header, Depends, HTTPException

from src import repositories
from src.db import get_db_session
from src.entities import Cart
from src.schemas import ChatRequest, ChatResponse, IntentResult
from src.services import menu_service, order_service
from src.services import cart_service
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
        cart = await cart_service.ensure_cart(db, session_id)
        await db.commit()
        response.headers["Session-Id"] = session_id
    else:
        cart = await cart_service.ensure_cart(db, session_id)

    cart_summary = cart_service.build_cart_summary(cart)
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
    *, db: AsyncSession, session_id: str, intent: IntentResult, cart: Cart
) -> ChatResponse:
    return await menu_service.search_menu(db, session_id, intent, cart)


async def _handle_cart_mutation(
    *, db: AsyncSession, session_id: str, intent: IntentResult, cart: Cart
) -> ChatResponse:
    return await cart_service.mutate_cart(db, session_id, intent)


async def _handle_show_cart(
    *, db: AsyncSession, session_id: str, intent: IntentResult, cart
) -> ChatResponse:
    cart_summary = cart_service.build_cart_summary(cart)
    return ChatResponse(
        session_id=session_id,
        reply=build_reply(intent, cart_summary),
        intent=intent,
        cart=cart_summary,
    )


async def _handle_checkout(
    *, db: AsyncSession, session_id: str, intent: IntentResult, cart: Cart
) -> ChatResponse:
    return await order_service.checkout(db, session_id, intent, cart)


_INTENT_HANDLERS = {
    "search_menu": _handle_search_menu,
    "add_item": _handle_cart_mutation,
    "update_item": _handle_cart_mutation,
    "remove_item": _handle_cart_mutation,
    "show_cart": _handle_show_cart,
    "checkout": _handle_checkout,
}
