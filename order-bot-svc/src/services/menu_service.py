from src import repositories
from src.schemas import IntentResult, ChatResponse, MenuItemOut
from src.services import cart_service
from src.services import response_builder
from sqlalchemy.ext.asyncio import AsyncSession


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
    cart_summary = cart_service.build_cart_summary(cart)
    return ChatResponse(
        session_id=cart.session_id,
        reply=response_builder.build_reply(intent, cart_summary),
        intent=intent,
        cart=cart_summary,
        menu_results=menu_out,
    )
