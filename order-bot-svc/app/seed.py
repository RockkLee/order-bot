from sqlalchemy.orm import Session
from app.models import MenuItem


def seed_menu(db: Session) -> None:
    existing = db.query(MenuItem).count()
    if existing:
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
