# Go Template Service

A lightweight Go (Fiber) template for building HTTP APIs with PostgreSQL, MinIO object storage, and OpenTelemetry tracing (Jaeger). It includes a simple auth flow for a root/admin user, database migrations, and common middlewares.

## Prerequisites

- Go 1.25+
- Docker and Docker Compose
- Make (optional; you can also run docker compose and Go commands directly)

## Quick Start

1) Start dependencies (PostgreSQL, MinIO, Jaeger)

```ps1
# Using Make (recommended)
make start-deps

# Or directly
docker compose up -d
```

2) Configure the application

- Edit cfg\config.yaml as needed.
- Default HTTP port is 9092 (API.HTTPServerPort).
- By default, docker-compose creates a Postgres database named "go-template". Make sure cfg\config.yaml -> Database.PostgreSQL.DBName matches it (set to "go-template") or adjust docker-compose accordingly.
- OpenTelemetry endpoint defaults to localhost:4318 (no config needed). You can override via OTEL_EXPORTER_OTLP_ENDPOINT.

3) Run the API locally

Option A: Using Go directly (works on Windows PowerShell)

```ps1
# From the repository root
$env:OTEL_EXPORTER_OTLP_ENDPOINT="localhost:4318"  # optional, this is the default
go run .\src\main.go serve-http-api --config .\cfg\config.yaml
```

Option B: Using Make (requires Bash/Git Bash/WSL for the script)

```sh
# This script builds and runs the API
make run-local
```

4) Stop dependencies when done

```ps1
make stop-deps
# or: docker compose down
```

## API Endpoints (default base URL: http://localhost:9092/api)

- GET /health-check
  - Basic liveness endpoint.

- POST /root-login
  - Body (JSON): { "username": "admin", "password": "P@ssw0rd" }
  - Returns: { "token": "..." }

- GET /me
  - Headers: Authorization: Bearer <token>

- POST /logout
  - Headers: Authorization: Bearer <token>

## Configuration

Edit cfg\config.yaml:
- Log: level, color, JSON format
- Database: PostgreSQL host/port/user/pass/dbname (set DBName to go-template for docker-compose default)
- Minio: endpoint, user, password, bucket, UseSSL
- API: HTTPServerPort (default 9092)
- Admin: root credentials used by /api/root-login
- HashiCorp (optional): commented examples for Vault integration

You can also override settings via environment variables (viper with dot->underscore replacement). For example: API.HTTPServerPort -> API_HTTPServerPort.

OpenTelemetry tracing:
- Jaeger UI: http://localhost:16686/search
- Exporter endpoint: set OTEL_EXPORTER_OTLP_ENDPOINT (default: localhost:4318)

## Database Migrations

Run migrations using the CLI:

```ps1
# Dry run (see which migrations would be applied)
go run .\src\main.go migrate-db --dry-run --config .\cfg\config.yaml

# Apply all migrations
go run .\src\main.go migrate-db --config .\cfg\config.yaml

# Force re-create all tables before migrating
go run .\src\main.go migrate-db --force-migrate --config .\cfg\config.yaml
```

## Docker

Build the image:

```ps1
make build-docker
# image tag: go-template
```

Run the container (example):

```ps1
docker run --rm -p 9092:9092 ^
  -v "%cd%\cfg":/app/cfg ^
  -e OTEL_EXPORTER_OTLP_ENDPOINT=host.docker.internal:4318 ^
  go-template serve-http-api --config /app/cfg/config.yaml
```

## Project Structure (high level)

- cfg\config.yaml: application configuration
- docker-compose.yaml: local dependencies (Postgres, MinIO, Jaeger)
- src\cmd: CLI commands (serve-http-api, migrate-db)
- src\core: handlers, middlewares, db, logging, utils
- src\service: business logic layer
- src\otel: OpenTelemetry setup
- src\main.go: entrypoint (Cobra CLI)

## Troubleshooting

- Cannot connect to Postgres: ensure cfg\config.yaml Database.PostgreSQL.DBName matches docker-compose (go-template by default) and containers are up (docker compose ps).
- Ports already in use: change mapped ports in docker-compose.yaml or API.HTTPServerPort in cfg\config.yaml.
- Traces not visible in Jaeger: confirm Jaeger is running and OTEL_EXPORTER_OTLP_ENDPOINT points to the Jaeger collector (localhost:4318 by default).
- On Windows, make run-local uses a Bash script; run it from Git Bash/WSL, or use the Go commands shown above in PowerShell.