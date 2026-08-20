# Payment Gateway

Multi-provider payment gateway credential management dashboard — deploy once, add as many payment gateways as you need from one place.

## Features

| Area | Details |
|------|---------|
| Multi-provider | Store credentials for Midtrans, Xendit, Tripay, Duitku, PayPal, Stripe and more |
| Credential security | AES-256-GCM encryption at rest — secrets are never returned in plaintext |
| Wallets | Track balances across multiple currencies |
| Team access | Role-based users (admin / viewer) with JWT auth |
| Single binary | Vue 3 frontend embedded into the Go binary — one container, one process |

## Stack

- **Backend**: Go + Fiber v3, GORM + SQLite, JWT (HS256), bcrypt, AES-256-GCM
- **Frontend**: Vue 3 + TypeScript, Vite, Tailwind CSS v4, Pinia, Vue Router, ky, lucide-vue-next
- **Deploy**: Multi-stage Dockerfile → single static binary

## Getting started

```bash
# backend
cp .env.example .env
go build -o run .
./run          # listens on :3005

# frontend (dev mode with hot reload)
cd web
npm install
npm run dev    # Vite on :5173, proxies /api to :3005
```

Default admin is seeded on first start (see console log). Change it immediately after first login.

## API

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/login` | Sign in, returns JWT |
| GET | `/api/v1/auth/me` | Current user |
| PATCH | `/api/v1/auth/password` | Change own password |
| GET/POST | `/api/v1/users` | List / create users |
| PATCH/DELETE | `/api/v1/users/:id` | Update / delete user |
| GET/POST | `/api/v1/wallets` | List / create wallets |
| PATCH/DELETE | `/api/v1/wallets/:id` | Update / delete wallet |
| GET/POST | `/api/v1/providers` | List / create providers |
| PATCH/DELETE | `/api/v1/providers/:id` | Update / delete provider |
| GET | `/api/v1/health` | Health check |

Provider credentials are encrypted before storage; list responses return masked values only.

## Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PAYMENT_BE` | `:3005` | HTTP listen address |
| `PAYMENT_DB_PATH` | `data/payment.db` | SQLite file location |
| `PAYMENT_JWT_SECRET` | dev fallback | JWT signing secret — set a strong value in production |
| `PAYMENT_ENC_KEY` | dev fallback | 32-byte AES key for credential encryption — set in production |

## Project layout

```
├── main.go          # bootstrap, DB, migration, seed
├── auth.go          # JWT, login, middleware, routes
├── models.go        # User / Wallet / Provider (UUID v7 PKs)
├── resources.go     # users + wallets handlers
├── providers.go     # provider handlers + masking
├── crypto.go        # AES-256-GCM encrypt/decrypt
├── spa.go           # embedded SPA serving
└── web/             # Vue 3 frontend (Vite + Tailwind v4)
```
