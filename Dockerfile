# Build stage
FROM golang:1.23-alpine3.20 AS builder

# Install build dependencies
RUN apk update && apk add --no-cache --virtual .build-deps \
    gcc \
    musl-dev \
    git \
    ca-certificates \
    && rm -rf /var/cache/apk/*

WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
WORKDIR /app/src
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOCACHE=/root/.cache/go-build

RUN --mount=type=cache,target="/root/.cache/go-build" \
    go build \
    -ldflags="-w -s -X go-template/src/version.GitCommit=$(git rev-parse HEAD) -X go-template/src/version.Version=$(git describe --tags --always)" \
    -o /build/app

# Final stage
FROM alpine:3.20
LABEL maintainer="yotsaphonliu <yotsaphonliu@gmail.com>" \
      description="Go Template Application" \
      version="1.0"

# Install runtime dependencies and create non-root user
RUN apk update && apk add --no-cache \
    bash \
    tzdata \
    postgresql-client \
    ca-certificates \
    && rm -rf /var/cache/apk/*

# Set security configurations
RUN umask 027 && echo "umask 0027" >> /etc/profile

# Copy binary from builder stage
COPY --from=builder /build/app /usr/local/bin/app

# Create app directory and set ownership
WORKDIR /app

# Set environment variables
ENV TZ="Asia/Bangkok"

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD app --health-check || exit 1

ENTRYPOINT ["app"]
