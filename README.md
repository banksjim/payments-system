# payments-system

A payments backend in Go. The course builds six microservices from an empty folder and runs them on Kubernetes. I'm building it while following the CodeWise course [Go Microservices Full Course: Build a Payments System From Scratch](https://www.youtube.com/watch?v=AVy0hy4qSnQ). The course has 24 episodes and runs about 20 hours.

This is a personal learning repo. It follows the video as I work through it and has no affiliation with the course author. The author's own code is at [dev-marker/payments](https://github.com/dev-marker/payments).

<!-- Source: repo About text, chapters.md. Open question: the About text names the course "Go payments system development course" and links a youtu.be URL. This intro uses the video title and a watch?v= URL. Confirm which to keep. -->

## Course progress

**Current position: 1:17:49 of 19:47:09 (about 7%).**

I've just started Episode 3, "Type-Safe SQL with sqlc + pgx." I finished Episode 2, "Designing a Double-Entry Ledger," and its last section, "Trying to break it," ran from 1:08:34 to 1:17:49. The next section is "Why sqlc," at 1:19:18.

For Episode 2, I've committed the ledger schema as a work-in-progress migration, with append-only triggers and the `account_balances` view, plus scratch SQL tests. I haven't committed anything from Episode 3 yet, including the Go ledger code.

<!-- Updated by temp-update-readme. Check the status paragraph against git before editing. -->

## What the course builds

<!-- Source: CodeWise course description. Only change when the course description changes. -->

- A double-entry ledger in PostgreSQL, with type-safe SQL using sqlc and pgx. Locks and idempotency keys control how money moves.
- JWT auth, gRPC between services, an API gateway, and events with NATS JetStream and the transactional outbox.
- Transfers as sagas, notifications with a dead-letter queue, and observability with OpenTelemetry, Prometheus, and Grafana.
- Testing with testcontainers and k6, distroless images, CI, signed webhooks, multi-currency FX, fraud checks, and reconciliation.
- Security with mTLS and least-privilege roles, performance tuning, Kubernetes with Helm, and chaos engineering.

### Episodes

<!-- Source: chapters.md. Regenerate when the course chapter list changes. -->

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

<!-- Source: git ls-files. Add a row when a tracked file is added. -->

Episode 1 covers the finished system, why Go and microservices, local infrastructure, and graceful shutdown. The repo has code for the last two of those topics: local infrastructure and graceful shutdown. It also has the first ledger migration from Episode 2:

| Path | What it is |
| ---- | ---------- |
| `cmd/wallet` | The wallet service. It serves `GET /healthz` on port 8081 and shuts down gracefully on SIGINT or SIGTERM. |
| `docker-compose.yml` | Local infrastructure: Postgres 17, Redis 8, and NATS 2.11 with JetStream enabled. |
| `deploy/postgres/01-create-databases.sql` | Creates the `wallet`, `auth`, and `transfers` databases the first time Postgres starts. Compose mounts the folder into the container's init directory. |
| `db/wallet/migrations` | Wallet database migrations. `000001_init_ledger` creates the `accounts`, `transactions`, and `entries` tables, with triggers that require every transaction to have entries. It is a work in progress and not ready to apply with `migrate`. |
| `docs/commands.md` | Command reference for Docker, golang-migrate, Postgres, and Redis. |
| `Makefile` | Shortcuts for the migration commands: `make migrate-up`, `make migrate-down`, `make migrate-version`, and `make migrate-new name=<name>`. |
| `sql_tests/` | Scratch SQL for the ledger schema. `scratch.sql` sets up accounts and a balance query. `scratch2.sql` writes transactions, including writes the schema should reject. |
| `.gitignore` | Files Git does not track, including the local VS Code workspace file. |

<!-- Source: docker-compose.yml. Regenerate when services or ports change. -->

Services started by `docker compose`:

| Service | Image | Ports |
| ------- | ----- | ----- |
| postgres | `postgres:17-alpine` | 5432 |
| redis | `redis:8-alpine` | 6379 |
| nats | `nats:2.11-alpine` | 4222 (client), 8222 (monitoring) |

<!-- Source: docker-compose.yml POSTGRES_DB, deploy/postgres/*.sql. Update when init scripts change. -->

Postgres has four databases besides the built-in `postgres` one. `payments` is the default, created from `POSTGRES_DB`. `wallet`, `auth`, and `transfers` get one database per service.

## Run it locally

<!-- Source: go.mod, docker version, docker compose version, go version -m migrate. -->

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

<!-- Source: docker-compose.yml credentials. -->

The credentials in `docker-compose.yml` are for local development only.

### Database migrations

<!-- Source: Makefile migrate targets, db/*/migrations state. Update when targets or migration state change. -->

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

The `Makefile` wraps these commands. `make migrate-up` applies pending migrations, `make migrate-down` rolls back one step, `make migrate-version` shows the current version, and `make migrate-new name=<name>` creates a new pair. Set `SERVICE` to target another database, for example `make migrate-up SERVICE=auth`.

`000001_init_ledger` is still a work in progress. Its down file is written but has not been tested, so don't rely on rolling it back yet. Don't run `up` against a database you care about until the migration is finished.

## Layout

<!-- Source: git ls-files. Regenerate when files change. -->

```text
.
├── .gitignore
├── Makefile               migration shortcuts
├── cmd/
│   └── wallet/            wallet service entry point
├── db/
│   └── wallet/
│       └── migrations/    wallet database migrations
├── deploy/
│   └── postgres/          init scripts, run on first Postgres start
├── docs/
│   └── commands.md        command reference
├── sql_tests/
│   ├── scratch.sql        ledger setup and balance query
│   └── scratch2.sql       ledger writes, including rejected ones
├── docker-compose.yml     local Postgres, Redis, NATS
├── go.mod
└── README.md
```
