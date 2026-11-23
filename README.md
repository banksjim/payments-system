# payments-system

A payments backend in Go, built from an empty folder to six microservices running on Kubernetes. I'm building it while following the CodeWise course [Go Microservices Full Course: Build a Payments System From Scratch](https://www.youtube.com/watch?v=AVy0hy4qSnQ). The course has 24 episodes and runs about 20 hours.

This is a personal learning repo. It follows the video as I work through it and has no affiliation with the course author. The author's own code is at [dev-marker/payments](https://github.com/dev-marker/payments).

## Course progress

**Current position: 36:28 of 19:47:09 (about 3%).**

That is the end of Episode 1, in the "Graceful shutdown" section (27:32 to 36:29). Episode 2, "Designing a Double-Entry Ledger," starts at 36:29.

The code in this repo matches the video up to that point and nothing beyond it. I update this line as I move forward, so check it before assuming a feature should exist.

## What the course builds

- A double-entry ledger in PostgreSQL, with type-safe SQL using sqlc and pgx. Money moves safely using locks and idempotency keys.
- JWT auth, gRPC between services, an API gateway, and events with NATS JetStream and the transactional outbox.
- Transfers as sagas, notifications with a dead-letter queue, and observability with OpenTelemetry, Prometheus, and Grafana.
- Testing with testcontainers and k6, distroless images, CI, signed webhooks, multi-currency FX, fraud checks, and reconciliation.
- Security with mTLS and least-privilege roles, performance tuning, Kubernetes with Helm, and chaos engineering.

### Episodes

| # | Episode | Starts at |
| - | ------- | --------- |
| 1 | Build a Production-Grade Payments Backend | 0:00 |
| 2 | Designing a Double-Entry Ledger | 36:29 |
| 3 | Type-Safe SQL with sqlc + pgx | 1:17:49 |
| 4 | Config, Logging & Errors | 2:05:17 |
| 5 | The Wallet API | 2:55:15 |
| 6 | Moving Money Safely | 3:41:13 |
| 7 | Idempotency Keys | 4:27:23 |
| 8 | The Auth Service | 5:15:08 |
| 9 | gRPC Between Services | 6:08:47 |
| 10 | The API Gateway | 7:01:26 |
| 11 | Events with NATS JetStream | 7:54:50 |
| 12 | Transfer Service & Sagas | 8:45:12 |
| 13 | Notification Service | 9:43:02 |
| 14 | Observability | 10:31:32 |
| 15 | Testing Like a Pro | 11:25:37 |
| 16 | Shipping It | 12:13:44 |
| 17 | Outbound Webhooks | 13:00:54 |
| 18 | Multi-currency & FX | 13:53:37 |
| 19 | Fraud & Risk Checks | 14:45:53 |
| 20 | Reconciliation & Reporting | 15:39:43 |
| 21 | Security Hardening | 16:30:03 |
| 22 | Performance | 17:24:50 |
| 23 | Kubernetes | 18:13:52 |
| 24 | Chaos & the Final Demo | 18:59:46 |

## What is in the repo now

Episode 1 covers the finished system, why Go and microservices, local infrastructure, and graceful shutdown. The repo contains the last two:

| Path | What it is |
| ---- | ---------- |
| `cmd/wallet` | The wallet service. It serves `GET /healthz` on port 8081 and shuts down gracefully on SIGINT or SIGTERM. |
| `docker-compose.yml` | Local infrastructure: Postgres 17, Redis 8, and NATS 2.11 with JetStream enabled. |
| `docs/commands.md` | Commands I use while working on the project. |

Services started by `docker compose`:

| Service | Image | Ports |
| ------- | ----- | ----- |
| postgres | `postgres:17-alpine` | 5432 |
| redis | `redis:8-alpine` | 6379 |
| nats | `nats:2.11-alpine` | 4222 (client), 8222 (monitoring) |

## Run it locally

You need Go (see `go.mod` for the version) and Docker with the compose plugin.

```bash
# Start Postgres, Redis, and NATS in the background
docker compose up -d

# Confirm all three are running (Postgres reports healthy)
docker compose ps

# Run the wallet service
go run ./cmd/wallet

# In another terminal, check it
curl -s localhost:8081/healthz
# {"status":"ok"}
```

Check the infrastructure directly:

```bash
docker compose exec postgres pg_isready -U payments
docker compose exec redis redis-cli ping
curl -s localhost:8222/healthz
```

Stop everything with `docker compose down`. Named volumes, including the Postgres data, survive that command. Add `-v` to delete them.

The credentials in `docker-compose.yml` are for local development only.

## Layout

```text
.
├── cmd/
│   └── wallet/        wallet service entry point
├── docs/
│   └── commands.md    command reference
├── docker-compose.yml local Postgres, Redis, NATS
└── go.mod
```
