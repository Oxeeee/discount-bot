# Stage 1: Build
FROM golang:1.23.4-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make build

# Stage 2: Runtime
FROM debian:buster-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/build/bin/ .

RUN echo 'CONFIG_PATH="./prod.yaml"' > .env

EXPOSE 3000

CMD sh -c 'echo "env: \"prod\"" > prod.yaml && \
    echo "telegram_token: \"${TELEGRAM_TOKEN}\"" >> prod.yaml && \
    echo "database_route: \"data.db\"" >> prod.yaml && \
    exec ./main'
