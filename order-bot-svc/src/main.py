from contextlib import asynccontextmanager
import logging
import os

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from src.config import settings
from src.db import engine, SessionLocal, Base
from src.api import routes
from src.api import health_route


_LOG_LEVEL = os.getenv("LOG_LEVEL", "INFO").upper()
logging.basicConfig(
    level=_LOG_LEVEL,
    format="%(asctime)s %(levelname)s %(name)s - %(message)s",
)


def create_app() -> FastAPI:
    @asynccontextmanager
    async def lifespan(_: FastAPI):
        async with engine.begin() as conn:
            await conn.run_sync(Base.metadata.create_all)
        yield

    app_root_path = settings.root_path
    app = FastAPI(
        title="Order Bot",
        description="Order Bot",
        root_path=app_root_path,
        lifespan=lifespan,
    )

    if settings.is_production:
        app.openapi_url = None
        app.docs_url = None
        app.redoc_url = None

    app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"],
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
        expose_headers=["*"],
    )


    app.include_router(routes.router, prefix=routes.API_PREFIX)
    app.include_router(health_route.router, prefix=health_route.API_PREFIX)

    return app


app = create_app()

if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        app,
        host=settings.host,
        port=settings.port,
        workers=1,
    )
