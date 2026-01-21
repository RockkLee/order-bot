from app.schemas import CartSummary, IntentResult


def build_reply(intent: IntentResult, cart: CartSummary) -> str:
    if not intent.valid:
        return "Sorry, I didn't understand that. Could you clarify?"

    match intent.intent_type:
        case "search_menu":
            return "Here are the menu items that match your search."
        case "add_item":
            return "Added the item to your cart."
        case "update_item":
            return "Updated the item in your cart."
        case "remove_item":
            return "Removed the item from your cart."
        case "show_cart":
            if not cart.items:
                return "Your cart is empty."
            return "Here is your current cart."
        case "checkout":
            if not cart.items:
                return "Your cart is empty. Add items before checkout."
            return "Ready to place your order. Please confirm to proceed."
        case _:
            return "Sorry, I didn't understand that."
