from sqlalchemy import select, func
from sqlalchemy.ext.asyncio import AsyncSession
from app.models import MenuItem


async def seed_menu(db: AsyncSession) -> None:
    result = await db.execute(select(func.count(MenuItem.id)))
    if result.scalar_one():
        return

    sample_items = [
        MenuItem(
            sku="coffee-latte",
            name="Latte",
            description="Espresso with steamed milk",
            price_cents=450,
        ),
        MenuItem(
            sku="coffee-espresso",
            name="Espresso",
            description="Single shot espresso",
            price_cents=300,
        ),
        MenuItem(
            sku="tea-matcha",
            name="Matcha Latte",
            description="Matcha with milk",
            price_cents=500,
        ),
        MenuItem(
            sku="pastry-croissant",
            name="Butter Croissant",
            description="Flaky pastry",
            price_cents=350,
        ),
    ]
    db.add_all(sample_items)
