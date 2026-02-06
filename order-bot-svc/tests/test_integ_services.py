import unittest

from fastapi import Depends

from src.db import get_db_session
from src.entities import Cart
from src.schemas import IntentResult
from src.services import order_service


class OrderServiceTests(unittest.IsolatedAsyncioTestCase):
    def setUp(self):
        self.session_id = ""

    def test_check_out(self):
        intenRes = IntentResult(valid=True, intent_type="checkout")
        cart = Cart()
        order_service.checkout(Depends(get_db_session), self.session_id, intenRes, cart)
