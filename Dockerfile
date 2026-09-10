FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o sports-data-platform ./cmd/sports/
RUN go build -o sports-data-platform-batch ./cmd/batch/

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl \
    && rm -rf /var/lib/apt/lists/* \
    && useradd --create-home --shell /bin/bash appuser

WORKDIR /app

COPY --from=builder /app/sports-data-platform .
COPY --from=builder /app/sports-data-platform-batch .

USER appuser

CMD ["./sports-data-platform"]