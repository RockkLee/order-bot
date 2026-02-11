import asyncio

from src.intent.intent import MCPIntentClient, IntentParser

if __name__ == "__main__":
    async def main():
        # client = MCPIntentClient()
        # tools = await client.get_tools()
        # print(tools)

        res = await IntentParser().parse("Can I order, please?", False)
        print(res)


    asyncio.run(main())
