# Build stage
FROM golang:1.20-bullseye AS builder

# Install required dependencies for librdkafka
RUN apt-get update && apt-get install -y \
    librdkafka-dev \
    pkg-config \
    build-essential \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bin/backyard cmd/main.go

# Run stage
FROM debian:bullseye-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    librdkafka1 \
    curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /root/

COPY --from=builder /app/bin/backyard .
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./backyard"]
