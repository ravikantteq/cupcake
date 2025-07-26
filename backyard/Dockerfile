# Build stage
FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/backyard cmd/main.go

# Run stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/bin/backyard .
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./backyard"]
