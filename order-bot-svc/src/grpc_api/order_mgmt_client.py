import os

import grpc

from src.config import settings
from src.grpc_api.order_mgmt_types import CheckoutOrderRequest

from orderbot.v1 import common_pb2
from orderbot.v1.order_bot_mgmt_svc import order_sync_service_pb2_grpc


class OrderMgmtGrpcClient:
    def __init__(self, target: str | None = None):
        self.target = target or settings.grpc_target_order_bot_mgmt_svc

    async def create_order_with_items(self, request: CheckoutOrderRequest) -> None:
        async with grpc.aio.insecure_channel(self.target) as channel:
            stub = order_sync_service_pb2_grpc.OrderSyncServiceStub(channel)
            response = await stub.CreateOrder(
                common_pb2.CreateOrderRequest(
                    order_id=request.order_id,
                    bot_id=request.bot_id,
                    cart_id=request.cart_id,
                    session_id=request.session_id,
                    total_scaled=request.total_scaled,
                    items=[
                        common_pb2.OrderItem(
                            menu_item_id=item.menu_item_id,
                            name=item.name,
                            quantity=item.quantity,
                            unit_price_scaled=item.unit_price_scaled,
                            total_price_scaled=item.total_price_scaled,
                        )
                        for item in request.items
                    ],
                )
            )
            if not response.ok:
                raise RuntimeError("OrderSyncService.CreateOrder returned ok=false")
