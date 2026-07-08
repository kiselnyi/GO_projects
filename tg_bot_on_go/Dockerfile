# === Stage 1: Сборка бинарника ===
FROM golang:1.26-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /duty-bot ./cmd/bot

# === Stage 2: Финальный минимальный образ ===
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /duty-bot .

CMD ["./duty-bot"]