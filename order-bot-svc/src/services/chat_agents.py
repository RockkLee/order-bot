from dataclasses import dataclass

from fastapi import HTTPException
from sqlalchemy.ext.asyncio import AsyncSession

from src import repositories
from src.schemas import ChatResponse, IntentResult, MenuItemOut
from src.services.cart_service import build_cart_summary, lock_cart, touch_cart
from src.services.response_builder import build_reply


@dataclass
class AgentContext:
    db: AsyncSession
    session_id: str
    intent: IntentResult


class BaseIntentAgent:
    async def run(self, context: AgentContext) -> ChatResponse:
        raise NotImplementedError


class SearchMenuAgent(BaseIntentAgent):
    async def run(self, context: AgentContext) -> ChatResponse:
        cart_summary = build_cart_summary(await repositories.get_cart_by_session(context.db, context.session_id))
        results = await repositories.get_menu_by_query(context.db, context.intent.query or "")
        menu_out = [
            MenuItemOut(
                sku=item.sku,
                name=item.name,
                description=item.description,
                price_cents=item.price_cents,
            )
            for item in results
        ]
        return ChatResponse(
            session_id=context.session_id,
            reply=build_reply(context.intent, cart_summary),
            intent=context.intent,
            cart=cart_summary,
            menu_results=menu_out,
        )


class AddItemAgent(BaseIntentAgent):
    async def run(self, context: AgentContext) -> ChatResponse:
        return await _apply_item_mutation(context, intent_type="add_item")


class UpdateItemAgent(BaseIntentAgent):
    async def run(self, context: AgentContext) -> ChatResponse:
        return await _apply_item_mutation(context, intent_type="update_item")


class RemoveItemAgent(BaseIntentAgent):
    async def run(self, context: AgentContext) -> ChatResponse:
        return await _apply_item_mutation(context, intent_type="remove_item")


async def _apply_item_mutation(context: AgentContext, intent_type: str) -> ChatResponse:
    async with context.db.begin():
        cart = await lock_cart(context.db, context.session_id)
        if cart.status != "OPEN":
            raise HTTPException(status_code=400, detail="Cart is closed")

        for item in context.intent.items:
            menu_item = await repositories.get_menu_item_by_sku(context.db, item.sku)
            if not menu_item or not menu_item.is_available:
                raise HTTPException(status_code=404, detail=f"SKU not found: {item.sku}")

            if intent_type == "remove_item":
                await repositories.remove_cart_item(context.db, cart, item.sku)
                continue

            if item.quantity <= 0:
                raise HTTPException(status_code=400, detail="Quantity must be positive")

            await repositories.upsert_cart_item(
                context.db,
                cart,
                sku=menu_item.sku,
                name=menu_item.name,
                quantity=item.quantity,
                unit_price_cents=menu_item.price_cents,
            )

        touch_cart(cart)

    cart = await repositories.get_cart_by_session(context.db, context.session_id)
    cart_summary = build_cart_summary(cart)
    return ChatResponse(
        session_id=context.session_id,
        reply=build_reply(context.intent, cart_summary),
        intent=context.intent,
        cart=cart_summary,
    )


class ShowCartAgent(BaseIntentAgent):
    async def run(self, context: AgentContext) -> ChatResponse:
        cart_summary = build_cart_summary(await repositories.get_cart_by_session(context.db, context.session_id))
        return ChatResponse(
            session_id=context.session_id,
            reply=build_reply(context.intent, cart_summary),
            intent=context.intent,
            cart=cart_summary,
        )


class CheckoutAgent(BaseIntentAgent):
    async def run(self, context: AgentContext) -> ChatResponse:
        cart = await repositories.get_cart_by_session(context.db, context.session_id)
        cart_summary = build_cart_summary(cart)
        if not context.intent.confirmed:
            return ChatResponse(
                session_id=context.session_id,
                reply="Please confirm checkout by replying with 'confirm' or 'yes'.",
                intent=context.intent,
                cart=cart_summary,
            )

        async with context.db.begin():
            cart = await lock_cart(context.db, context.session_id)
            if cart.status != "OPEN":
                raise HTTPException(status_code=400, detail="Cart is closed")
            if not cart.items:
                raise HTTPException(status_code=400, detail="Cart is empty")

            total_cents = sum(item.line_total_cents for item in cart.items)
            order = await repositories.insert_order(context.db, cart, total_cents)
            await repositories.insert_order_items(context.db, order, cart.items)

            cart.status = "CLOSED"
            touch_cart(cart)
            cart.closed_at = cart.updated_at

        cart = await repositories.get_cart_by_session(context.db, context.session_id)
        cart_summary = build_cart_summary(cart)
        return ChatResponse(
            session_id=context.session_id,
            reply=f"Order placed! Your order id is {order.id}.",
            intent=context.intent,
            cart=cart_summary,
            order_id=order.id,
        )


INTENT_AGENTS: dict[str, BaseIntentAgent] = {
    "search_menu": SearchMenuAgent(),
    "add_item": AddItemAgent(),
    "update_item": UpdateItemAgent(),
    "remove_item": RemoveItemAgent(),
    "show_cart": ShowCartAgent(),
    "checkout": CheckoutAgent(),
}
