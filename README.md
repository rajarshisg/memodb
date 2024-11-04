# MemoDB

MemoDB is a toy key-value database. It's essentially an implementation of the popular in-memory data store, Redis, written in Go.

This database can act as a drop-in replacement for a Redis server as any Redis client in the world would be able to interact with it.

> Note: This project is an attempt at diving into the internals of Redis and is for learning purposes only. It's not intended for any production level use.

## Run Locally

Clone the project

```bash
  git clone https://github.com/rajarshisg/memodb.git
```

Go to the project directory

```bash
  cd memodb
```

Run the server

```bash
  make
```
## MemoDB in action
Server:
<img width="621" alt="Screenshot 2024-11-04 at 4 24 32 PM" src="https://github.com/user-attachments/assets/a116ece4-411f-44f8-aa11-fb15428ba577">
Using redis-cli:
<img width="239" alt="Screenshot 2024-11-04 at 4 24 53 PM" src="https://github.com/user-attachments/assets/aa7e1c55-891f-4c78-aff2-eb55683a6322">

## Utility commands for local development

Build the Docker image

```bash
  make build
```

Run the Docker container

```bash
  make run
```

Stop the Docker container

```bash
  make stop
```

Clean-up the Docker image

```bash
  make clean
```

Force Docker image re-build

```bash
  make rebuild
```
