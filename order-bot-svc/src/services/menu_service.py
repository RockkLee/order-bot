from src import repositories
from src.schemas import IntentResult, ChatResponse, MenuItemOut, MenuItemIntent
from src.services import cart_service
from src.services import response_builder
from sqlalchemy.ext.asyncio import AsyncSession


async def search_menu_for_intent(db: AsyncSession, menu_id: str) -> list[MenuItemIntent]:
    menu_items = await repositories.get_menu_by_query(db, menu_id)
    menu_item_intents = [
        MenuItemIntent(
            menu_item_id=item.id,
            name=item.name,
            price=item.price
        )
        for item in menu_items
    ]
    return menu_item_intents


async def search_menu(
    db: AsyncSession, menu_id: str, intent: IntentResult, cart
) -> ChatResponse:
    results = await repositories.get_menu_by_query(db, menu_id)
    menu_out = [
        MenuItemOut(
            name=item.name,
            price=item.price,
        )
        for item in results
    ]
    cart_summary = await cart_service.build_cart_summary(cart)
    return ChatResponse(
        session_id=cart.session_id,
        reply=response_builder.build_reply(intent, cart_summary),
        intent=intent,
        cart=cart_summary,
        menu_results=menu_out,
    )
