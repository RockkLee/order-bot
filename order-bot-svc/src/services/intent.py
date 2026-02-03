import re
from src.schemas import IntentResult, IntentItem


SKU_PATTERN = re.compile(r"sku[-_ ]?([\\w-]+)", re.IGNORECASE)
QUANTITY_PATTERN = re.compile(r"(\d+)")


class IntentParser:
    def parse(self, message: str, has_cart_items: bool) -> IntentResult:
        text = message.strip().lower()
        if not text:
            return IntentResult(valid=False, intent_type="unknown", reason="empty")

        if any(word in text for word in ["show cart", "view cart", "my cart", "cart"]):
            return IntentResult(valid=True, intent_type="show_cart")

        if any(word in text for word in ["checkout", "place order", "buy now", "confirm order", "confirm"]):
            confirmed = any(word in text for word in ["yes", "confirm", "place", "go ahead"])
            return IntentResult(valid=True, intent_type="checkout", confirmed=confirmed)

        if any(word in text for word in ["search", "menu", "have", "what"]):
            query = text.replace("search", "").replace("menu", "").strip()
            return IntentResult(valid=True, intent_type="search_menu", query=query or None)

        if any(word in text for word in ["remove", "delete", "drop"]):
            items = self._extract_items(text)
            if not items:
                return IntentResult(valid=False, intent_type="remove_item", reason="missing sku")
            return IntentResult(valid=True, intent_type="remove_item", items=items)

        if any(word in text for word in ["add", "update", "set"]):
            items = self._extract_items(text)
            if not items:
                return IntentResult(valid=False, intent_type="add_item", reason="missing sku")
            intent_type = "update_item" if "update" in text or "set" in text else "add_item"
            return IntentResult(valid=True, intent_type=intent_type, items=items)

        # fallback: if cart has items, treat as show cart clarification
        if has_cart_items:
            return IntentResult(valid=True, intent_type="show_cart")

        return IntentResult(valid=False, intent_type="unknown", reason="unrecognized")

    def _extract_items(self, text: str) -> list[IntentItem]:
        skus = SKU_PATTERN.findall(text)
        quantities = QUANTITY_PATTERN.findall(text)
        quantity = int(quantities[0]) if quantities else 1
        return [IntentItem(sku=sku, quantity=quantity) for sku in skus]
