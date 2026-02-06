import unittest
from datetime import datetime, timedelta

from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker
from sqlalchemy.pool import StaticPool

from src.db import Base
from src.entities import Cart, CartItem, MenuItem
from src.schemas import IntentItem, IntentResult
from src.services import cart_service, menu_service, order_service


if not hasattr(CartItem, "sku"):
    CartItem.sku = property(lambda self: self.menu_item_id)


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
        item_id="menu-item-1",
        menu_id="latte",
        name="Latte",
        price=450,
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
        cart_item = CartItem(
            cart_id=cart.id,
            menu_item_id=menu_item.id,
            name=menu_item.name,
            quantity=quantity,
            unit_price_cents=menu_item.price,
            line_total_cents=quantity * menu_item.price,
        )
        session.add(cart_item)
        await session.flush()
        return cart_item


class CartServiceTests(AsyncServiceTestCase):
    async def test_build_cart_summary_aggregates_items(self):
        cart = Cart(session_id="session-1")
        item_one = CartItem(
            cart_id="cart-1",
            menu_item_id="sku-1",
            name="Latte",
            quantity=2,
            unit_price_cents=450,
            line_total_cents=900,
        )
        item_two = CartItem(
            cart_id="cart-1",
            menu_item_id="sku-2",
            name="Mocha",
            quantity=1,
            unit_price_cents=500,
            line_total_cents=500,
        )
        cart.items = [item_one, item_two]

        summary = cart_service.build_cart_summary(cart)

        self.assertEqual(summary.total_cents, 1400)
        self.assertEqual(summary.items[0].sku, "sku-1")
        self.assertEqual(summary.items[1].sku, "sku-2")

    async def test_touch_cart_updates_timestamp(self):
        cart = Cart(session_id="session-1")
        cart.updated_at = datetime.utcnow() - timedelta(minutes=5)

        cart_service.touch_cart(cart)

        self.assertGreater(cart.updated_at, datetime.utcnow() - timedelta(minutes=1))

    async def test_ensure_cart_creates_when_missing(self):
        async with self.SessionLocal() as session:
            cart = await cart_service.ensure_cart(session, "session-1")

        self.assertEqual(cart.session_id, "session-1")

    async def test_lock_cart_returns_existing(self):
        async with self.SessionLocal() as session:
            existing = await self.create_cart(session, "session-1")
            cart = await cart_service.lock_cart(session, "session-1")

        self.assertEqual(cart.id, existing.id)

    async def test_mutate_cart_adds_item(self):
        async with self.SessionLocal() as session:
            menu_item = await self.create_menu_item(session, item_id="sku-1")
            intent = IntentResult(
                valid=True,
                intent_type="add_item",
                items=[IntentItem(sku=menu_item.id, quantity=2)],
            )

            response = await cart_service.mutate_cart(session, "session-1", intent)

        self.assertEqual(response.cart.items[0].sku, "sku-1")
        self.assertEqual(response.cart.items[0].line_total_cents, 900)

    async def test_mutate_cart_removes_item(self):
        async with self.SessionLocal() as session:
            menu_item = await self.create_menu_item(session, item_id="sku-2")
            cart = await self.create_cart(session, "session-2")
            await self.add_cart_item(session, cart, menu_item, quantity=1)
            intent = IntentResult(
                valid=True,
                intent_type="remove_item",
                items=[IntentItem(sku=menu_item.id, quantity=1)],
            )

            response = await cart_service.mutate_cart(session, "session-2", intent)

        self.assertEqual(response.cart.items, [])


class MenuServiceTests(AsyncServiceTestCase):
    async def test_search_menu_returns_results(self):
        async with self.SessionLocal() as session:
            await self.create_menu_item(session, item_id="latte-1", name="Latte", menu_id="latte")
            await self.create_menu_item(session, item_id="mocha-1", name="Mocha", menu_id="mocha")
            cart = await self.create_cart(session, "session-3")

            intent = IntentResult(valid=True, intent_type="search_menu", query="Latte")
            response = await menu_service.search_menu(session, "session-3", intent, cart)

        self.assertEqual(len(response.menu_results), 1)
        self.assertEqual(response.menu_results[0].name, "Latte")


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
            _ = cart.items
            intent = IntentResult(valid=True, intent_type="checkout", confirmed=True)

            response = await order_service.checkout(session, "session-5", intent, cart)

        self.assertIsNotNone(response.order_id)
        self.assertEqual(response.cart.status, "CLOSED")
