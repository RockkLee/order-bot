from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from src.config import settings
from src.db import engine, SessionLocal, Base
from src.seed import seed_menu
from src.api.routes import router


def create_app() -> FastAPI:
    @asynccontextmanager
    async def lifespan(_: FastAPI):
        async with engine.begin() as conn:
            await conn.run_sync(Base.metadata.create_all)
        if settings.seed_menu:
            async with SessionLocal() as session:
                await seed_menu(session)
                await session.commit()
        yield

    app_root_path = settings.root_path
    app = FastAPI(
        title="Justka Q&A Bot",
        description="Justka Q&A Bot",
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

    app.include_router(router, prefix=settings.api_prefix)
    return app


app = create_app()

if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        app,
        host=settings.host,
        port=settings.port,
        log_config=settings.logger_settings,
        workers=1,
    )
