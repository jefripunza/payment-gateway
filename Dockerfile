# syntax=docker/dockerfile:1.7

# ---------- Stage 1: build Vue 3 frontend ----------
FROM node:22-alpine AS fe-builder

WORKDIR /fe

COPY web/package.json web/package-lock.json* ./
RUN npm ci --legacy-peer-deps || npm install --legacy-peer-deps

COPY web/ ./
RUN npm run build

# ---------- Stage 2: build Go backend (embeds web/dist via go:embed) ----------
FROM golang:1.25-alpine AS be-builder

WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
COPY --from=fe-builder /fe/dist ./web/dist

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -trimpath -ldflags "-s -w" -o /out/payment .

# ---------- Stage 3: minimal runner ----------
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S app && adduser -S -G app app

WORKDIR /app
COPY --from=be-builder /out/payment /app/payment

RUN mkdir -p /app/data && chown -R app:app /app && chmod -R 777 /app/data

USER app

ENV PAYMENT_BE=:3005
EXPOSE 3005

VOLUME ["/app/data"]

ENTRYPOINT ["/app/payment"]
