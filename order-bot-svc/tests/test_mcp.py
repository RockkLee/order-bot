import asyncio

from src.intent.intent import MCPIntentClient

if __name__ == "__main__":
    async def main():
        client = MCPIntentClient()
        tools = await client.get_tools()
        print(tools)


    asyncio.run(main())
