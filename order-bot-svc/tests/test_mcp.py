import asyncio

from src.intent.intent import MCPIntentClient, IntentParser
from src.schemas import MenuItemIntent, CartItemIntent

if __name__ == "__main__":
    async def main():
        # client = MCPIntentClient()
        # tools = await client.get_tools()
        # print(tools)

        menu_item_intents = [
            MenuItemIntent(
                menu_item_id="menu_item_id_1",
                name="Latte",
                price=4.5
            ),
            MenuItemIntent(
                menu_item_id="menu_item_id_2",
                name="Espresso",
                price=3.5
            ),
        ]
        cart_item_intents = [
            CartItemIntent(
                menu_item_id="menu_item_id_1",
                quantity=1
            )
        ]
        mutate_cart_msg = "I want one cup of latte and 2 espresso"
        show_cart_msg = "I want to know what did I order"
        checkout_msg = "Checkout, please"
        unknown_msg = "abbabllr"
        res = await IntentParser().parse(
            unknown_msg,
            menu_item_intents,
            True,
            cart_item_intents
        )
        print(res)


    asyncio.run(main())
