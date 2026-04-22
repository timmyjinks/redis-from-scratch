# redis-from-scratch

A lightweight Redis server implementation written from scratch in Go. This project implements the core Redis communication protocol (RESP) and a subset of Redis commands over a raw TCP connection — no external dependencies required.

## How It Works

The server listens on port **6379** (Redis's default port) and handles each client connection in its own goroutine. Incoming data is parsed using a custom RESP (REdis Serialization Protocol) reader, dispatched to a command handler, and the response is written back using a RESP writer.

### Project Structure

| File | Description |
|------|-------------|
| `main.go` | TCP server setup, connection loop, and request dispatcher |
| `handler.go` | Command handler registry and implementations (`PING`, `SET`, `GET`) |
| `resp.go` | RESP protocol parser and writer |

## Supported Commands

| Command | Usage | Description |
|---------|-------|-------------|
| `PING` | `PING [message]` | Returns `PONG`, or echoes back the optional message |
| `SET` | `SET key value` | Stores a key-value pair |
| `GET` | `GET key` | Retrieves the value for a given key |

## Getting Started

**Prerequisites:** Go 1.18+

```bash
# Clone the repo
git clone https://github.com/timmyjinks/redis-from-scratch.git
cd redis-from-scratch

# Run the server
go run .
```

The server will start listening on `localhost:6379`.

## Usage

You can interact with the server using the official `redis-cli` or any Redis client:

```bash
redis-cli PING
# => PONG

redis-cli SET hello world
# => OK

redis-cli GET hello
# => "world"
```
