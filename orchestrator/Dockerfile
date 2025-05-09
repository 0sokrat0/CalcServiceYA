FROM golang:1.24 AS builder

WORKDIR /orchestrator

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# RUN go test ./... -v || echo "Tests failed, but skipping for Docker build"

RUN CGO_ENABLED=1 go build -a -installsuffix cgo -o orchestrator ./cmd/orchestrator

FROM alpine:latest
RUN apk add --no-cache libc6-compat sqlite-libs


WORKDIR /root/

COPY --from=builder /orchestrator/orchestrator .
COPY --from=builder /orchestrator/config ./config
RUN chmod +x orchestrator

# Запуск
CMD ["./orchestrator"]