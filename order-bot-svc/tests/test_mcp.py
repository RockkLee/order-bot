import asyncio

from src.intent.intent import MCPIntentClient, IntentParser
from src.schemas import MenuItemIntent

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
        res = await IntentParser().parse("abbabllr", False, menu_item_intents)
        print(res)


    asyncio.run(main())
