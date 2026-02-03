import uuid
from fastapi import APIRouter, Depends, Header, HTTPException
from starlette.responses import Response
from sqlalchemy.orm import Session
from app.db import get_db_session
from app import repositories
from app.schemas import ChatRequest, ChatResponse, IntentResult, MenuItemOut
from app.services.intent import IntentParser
from app.services.cart_service import build_cart_summary, ensure_cart, lock_cart, touch_cart
from app.services.response_builder import build_reply

router = APIRouter()
intent_parser = IntentParser()


@router.post("/chat", response_model=ChatResponse)
def chat(
    payload: ChatRequest,
    response: Response,
    session_id: str | None = Header(default=None, alias="Session-Id"),
    db: Session = Depends(get_db_session),
):
    if not session_id:
        session_id = str(uuid.uuid4())
        cart = ensure_cart(db, session_id)
        db.commit()
        response.headers["Session-Id"] = session_id
    else:
        cart = ensure_cart(db, session_id)

    cart_summary = build_cart_summary(cart)
    intent = intent_parser.parse(payload.message, has_cart_items=bool(cart.items))

    if not intent.valid:
        return ChatResponse(
            session_id=session_id,
            reply=build_reply(intent, cart_summary),
            intent=intent,
            cart=cart_summary,
        )

    if intent.intent_type == "search_menu":
        results = repositories.get_menu_by_query(db, intent.query or "")
        menu_out = [
            MenuItemOut(
                sku=item.sku,
                name=item.name,
                description=item.description,
                price_cents=item.price_cents,
            )
            for item in results
        ]
        reply = build_reply(intent, cart_summary)
        return ChatResponse(
            session_id=session_id,
            reply=reply,
            intent=intent,
            cart=cart_summary,
            menu_results=menu_out,
        )

    if intent.intent_type in {"add_item", "update_item", "remove_item"}:
        with db.begin():
            cart = lock_cart(db, session_id)
            if cart.status != "OPEN":
                raise HTTPException(status_code=400, detail="Cart is closed")

            for item in intent.items:
                menu_item = repositories.get_menu_item_by_sku(db, item.sku)
                if not menu_item or not menu_item.is_available:
                    raise HTTPException(status_code=404, detail=f"SKU not found: {item.sku}")

                if intent.intent_type == "remove_item":
                    repositories.remove_cart_item(db, cart, item.sku)
                    continue

                if item.quantity <= 0:
                    raise HTTPException(status_code=400, detail="Quantity must be positive")

                repositories.upsert_cart_item(
                    db,
                    cart,
                    sku=menu_item.sku,
                    name=menu_item.name,
                    quantity=item.quantity,
                    unit_price_cents=menu_item.price_cents,
                )

            touch_cart(cart)

        cart_summary = build_cart_summary(cart)
        reply = build_reply(intent, cart_summary)
        return ChatResponse(
            session_id=session_id,
            reply=reply,
            intent=intent,
            cart=cart_summary,
        )

    if intent.intent_type == "show_cart":
        cart_summary = build_cart_summary(cart)
        reply = build_reply(intent, cart_summary)
        return ChatResponse(
            session_id=session_id,
            reply=reply,
            intent=intent,
            cart=cart_summary,
        )

    if intent.intent_type == "checkout":
        if not intent.confirmed:
            reply = "Please confirm checkout by replying with 'confirm' or 'yes'."
            return ChatResponse(
                session_id=session_id,
                reply=reply,
                intent=intent,
                cart=cart_summary,
            )

        with db.begin():
            cart = lock_cart(db, session_id)
            if cart.status != "OPEN":
                raise HTTPException(status_code=400, detail="Cart is closed")
            if not cart.items:
                raise HTTPException(status_code=400, detail="Cart is empty")

            total_cents = sum(item.line_total_cents for item in cart.items)
            order = repositories.insert_order(db, cart, total_cents)
            repositories.insert_order_items(db, order, cart.items)

            cart.status = "CLOSED"
            touch_cart(cart)
            cart.closed_at = cart.updated_at

        cart_summary = build_cart_summary(cart)
        reply = f"Order placed! Your order id is {order.id}."
        return ChatResponse(
            session_id=session_id,
            reply=reply,
            intent=intent,
            cart=cart_summary,
            order_id=order.id,
        )

    return ChatResponse(
        session_id=session_id,
        reply=build_reply(IntentResult(valid=False, intent_type="unknown"), cart_summary),
        intent=intent,
        cart=cart_summary,
    )
