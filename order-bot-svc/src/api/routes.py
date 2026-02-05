import uuid
from fastapi import APIRouter, Depends, Header, Response
from sqlalchemy.ext.asyncio import AsyncSession

from src.db import get_db_session
from src.schemas import ChatRequest, ChatResponse, IntentResult
from src.services.cart_service import build_cart_summary, ensure_cart
from src.services.chat_agents import AgentContext, INTENT_AGENTS
from src.services.intent import IntentParser, LLMIntentClient
from src.services.response_builder import build_reply

router = APIRouter()
intent_parser = IntentParser()
llm_intent_client = LLMIntentClient()


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
    message = payload.message.strip()
    if not message:
        intent = IntentResult(valid=False, intent_type="unknown", reason="empty")
    else:
        llm_response = await llm_intent_client.infer_intent(message, has_cart_items=bool(cart.items))
        intent = await intent_parser.parse(llm_response, has_cart_items=bool(cart.items))

    if not intent.valid:
        return ChatResponse(
            session_id=session_id,
            reply=build_reply(intent, cart_summary),
            intent=intent,
            cart=cart_summary,
        )

    agent = INTENT_AGENTS.get(intent.intent_type)
    if agent:
        return await agent.run(AgentContext(db=db, session_id=session_id, intent=intent))

    return ChatResponse(
        session_id=session_id,
        reply=build_reply(IntentResult(valid=False, intent_type="unknown"), cart_summary),
        intent=intent,
        cart=cart_summary,
    )
