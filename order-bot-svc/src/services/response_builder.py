import json

from src.schemas import CartSummary, IntentResult


def build_reply(intent: IntentResult, cart: CartSummary) -> str:
    if not intent.valid:
        return "\n".join(
            [
                "Sorry, I didn't understand that. Could you clarify?",
                "I can only help you place an order.",
                "",
                "If you want to see the menu, click the button in the top-right corner.",
            ]
        )

    match intent.intent_type:
        case "search_menu":
            return "You can check the menu by clicking the button in the top-right corner."
        case "mutate_cart_items":
            return "\n".join(
                [
                    "Sure! I've updated your order.",
                    "",
                    "--------------------------------",
                    "",
                    "Here is your current order:",
                    _cart_pretty_json(cart),
                    "",
                    "If you'd like to submit your order, I'm happy to help.",
                ]
            )
        case "show_cart":
            if not cart.items:
                return "Sorry, you don't have any items in your order yet."
            return "\n".join(
                [
                    "Here is your current order:",
                    _cart_pretty_json(cart),
                    "",
                    "If you'd like to submit your order, I'm happy to help.",
                ]
            )
        case "checkout":
            if not cart.items:
                return "\n".join(
                    [
                        "Sorry, you don't have any items in your order yet.",
                        "Please order some items before checkout.",
                    ]
                )
            return "Ready to place your order. Please confirm to proceed."
        case _:
            return "Sorry, I didn't understand that."


def _cart_pretty_json(cart: CartSummary) -> str:
    payload = {
        "items": [
            {
                "name": i.name,
                "unit price": i.unit_price,
                "quantity": i.quantity,
                "total price": i.total_price,
            }
            for i in cart.items
        ],
        "total_price": cart.total_price,
    }
    return json.dumps(payload, ensure_ascii=False, indent=2)

