# MemoDB

A Redis-compatible in-memory key-value store built from scratch in Go — implementing the RESP 2.0 wire protocol, master-replica replication, and RDB persistence.

Any standard Redis client (including `redis-cli`) can connect to MemoDB without modification.

> Built as a deep-dive into how Redis actually works internally — TCP server design, binary protocol parsing, concurrent connection handling, and distributed replication.

---

## Why I Built This

I wanted to understand what happens below the Redis API — how the RESP protocol works at the byte level, how a TCP server handles thousands of concurrent connections efficiently, how replication propagates commands from master to replicas, and how RDB files persist state to disk. Building it from scratch was the only way to actually learn this.

---

## Architecture

```
                        ┌─────────────────────────────────┐
                        │           MemoDB Server          │
                        │                                  │
  redis-cli ──TCP──►    │  net.Listen("tcp", port)         │
  any Redis client       │         │                        │
                        │         ▼                        │
                        │  tcpListener.Accept()            │
                        │         │                        │
                        │         ▼ (per connection)       │
                        │  go handleConnection()  ◄─────── goroutine per client
                        │         │                        │
                        │         ▼                        │
                        │  RESP Parser (commands pkg)      │
                        │         │                        │
                        │         ▼                        │
                        │  In-Memory Store (store pkg)     │
                        │         │                        │
                        │         ▼                        │
                        │  worker.PropagateCommand() ─────► replica 1
                        │                                  │ replica 2
                        └─────────────────────────────────┘
                                  │
                                  ▼ (on startup if RDB configured)
                          store.LoadRdbInStore()
```

---

## How It Works

### Concurrency — Goroutine Per Connection

Each incoming TCP connection is handled in a dedicated goroutine:

```go
go handleConnection(clientConn, false)
```

This is idiomatic Go — goroutines are lightweight (a few KB of stack vs MB for OS threads), making it practical to spawn one per connection. Connections are kept alive as long as data flows; idle connections with no activity for 30 seconds are automatically closed.

Replica connections are treated as persistent and handled with `persist: true` — they never time out.

### RESP 2.0 Protocol

MemoDB implements the Redis Serialization Protocol (RESP 2.0) — the binary protocol Redis uses for all client-server communication. This means:

- Any Redis client in the world can connect to MemoDB without code changes
- Commands are parsed from raw TCP bytes into structured command objects
- Responses are serialised back to RESP format before sending

### Master-Replica Replication

MemoDB supports a master-replica topology. Start a replica by passing the `--replicaof` flag:

```bash
# Start master on 6379
./memodb --port 6379

# Start replica pointing to master
./memodb --port 6380 --replicaof "127.0.0.1 6379"
```

Write commands (SET, DEL etc.) executed on the master are propagated to all connected replicas via `worker.PropagateCommand()`. Replica connections are maintained as persistent goroutines.

### RDB Persistence

MemoDB can load state from an RDB dump file on startup — the same binary format Redis uses for snapshots:

```bash
./memodb --dir /path/to/backups --dbfilename dump.rdb
```

On startup, `store.LoadRdbInStore()` parses the RDB file and restores the key-value store to its last persisted state.

---

## Supported Commands

| Command | Syntax | Description |
|---------|--------|-------------|
| `SET` | `SET key value` | Set a key to a string value |
| `GET` | `GET key` | Get the value of a key |
| `DEL` | `DEL key [key ...]` | Delete one or more keys |
| `EXISTS` | `EXISTS key` | Check if a key exists |
| `EXPIRE` | `EXPIRE key seconds` | Set TTL on a key |
| `TTL` | `TTL key` | Get remaining TTL of a key |
| `KEYS` | `KEYS pattern` | Find all keys matching a pattern |
| `CONFIG GET` | `CONFIG GET param` | Get server configuration |
| `PING` | `PING` | Test connection |
| `ECHO` | `ECHO message` | Echo a message back |

---

## Getting Started

**Prerequisites:** Go 1.21+ or Docker

### Run with Go

```bash
git clone https://github.com/rajarshisg/memodb.git
cd memodb
make
```

### Run with Docker

```bash
# Build image
make build

# Run container
make run

# Stop container
make stop
```

### Connect with redis-cli

```bash
redis-cli -p 6379
127.0.0.1:6379> SET name "Rajarshi"
OK
127.0.0.1:6379> GET name
"Rajarshi"
127.0.0.1:6379> EXPIRE name 60
(integer) 1
127.0.0.1:6379> TTL name
(integer) 58
```

### Run with RDB persistence

```bash
./memodb --port 6379 --dir ./backups --dbfilename dump.rdb
```

### Run master-replica setup

```bash
# Terminal 1 — start master
./memodb --port 6379

# Terminal 2 — start replica
./memodb --port 6380 --replicaof "127.0.0.1 6379"
```

---

## Project Structure

```
memodb/
├── server.go              # Entry point — TCP listener, connection handler, flag parsing
├── internal/
│   ├── commands/          # RESP protocol parser + command handlers
│   ├── store/             # In-memory key-value store + RDB persistence
│   ├── tcp/               # TCP connection management
│   └── worker/            # Master-replica replication worker
├── Dockerfile
├── Makefile
└── go.mod
```

---

## Utility Commands

```bash
make          # Build and run
make build    # Build Docker image
make run      # Run Docker container
make stop     # Stop Docker container
make clean    # Remove Docker image
make rebuild  # Force Docker image rebuild
```

---

## What I Learned

- **RESP protocol internals** — parsing a binary protocol from raw TCP bytes, handling bulk strings, arrays, and simple strings
- **TCP server design in Go** — goroutine-per-connection model, connection lifecycle management, idle timeout handling
- **Distributed replication** — how commands propagate from master to replicas, maintaining persistent replica connections
- **RDB file format** — Redis's binary snapshot format for state persistence
- **Go concurrency primitives** — goroutines and channels for concurrent connection handling without the overhead of OS threads

---

## Limitations

This is a learning project, not production-ready software:

- In-memory only — no AOF (Append Only File) persistence
- Single-threaded store — no concurrent write protection
- Limited command set — core commands only, no Pub/Sub, Streams, or Lua scripting
- No authentication or TLS
- Replication is best-effort — no acknowledgment or consistency guarantees

---

## Related Reading

- [Redis RESP Protocol Specification](https://redis.io/docs/reference/protocol-spec/)
- [Redis RDB File Format](https://rdb.fnordig.de/file_format.html)
- [Redis Replication Documentation](https://redis.io/docs/management/replication/)
