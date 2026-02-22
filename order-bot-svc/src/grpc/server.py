import json
import logging

import grpc

from src.config import settings
from src.db import SessionLocal
from src.enums import OrderStatus
from src.grpc.contracts import ORDER_CALLBACK_METHOD
from src import repositories

logger = logging.getLogger(__name__)


class OrderCallbackServer:
    def __init__(self):
        self._server = grpc.aio.server()
        handler = grpc.unary_unary_rpc_method_handler(
            self._update_order_status,
            request_deserializer=lambda data: json.loads(data.decode("utf-8")),
            response_serializer=lambda resp: json.dumps(resp).encode("utf-8"),
        )
        service = grpc.method_handlers_generic_handler(
            "ordersync.OrderCallbackService",
            {"UpdateOrderStatus": handler},
        )
        self._server.add_generic_rpc_handlers((service,))
        self._server.add_insecure_port(settings.order_svc_grpc_addr)

    async def _update_order_status(self, request: dict, context: grpc.aio.ServicerContext) -> dict:
        order_id = request.get("order_id", "")
        async with SessionLocal() as db:
            updated = await repositories.update_order_status(
                db,
                order_id,
                OrderStatus.PROCESSING,
                OrderStatus.DONE,
            )
            await db.commit()
        return {"updated": updated}

    async def start(self) -> None:
        await self._server.start()
        logger.info("order callback grpc server started at %s", settings.order_svc_grpc_addr)

    async def stop(self) -> None:
        await self._server.stop(5)
