FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o sports-data-platform ./cmd/sports/

FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/sports-data-platform .
CMD ["./sports-data-platform"]