# ─── Stage 1: Build ────────────────────────────────────────────────────────────
FROM golang:1.24-alpine As builder

RUN apk add --no-cache git

WORKDIR /app

# Copy everything (including docs/)
COPY . .

# Build it
RUN go mod tidy && go build -o wallet-service

# ─── Stage 2: Minimal runtime ──────────────────────────────────────────────────
FROM alpine:latest

WORKDIR /app

# Copy binary
COPY --from=builder /app/wallet-service .

# Copy Swagger docs so http‑swagger can serve them
COPY --from=builder /app/docs ./docs

# Copy your env file
COPY .env .

EXPOSE 8080

ENTRYPOINT ["./wallet-service"]
