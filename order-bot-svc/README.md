# Order Bot Service

A FastAPI-based order bot that stores carts/orders in PostgreSQL and follows the flow in the provided Mermaid diagram. The LLM step is modeled by a simple intent parser that can be swapped with a real LLM integration.

## Features
- Session-aware cart management via `Session-Id` header
- Menu search, add/update/remove items, show cart, and checkout
- PostgreSQL persistence for carts and orders
- Unit tests using `unittest`

## Quick Start

```bash
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt

# Update DATABASE_URL if needed
export DATABASE_URL="postgresql+psycopg://postgres:postgres@localhost:5432/order_bot"

uvicorn app.main:app --reload
```

## Example Request

```bash
curl -X POST http://127.0.0.1:8000/api/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "show cart"}'
```

The response includes a `session_id` and also returns the same value in the `Session-Id` response header. Pass it back in future calls.

## Running Tests

```bash
python -m unittest
```
