FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o sports-data-platform ./cmd/sports/

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*
    
WORKDIR /app
COPY --from=builder /app/sports-data-platform .
CMD ["./sports-data-platform"]