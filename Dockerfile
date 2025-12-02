FROM golang:1.25.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server .

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*
    
COPY --from=builder /app/server .

EXPOSE 3000

CMD ["./server"]



