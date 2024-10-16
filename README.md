# Redis Clone

An implementation of the popular in-memory data store, Redis, written in Go.

This database can act as a drop-in replacement for a Redis server as any Redis client in the world would be able to interact with it.

> Note: This project is an attempt at diving into the internals of Redis and is for learning purposes only. It's not intended for any production level use.

## Run Locally

Clone the project

```bash
  git clone https://github.com/rajarshisg/redis-clone.git
```

Go to the project directory

```bash
  cd redis-clone
```

Build and Start the server

```bash
  make
```

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
