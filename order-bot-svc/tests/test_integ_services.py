import unittest
from typing import cast
from datetime import datetime, timedelta, UTC

from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker
from sqlalchemy.pool import StaticPool

from src.db import Base
from src.entities import Cart, CartItem, MenuItem
from src.enums import CartStatus
from src.schemas import IntentItem, IntentResult
from src.services import cart_service, menu_service, order_service
from src.utils import money_util


class AsyncServiceTestCase(unittest.IsolatedAsyncioTestCase):
    async def asyncSetUp(self):
        self.engine = create_async_engine(
            "sqlite+aiosqlite:///:memory:",
            connect_args={"check_same_thread": False},
            poolclass=StaticPool,
        )
        async with self.engine.begin() as conn:
            await conn.run_sync(Base.metadata.create_all)
        self.SessionLocal = async_sessionmaker(bind=self.engine, expire_on_commit=False)

    async def asyncTearDown(self):
        await self.engine.dispose()

    async def create_menu_item(
        self,
        session,
        item_id,
        menu_id,
        name,
        price,
    ):
        menu_item = MenuItem(id=item_id, menu_id=menu_id, name=name, price=price)
        session.add(menu_item)
        await session.flush()
        return menu_item

    async def create_cart(self, session, session_id="session-1"):
        cart = Cart(session_id=session_id)
        session.add(cart)
        await session.flush()
        return cart

    async def add_cart_item(self, session, cart, menu_item, quantity=2):
        scaled_val = money_util.to_scaled_val(menu_item.price)
        cart_item = CartItem(
            cart_id=cart.id,
            menu_item_id=menu_item.id,
            name=menu_item.name,
            quantity=quantity,
            unit_price_scaled=scaled_val,
            total_price_scaled=quantity * scaled_val,
        )
        session.add(cart_item)
        await session.flush()
        return cart_item


class CartServiceTests(AsyncServiceTestCase):
    async def test_build_cart_summary_aggregates_items(self):
        cart = Cart(
            id="e4d61b28-3dfa-4823-be10-6d53d56ea4d0",
            session_id="56c96d21-439d-4601-835b-df9127ac2e35",
            status=CartStatus.OPEN.value,
            total_scaled=1800,
            closed_at=None,
        )
        item_one = CartItem(
            cart_id="9f7480df-4b4d-4882-8e6e-6f517eee7f82",
            menu_item_id="6a569db8-74ef-4ba5-85aa-16102cc30f69",
            name="Latte",
            quantity=3,
            unit_price_scaled=450,
            total_price_scaled=1350,
        )
        item_two = CartItem(
            cart_id="cb484e7b-cc4f-446d-a4b3-4e47adc17dfc",
            menu_item_id="cf73287b-6d5d-4839-80cf-94a8e5ff6a7e",
            name="Mocha",
            quantity=1,
            unit_price_scaled=450,
            total_price_scaled=450,
        )
        cart.items = [item_one, item_two]

        summary = await cart_service.build_cart_summary(cart)

        self.assertEqual(summary.total_price_scaled, 1800)
        self.assertEqual(summary.total_price, 18.0)
        self.assertEqual(summary.items[0].menu_item_id, "6a569db8-74ef-4ba5-85aa-16102cc30f69")
        self.assertEqual(summary.items[1].menu_item_id, "cf73287b-6d5d-4839-80cf-94a8e5ff6a7e")

    async def test_ensure_cart_creates_when_missing(self):
        async with self.SessionLocal() as session:
            cart = await cart_service.get_cart(session, "session-1")

        self.assertEqual(cart.session_id, "session-1")

    async def test_lock_cart_returns_existing(self):
        async with self.SessionLocal() as session:
            existing = await self.create_cart(session, "session-1")
            cart = await cart_service.lock_cart(session, "session-1")

        self.assertEqual(cart.id, existing.id)

    async def test_mutate_cart_adds_item(self):
        async with self.SessionLocal() as session:
            menu_item = await self.create_menu_item(
                session, item_id="item_id-1", menu_id="menu_id-1", name="menu_item_name", price=3.5)
            intent = IntentResult(
                valid=True,
                intent_type="mutate_item",
                items=[IntentItem(menu_item_id=menu_item.id, quantity=2)],
            )

            response = await cart_service.mutate_cart(session, "session-1", intent)

        self.assertEqual(response.cart.items[0].sku, "sku-1")
        self.assertEqual(response.cart.items[0].line_total_cents, 700)

    async def test_mutate_cart_removes_item(self):
        async with self.SessionLocal() as session:
            menu_item = await self.create_menu_item(
                session, item_id="item_id-1", menu_id="menu_id-1", name="menu_item_name", price=3.5)
            cart = await self.create_cart(session, "session-2")
            await self.add_cart_item(session, cart, menu_item, quantity=1)
            intent = IntentResult(
                valid=True,
                intent_type="mutate_item",
                items=[IntentItem(menu_item_id=menu_item.id, quantity=1)],
            )

            response = await cart_service.mutate_cart(session, "session-2", intent)

        self.assertEqual(response.cart.items, [])


class MenuServiceTests(AsyncServiceTestCase):
    async def test_search_menu_returns_results(self):
        async with self.SessionLocal() as session:
            await self.create_menu_item(session, item_id="latte-1", name="Latte", menu_id="menu_id_1", price=4.5)
            await self.create_menu_item(session, item_id="mocha-1", name="Mocha", menu_id="menu_id_2", price=5.5)
            cart = await self.create_cart(session, "session-3")

            intent = IntentResult(valid=True, intent_type="search_menu", query="Latte")
            response = await menu_service.search_menu(session, "menu_id_1", intent, cart)

        self.assertEqual(len(response.menu_results), 1)
        self.assertEqual(response.menu_results[0].name, "Latte")
        self.assertEqual(response.menu_results[0].price, 4.5)


class OrderServiceTests(AsyncServiceTestCase):
    async def test_checkout_requires_confirmation(self):
        cart = Cart(session_id="session-4")
        cart.items = [
            CartItem(
                cart_id="cart-4",
                menu_item_id="sku-4",
                name="Latte",
                quantity=1,
                unit_price_cents=450,
                line_total_cents=450,
            )
        ]
        intent = IntentResult(valid=True, intent_type="checkout", confirmed=False)

        response = await order_service.checkout(None, "session-4", intent, cart)

        self.assertIn("confirm", response.reply.lower())
        self.assertIsNone(response.order_id)

    async def test_checkout_places_order(self):
        async with self.SessionLocal() as session:
            menu_item = await self.create_menu_item(session, item_id="sku-5", menu_id="latte-5")
            cart = await self.create_cart(session, "session-5")
            await self.add_cart_item(session, cart, menu_item, quantity=2)
            _ = await cart.awaitable_attrs.items
            intent = IntentResult(valid=True, intent_type="checkout", confirmed=True)

            response = await order_service.checkout(session, "session-5", intent, cart)

        self.assertIsNotNone(response.order_id)
        self.assertEqual(response.cart.status, CartStatus.CLOSED)
