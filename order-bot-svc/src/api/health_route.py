from fastapi import APIRouter

API_PREFIX = "/health"

router = APIRouter()


@router.get("/chk")
async def health_check() -> dict[str, str]:
    return {"status": "ok"}
