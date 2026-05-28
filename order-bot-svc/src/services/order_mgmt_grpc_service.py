from src.grpc.order_mgmt_client import OrderMgmtGrpcClient
from src.grpc.order_mgmt_types import CheckoutOrderItem, CheckoutOrderRequest


class OrderMgmtGrpcService:
    """gRPC integration boundary for order-bot-mgmt-svc.

    This wraps the generated gRPC client used by order_service.
    """

    def __init__(self, client: OrderMgmtGrpcClient | None = None):
        self.client = client or OrderMgmtGrpcClient()

    async def create_order_with_items(self, request: CheckoutOrderRequest) -> None:
        await self.client.create_order_with_items(request)

    async def notify_order_completed(self, order_id: str) -> None:
        _ = order_id
