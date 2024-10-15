# Redis Clone

An implementation of the popular in-memory data store, Redis, written in Go.

Note: This project is an attempt at diving into the internals of Redis and is still in development. It's not intended for any production level.

## Run Locally

Clone the project

```bash
  git clone https://github.com/rajarshisg/redis-clone.git
```

Go to the project directory

```bash
  cd redis-clone
```

Build the project

```bash
  docker build -t redis-clone .
```

Start the server

```bash
  docker run redis-clone
```
