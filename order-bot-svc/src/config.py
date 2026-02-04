import logging
import os
from typing import Any

from dotenv import load_dotenv
from pydantic import BaseModel, Field


load_dotenv()
logger = logging.getLogger(__name__)


class AppSettings(BaseModel):
    app_name: str = "Order Bot"
    api_prefix: str = "/api"
    host: str = "127.0.0.1"
    port: int = 8000
    root_path: str = Field(default_factory=lambda: os.getenv("ROOT_PATH", ""))
    is_production: bool = Field(
        default_factory=lambda: os.getenv("IS_PRODUCTION", "false").lower() == "true"
    )
    logger_settings: dict[str, Any] | None = None
    database_url: str = Field(
        default_factory=lambda: os.getenv(
            "DATABASE_URL",
            "postgresql+asyncpg://postgres:postgres@localhost:5432/order_bot",
        )
    )
    seed_menu: bool = Field(
        default_factory=lambda: os.getenv("SEED_MENU", "true").lower() == "true"
    )
    mistral_api_key: str = Field(default_factory=lambda: os.getenv("MISTRAL_API_KEY", ""))
    mistral_model: str = Field(
        default_factory=lambda: os.getenv("MISTRAL_MODEL", "mistral-large-latest")
    )


settings = AppSettings()
logger.info(
    "Loaded configuration from environment: %s",
    settings.model_dump(exclude={"mistral_api_key"}),
)
