from dataclasses import dataclass
import os


@dataclass
class CheckoutOrderItem:
    menu_item_id: str
    name: str
    quantity: int
    unit_price_scaled: int
    total_price_scaled: int


@dataclass
class CheckoutOrderRequest:
    order_id: str
    bot_id: str
    session_id: str
    total_scaled: int
    items: list[CheckoutOrderItem]


class OrderMgmtGrpcService:
    """gRPC integration boundary for order-bot-mgmt-svc.

    The actual gRPC transport can be plugged in later. For now this class contains
    the business-flow oriented method signatures used by order_service.
    """

    async def create_order_with_items(self, request: CheckoutOrderRequest) -> None:
        try:
            import grpc
        except ModuleNotFoundError:
            return
        target = os.getenv("ORDER_MGMT_GRPC_TARGET", "order-bot-mgmt-svc:50051")
        async with grpc.aio.insecure_channel(target) as channel:
            # TODO: replace with generated stub:
            # stub = grpcpb.OrderSyncServiceStub(channel)
            # await stub.CreateOrder(grpcpb.CreateOrderRequest(...))
            _ = (channel, request)

    async def notify_order_completed(self, order_id: str) -> None:
        try:
            import grpc
        except ModuleNotFoundError:
            return
        target = os.getenv("ORDER_MGMT_GRPC_TARGET", "order-bot-mgmt-svc:50051")
        async with grpc.aio.insecure_channel(target) as channel:
            # TODO: replace with generated stub:
            # stub = grpcpb.OrderStatusServiceStub(channel)
            # await stub.MarkOrderCompleted(grpcpb.MarkOrderCompletedRequest(order_id=order_id))
            _ = (channel, order_id)
