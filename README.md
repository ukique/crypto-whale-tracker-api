# crypto-whale-tracker-api

A Go backend service that detects whale activity on MEXC exchange in real-time and exposes it via REST API.

## What is a whale?

A whale is a trade where the USD value exceeds a configured threshold (`price * quantity >= $50,000`).

## How it works

```
MEXC WebSocket → exchange/websocket.go → service.go (whale detection) → GET /whale
```

1. Connects to MEXC WebSocket and receives real-time trade data in Protobuf format
2. For each trade, calculates `USD = price * quantity`
3. If USD exceeds the threshold — it's a whale
4. Whale data is sent to the HTTP handler via Go channel
5. Client receives whale data via `GET /whale`

## Architecture

```
internal/
  core/
    config/       # Environment variables and thresholds
    errors/       # Error handling
  features/
    whale/
      models/     # Whale data structure
      service/    # Whale detection logic
      transport/
        exchange/ # MEXC WebSocket connection (read-only)
        handlers/ # HTTP handler, exposes /whale endpoint
proto/
  websocket-proto/ # MEXC Protobuf schemas
```

## Configuration

Create a `.env` file in the root:

```env
MEXC_WS_URL=wss://wbs-api.mexc.com/ws
TRADING_PAIR=BTCUSDT
WHALE_THRESHOLD=50000
```

## Running

```bash
go mod download
go run cmd/crypto-whale-tracker-api/main.go
```

## API

### GET /whale

Waits and returns the next detected whale trade.

> Note: This is a blocking request — it will wait until a whale trade is detected.

**Response:**
```json
{
  "whale": {
    "price": "66327.77",
    "quantity": "0.78345683",
    "symbol": "BTCUSDT",
    "time": 1772358404090
  }
}
```

## Notes

- MEXC WebSocket is **read-only** — we only consume market data
- Connection stays alive via PING every 15 seconds
- MEXC sends data in **Protobuf** format for market data, JSON for system messages
- Designed to be extended: add new exchanges in `transport/exchange/`, add new thresholds in config
