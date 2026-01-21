import os
from pydantic import BaseModel, Field


class AppSettings(BaseModel):
    app_name: str = "Order Bot"
    api_prefix: str = "/api"
    host: str = "127.0.0.1"
    port: int = 8000
    database_url: str = Field(
        default_factory=lambda: os.getenv(
            "DATABASE_URL",
            "postgresql+psycopg://postgres:postgres@localhost:5432/order_bot",
        )
    )
    seed_menu: bool = Field(
        default_factory=lambda: os.getenv("SEED_MENU", "true").lower() == "true"
    )


settings = AppSettings()
