FROM golang:1.24.2 AS builder
WORKDIR /auth-service

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o grpc-auth ./cmd/main.go

FROM alpine:3.21
WORKDIR /auth-service

ENV CONFIG_PATH=config/config.yaml

COPY --from=builder /auth-service/migrations /root/migrations
COPY --from=builder /auth-service/config/config.yaml ${CONFIG_PATH}
COPY --from=builder /auth-service/grpc-auth .
CMD ["./grpc-auth"]