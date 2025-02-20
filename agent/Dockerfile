# Этап 1. Сборка и запуск тестов
FROM golang:1.24 AS builder

WORKDIR /agent

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go test ./... -v

# Собираем статически слинкованный бинарник
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o agent ./cmd/agent/...

# Этап 2. Финальный образ
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /agent/agent .
COPY --from=builder /agent/config ./config
RUN chmod +x agent
CMD ["./agent"]