import os

import grpc
from src import repositories
from src.db import SessionLocal
from src.entities import Order
from src.enums import OrderStatus

from orderbot.v1 import common_pb2
from orderbot.v1.order_bot_svc import order_status_service_pb2_grpc


class OrderStatusGrpcServer(order_status_service_pb2_grpc.OrderStatusServiceServicer):
    """gRPC server boundary for callbacks from order-bot-mgmt-svc."""

    async def mark_order_completed(self, order_id: str) -> None:
        async with SessionLocal() as db:
            order = await db.get_one(Order, order_id)
            await repositories.update_order_status(db, order, OrderStatus.COMPLETED)
            await db.commit()

    async def MarkOrderCompleted(self, request, context):
        await self.mark_order_completed(request.order_id)
        return common_pb2.MarkOrderCompletedResponse(ok=True)


async def serve_grpc() -> grpc.aio.Server:
    server = grpc.aio.server()
    order_status_service_pb2_grpc.add_OrderStatusServiceServicer_to_server(OrderStatusGrpcServer(), server)
    server.add_insecure_port(f"[::]:{os.getenv('ORDER_BOT_GRPC_PORT', '50052')}")
    return server
