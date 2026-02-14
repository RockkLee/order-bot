import logging
import os
from typing import Any

from dotenv import load_dotenv
from pydantic import BaseModel


load_dotenv()
logger = logging.getLogger(__name__)


class AppSettings(BaseModel):
    app_name: str = "Order Bot"
    host: str = os.environ["HOST"]
    port: int = int(os.environ["PORT"])
    root_path: str = os.environ["ROOT_PATH"]
    is_production: bool = os.environ["IS_PRODUCTION"].lower() == "true"
    database_url: str = os.environ["DATABASE_URL"]
    database_schema: str = os.environ["DATABASE_SCHEMA"]
    seed_menu: bool = os.environ["SEED_MENU"].lower() == "true"
    mistral_api_key: str = os.environ["MISTRAL_API_KEY"]
    mistral_model: str = os.environ["MISTRAL_MODEL"]


settings = AppSettings()
logger.info(
    "Loaded configuration from environment: %s",
    settings.model_dump(exclude={"mistral_api_key"}),
)
