import asyncio
import os

os.environ.setdefault("OLLAMA_BASE_URL", "http://localhost:11434")

from langchain_ollama import ChatOllama
from browser_use import Agent


async def main():
    print("Starting BrowserUse agent...")

    llm = ChatOllama(
        model="llama3.2",
        temperature=0,
    )

    agent = Agent(
        task="Go to https://grok.com and tell me what you see on the page. Just describe if there's a chat input box visible.",
        llm=llm,
    )

    result = await agent.run()
    print("Result:", result)


if __name__ == "__main__":
    asyncio.run(main())
