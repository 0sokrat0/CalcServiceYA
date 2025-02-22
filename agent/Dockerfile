FROM golang:1.24 AS builder

WORKDIR /agent

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go test ./... -v

RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o agent ./cmd/agent/...

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /agent/agent .
COPY --from=builder /agent/config ./config
RUN chmod +x agent
CMD ["./agent"]