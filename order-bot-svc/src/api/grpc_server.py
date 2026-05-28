import grpc
from src import repositories
from src.db import SessionLocal
from src.entities import Order
from src.enums import OrderStatus


class OrderStatusGrpcServer:
    """gRPC server boundary for callbacks from order-bot-mgmt-svc.

    TODO: bind to generated OrderStatusServiceServicer from proto/order_sync.proto.
    """

    async def mark_order_completed(self, order_id: str) -> None:
        async with SessionLocal() as db:
            order = await db.get_one(Order, order_id)
            await repositories.update_order_status(db, order, OrderStatus.COMPLETED)
            await db.commit()


async def serve_grpc() -> grpc.aio.Server:
    server = grpc.aio.server()
    # TODO: register generated servicer here.
    return server
