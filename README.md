# payments-system

A payments backend in Go. The course builds six microservices from an empty folder and runs them on Kubernetes. I'm building it while following the CodeWise course [Go Microservices Full Course: Build a Payments System From Scratch](https://www.youtube.com/watch?v=AVy0hy4qSnQ). The course has 24 episodes and runs about 20 hours.

This is a personal learning repo. It follows the video as I work through it and has no affiliation with the course author. The author's own code is at [dev-marker/payments](https://github.com/dev-marker/payments).

## Course progress

**Current position: 1:05:29 of 19:47:09 (about 6%).**

That is inside "Append-only & derived balances" in Episode 2, "Designing a Double-Entry Ledger." That section runs from 1:01:27 to 1:08:34. Episode 2 began at 36:29, and its "Accounts" section ran from 52:53 to 1:01:27.

The committed code covers Episode 1 and the start of Episode 2. The ledger migration is a work in progress. I update this line as I move forward, so check it before assuming a feature should exist.

## What the course builds

- A double-entry ledger in PostgreSQL, with type-safe SQL using sqlc and pgx. Locks and idempotency keys control how money moves.
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

Episode 1 covers the finished system, why Go and microservices, local infrastructure, and graceful shutdown. The repo contains the last two, plus the first ledger migration from Episode 2:

| Path | What it is |
| ---- | ---------- |
| `cmd/wallet` | The wallet service. It serves `GET /healthz` on port 8081 and shuts down gracefully on SIGINT or SIGTERM. |
| `docker-compose.yml` | Local infrastructure: Postgres 17, Redis 8, and NATS 2.11 with JetStream enabled. |
| `deploy/postgres/01-create-databases.sql` | Creates the `wallet`, `auth`, and `transfers` databases the first time Postgres starts. Compose mounts the folder into the container's init directory. |
| `db/wallet/migrations` | Wallet database migrations. `000001_init_ledger` creates the `accounts`, `transactions`, and `entries` tables, with triggers that require every transaction to have entries. It is a work in progress and not ready to apply with `migrate`. |
| `docs/commands.md` | Commands I use while working on the project. |

Services started by `docker compose`:

| Service | Image | Ports |
| ------- | ----- | ----- |
| postgres | `postgres:17-alpine` | 5432 |
| redis | `redis:8-alpine` | 6379 |
| nats | `nats:2.11-alpine` | 4222 (client), 8222 (monitoring) |

Postgres has four databases besides the built-in `postgres` one. `payments` is the default, created from `POSTGRES_DB`. `wallet`, `auth`, and `transfers` get one database per service.

## Run it locally

### Prerequisites

| Tool | Version | Used for |
| ---- | ------- | -------- |
| Go | 1.27.4 (see `go.mod`) | Building and running the services |
| Docker with the Compose plugin | Docker 29.x, Compose v5.x | Running Postgres, Redis, and NATS |
| golang-migrate CLI (`migrate`) | v4.20.1, built with `-tags postgres` | Applying and rolling back the database migrations in `db/` |

Install `migrate` with Go:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.20.1
```

This puts the binary in `~/go/bin`. Make sure that directory is on your `PATH`, then check it with `migrate -version`.

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

List the Postgres databases. Expect `auth`, `payments`, `postgres`, `transfers`, and `wallet`:

```bash
docker compose exec -T postgres psql -U payments -c 'SELECT datname FROM pg_database WHERE NOT datistemplate ORDER BY 1;'
```

Stop everything with `docker compose down`. Named volumes, including the Postgres data, survive that command. Add `-v` to delete them.

The init script in `deploy/postgres` runs only when the Postgres volume is first created. If you add a database to that file, run `docker compose down -v` and then `docker compose up -d` to recreate the volume, which deletes the data. To keep the data, create the database in `psql` with `CREATE DATABASE`.

The credentials in `docker-compose.yml` are for local development only.

### Database migrations

Each database change is a pair of files in `db/<service>/migrations`: an `.up.sql` file that applies the change and a `.down.sql` file that reverses it. `migrate` tracks the applied version in a `schema_migrations` table in each database.

Create a new pair. `-seq` numbers the files in order:

```bash
migrate create -ext sql -dir db/wallet/migrations -seq <name>
```

Apply all pending migrations to the `wallet` database, using the local compose credentials:

```bash
migrate -path db/wallet/migrations \
  -database "postgres://payments:payments@localhost:5432/wallet?sslmode=disable" up
```

Roll back the most recent migration with `down 1`, or check the current version with `version`:

```bash
migrate -path db/wallet/migrations -database "<same URL>" down 1
migrate -path db/wallet/migrations -database "<same URL>" version
```

`000001_init_ledger` is still a work in progress, and its down file is empty, so it cannot be rolled back yet. Do not run `up` against a database you care about until the migration is finished.

## Layout

```text
.
├── cmd/
│   └── wallet/            wallet service entry point
├── db/
│   └── wallet/
│       └── migrations/    wallet database migrations
├── deploy/
│   └── postgres/          init scripts, run on first Postgres start
├── docs/
│   └── commands.md        command reference
├── docker-compose.yml     local Postgres, Redis, NATS
├── go.mod
└── README.md
```
