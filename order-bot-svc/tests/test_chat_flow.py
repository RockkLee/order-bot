import unittest
from fastapi import FastAPI
from starlette.testclient import TestClient
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker

from src.api.routes import router
from src.db import Base, get_db_session
from src.entities import MenuItem


class OrderBotServiceTests(unittest.TestCase):
    def setUp(self):
        self.engine = create_engine("sqlite+pysqlite:///:memory:")
        Base.metadata.create_all(bind=self.engine)
        self.SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=self.engine)

        app = FastAPI()
        app.include_router(router, prefix="/api")

        def override_get_db():
            db = self.SessionLocal()
            try:
                yield db
            finally:
                db.close()

        app.dependency_overrides[get_db_session] = override_get_db
        self.client = TestClient(app)

        with self.SessionLocal() as session:
            session.add(
                MenuItem(
                    sku="coffee-latte",
                    name="Latte",
                    description="Espresso with milk",
                    price_cents=450,
                )
            )
            session.commit()

    def test_session_created_when_missing(self):
        response = self.client.post("/api/chat", json={"message": "show cart"})
        self.assertEqual(response.status_code, 200)
        payload = response.json()
        self.assertIn("session_id", payload)
        self.assertIn("Session-Id", response.headers)
        self.assertEqual(payload["cart"]["items"], [])

    def test_add_and_checkout_flow(self):
        create_response = self.client.post("/api/chat", json={"message": "show cart"})
        session_id = create_response.json()["session_id"]

        add_response = self.client.post(
            "/api/chat",
            json={"message": "add 2 sku coffee-latte"},
            headers={"Session-Id": session_id},
        )
        self.assertEqual(add_response.status_code, 200)
        add_payload = add_response.json()
        self.assertEqual(add_payload["cart"]["items"][0]["quantity"], 2)

        checkout_response = self.client.post(
            "/api/chat",
            json={"message": "checkout"},
            headers={"Session-Id": session_id},
        )
        self.assertEqual(checkout_response.status_code, 200)
        self.assertIn("confirm", checkout_response.json()["reply"].lower())

        confirm_response = self.client.post(
            "/api/chat",
            json={"message": "confirm"},
            headers={"Session-Id": session_id},
        )
        self.assertEqual(confirm_response.status_code, 200)
        confirm_payload = confirm_response.json()
        self.assertIsNotNone(confirm_payload.get("order_id"))
        self.assertEqual(confirm_payload["cart"]["status"], "CLOSED")


if __name__ == "__main__":
    unittest.main()
