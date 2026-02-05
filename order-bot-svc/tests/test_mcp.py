import asyncio

from src.intent.intent import MCPIntentClient

if __name__ == "__main__":
    async def main():
        client = MCPIntentClient()
        result = await client.infer_intent("Hello World", True)
        print(result.model_dump())


    asyncio.run(main())
