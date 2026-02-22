from enum import Enum


class CartStatus(Enum):
    OPEN = "OPEN"
    CLOSED = "CLOSED"


class OrderStatus(Enum):
    PROCESSING = "PROCESSING"
    DONE = "DONE"
