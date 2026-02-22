import asyncio
import json
import logging

import grpc

from src.config import settings
from src.entities import Order, OrderItem
from src.grpc.contracts import ORDER_SYNC_METHOD

logger = logging.getLogger(__name__)


async def submit_order_async(order: Order, order_items: list[OrderItem]) -> None:
    payload = {
        "order_id": order.id,
        "cart_id": order.cart_id,
        "session_id": order.session_id,
        "total_scaled": order.total_scaled,
        "items": [
            {
                "id": item.id,
                "menu_item_id": item.menu_item_id,
                "name": item.name,
                "quantity": item.quantity,
                "unit_price_scaled": item.unit_price_scaled,
                "total_price_scaled": item.total_price_scaled,
            }
            for item in order_items
        ],
    }

    async with grpc.aio.insecure_channel(settings.order_mgmt_grpc_addr) as channel:
        stub = channel.unary_unary(
            ORDER_SYNC_METHOD,
            request_serializer=lambda req: json.dumps(req).encode("utf-8"),
            response_deserializer=lambda data: json.loads(data.decode("utf-8")),
        )
        await stub(payload)


def submit_order_background(order: Order, order_items: list[OrderItem]) -> None:
    async def _task() -> None:
        try:
            await submit_order_async(order, order_items)
        except Exception:
            logger.exception("failed to send order %s to mgmt grpc", order.id)

    asyncio.create_task(_task())
