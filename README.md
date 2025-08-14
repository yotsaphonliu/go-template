# OA Service

## Prerequisites

1. Golang 1.24.4

## How to run

1. Start dependencies

```sh
make start-deps
```

2. Run service

```sh
make setup-ssh-tunnel
make run-local
```

## Using Jaeger on Localhost

1. start docker compose
2. start api
3. go to http://localhost:16686/search