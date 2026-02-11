import uuid
from datetime import datetime, UTC
from sqlalchemy import String, Integer ,Double, Boolean, ForeignKey, DateTime, UniqueConstraint
from sqlalchemy.orm import Mapped, mapped_column, relationship
from src.db import Base


def _uuid_str() -> str:
    return str(uuid.uuid4())


class MenuItem(Base):
    __tablename__ = "menu_item"

    id: Mapped[str] = mapped_column(String(36), primary_key=True, default=_uuid_str)
    menu_id: Mapped[str] = mapped_column(String(64), unique=True, index=True)
    name: Mapped[str] = mapped_column(String(200), name="menu_item_name")
    price: Mapped[float] = mapped_column(Double)


class Cart(Base):
    __tablename__ = "cart"

    id: Mapped[str] = mapped_column(String(36), primary_key=True, default=_uuid_str)
    session_id: Mapped[str] = mapped_column(String(36), unique=True, index=True)
    status: Mapped[str] = mapped_column(String(20), default="OPEN")
    total_scaled: Mapped[int] = mapped_column(Integer, default=0)
    closed_at: Mapped[datetime | None] = mapped_column(DateTime, nullable=True)

    items: Mapped[list["CartItem"]] = relationship(
        "CartItem",
        back_populates="cart",
        cascade="all, delete-orphan",
    )


class CartItem(Base):
    __tablename__ = "cart_item"
    __table_args__ = (
        UniqueConstraint("cart_id", "menu_item_id", name="menu_item_id"),
    )

    id: Mapped[str] = mapped_column(String(36), primary_key=True, default=_uuid_str)
    cart_id: Mapped[str] = mapped_column(String(36), ForeignKey("cart.id"))
    menu_item_id: Mapped[str] = mapped_column(String(64))
    name: Mapped[str] = mapped_column(String(200))
    quantity: Mapped[int] = mapped_column(Integer)
    unit_price_scaled: Mapped[int] = mapped_column(Integer)
    total_price_scaled: Mapped[int] = mapped_column(Integer)

    cart: Mapped[Cart] = relationship("Cart", back_populates="items")


class Order(Base):
    __tablename__ = "orders"

    id: Mapped[str] = mapped_column(String(36), primary_key=True, default=_uuid_str)
    cart_id: Mapped[str] = mapped_column(String(36), ForeignKey("cart.id"))
    session_id: Mapped[str] = mapped_column(String(36), index=True)
    total_scaled: Mapped[int] = mapped_column(Integer)

    order_items: Mapped[list["OrderItem"]] = relationship(
        "OrderItem",
        back_populates="order",
        cascade="all, delete-orphan",
    )


class OrderItem(Base):
    __tablename__ = "order_item"

    id: Mapped[str] = mapped_column(String(36), primary_key=True, default=_uuid_str)
    order_id: Mapped[str] = mapped_column(String(36), ForeignKey("orders.id"))
    menu_item_id: Mapped[str] = mapped_column(String(64))
    name: Mapped[str] = mapped_column(String(200))
    quantity: Mapped[int] = mapped_column(Integer)
    unit_price_scaled: Mapped[int] = mapped_column(Integer)
    total_price_scaled: Mapped[int] = mapped_column(Integer)

    order: Mapped[Order] = relationship("Order", back_populates="order_items")
