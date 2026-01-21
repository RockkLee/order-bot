from fastapi import FastAPI
from app.config import settings
from app.db import engine, SessionLocal, Base
from app.seed import seed_menu
from app.api.routes import router


def create_app() -> FastAPI:
    app = FastAPI(title=settings.app_name)

    @app.on_event("startup")
    def startup() -> None:
        Base.metadata.create_all(bind=engine)
        if settings.seed_menu:
            with SessionLocal() as session:
                seed_menu(session)
                session.commit()

    app.include_router(router, prefix=settings.api_prefix)
    return app


app = create_app()
