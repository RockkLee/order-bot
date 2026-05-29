from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class OrderItem(_message.Message):
    __slots__ = ("id", "menu_item_id", "name", "quantity", "unit_price_scaled", "total_price_scaled")
    ID_FIELD_NUMBER: _ClassVar[int]
    MENU_ITEM_ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    QUANTITY_FIELD_NUMBER: _ClassVar[int]
    UNIT_PRICE_SCALED_FIELD_NUMBER: _ClassVar[int]
    TOTAL_PRICE_SCALED_FIELD_NUMBER: _ClassVar[int]
    id: str
    menu_item_id: str
    name: str
    quantity: int
    unit_price_scaled: int
    total_price_scaled: int
    def __init__(self, id: _Optional[str] = ..., menu_item_id: _Optional[str] = ..., name: _Optional[str] = ..., quantity: _Optional[int] = ..., unit_price_scaled: _Optional[int] = ..., total_price_scaled: _Optional[int] = ...) -> None: ...

class CreateOrderRequest(_message.Message):
    __slots__ = ("order_id", "bot_id", "cart_id", "session_id", "total_scaled", "items")
    ORDER_ID_FIELD_NUMBER: _ClassVar[int]
    BOT_ID_FIELD_NUMBER: _ClassVar[int]
    CART_ID_FIELD_NUMBER: _ClassVar[int]
    SESSION_ID_FIELD_NUMBER: _ClassVar[int]
    TOTAL_SCALED_FIELD_NUMBER: _ClassVar[int]
    ITEMS_FIELD_NUMBER: _ClassVar[int]
    order_id: str
    bot_id: str
    cart_id: str
    session_id: str
    total_scaled: int
    items: _containers.RepeatedCompositeFieldContainer[OrderItem]
    def __init__(self, order_id: _Optional[str] = ..., bot_id: _Optional[str] = ..., cart_id: _Optional[str] = ..., session_id: _Optional[str] = ..., total_scaled: _Optional[int] = ..., items: _Optional[_Iterable[_Union[OrderItem, _Mapping]]] = ...) -> None: ...

class CreateOrderResponse(_message.Message):
    __slots__ = ("ok",)
    OK_FIELD_NUMBER: _ClassVar[int]
    ok: bool
    def __init__(self, ok: bool = ...) -> None: ...

class MarkOrderCompletedRequest(_message.Message):
    __slots__ = ("order_id",)
    ORDER_ID_FIELD_NUMBER: _ClassVar[int]
    order_id: str
    def __init__(self, order_id: _Optional[str] = ...) -> None: ...

class MarkOrderCompletedResponse(_message.Message):
    __slots__ = ("ok",)
    OK_FIELD_NUMBER: _ClassVar[int]
    ok: bool
    def __init__(self, ok: bool = ...) -> None: ...
