# Commands

Run these from the repo root. The usual order is: start the stack with Docker, apply migrations with golang-migrate, then inspect Postgres and Redis.

## Docker

Creates and starts the compose project's containers in the background (`-d`), pulling images and creating networks and volumes if they don't exist yet. If they're already running, it leaves them alone. Run it again after `docker compose down` to bring the stack back up.

```bash
docker compose up -d
```

Shows status of this project's compose services (service, state, uptime/health).

```bash
docker compose ps --format "table {{.Service}}\t{{.State}}\t{{.Status}}"
```

Checks the **NATS** monitoring endpoint on port 8222 to see whether the server is healthy, returning `{"status":"ok"}` if it is.

```bash
curl -s localhost:8222/healthz
```

Stops and removes the containers and the default network for this compose project. Named volumes such as Postgres data are kept unless you add `-v` (see below). Doesn't touch the code or images.

```bash
docker compose down
```

Stops and removes the containers and network for this compose project, and also deletes its named volumes. This wipes the Postgres and other stored data, so it can't be undone. Use it only when you want a clean slate.

```bash
docker compose down -v
```

## Golang Migrate

These run the `migrate` CLI installed in WSL. It connects to Postgres through port 5432, which the compose `postgres` service publishes to `localhost`, so the stack must be up first.

Applies every migration in `db/wallet/migrations` that hasn't run yet against the `wallet` database, in version order.

```bash
migrate -path db/wallet/migrations -database "postgres://payments:payments@localhost:5432/wallet?sslmode=disable" up
```

Rolls back every applied migration by running the `.down.sql` files in reverse order. It asks for confirmation first, because it removes the tables and their data.

```bash
migrate -path db/wallet/migrations -database "postgres://payments:payments@localhost:5432/wallet?sslmode=disable" down
```

Sets the recorded migration version to 0 and clears the `dirty` flag without running any SQL. Use it when a failed migration has left the database marked dirty and `up`/`down` refuse to run. Clean up any half-applied changes by hand first, because `force` only changes the bookkeeping.

```bash
migrate -path db/wallet/migrations -database "postgres://payments:payments@localhost:5432/wallet?sslmode=disable" force 0
```

## Postgres

These run inside the compose `postgres` container. `-T` turns off TTY allocation, so the output prints cleanly and the command also works in scripts.

Checks that the **Postgres server** is accepting connections.

```bash
docker compose exec postgres pg_isready -U payments
```

Lists the tables in the `wallet` database.

```bash
docker compose exec -T postgres psql -U payments -d wallet -c '\dt'
```

Lists the views in the `wallet` database.

```bash
docker compose exec -T postgres psql -U payments -d wallet -c '\dv'
```

Shows golang-migrate's bookkeeping table: the current migration version, and whether it's `dirty` from a failed run.

```bash
docker compose exec -T postgres psql -U payments -d wallet -c 'TABLE schema_migrations'
```

Shows every row in the `transactions` table. `TABLE x` is shorthand for `SELECT * FROM x`.

```bash
docker compose exec -T postgres psql -U payments -d wallet -c 'TABLE transactions'
```

## Redis

These run inside the compose `redis` container.

Checks that the **Redis server** responds, returning `PONG` if it's up.

```bash
docker compose exec redis redis-cli ping
```

## Docker quick reference

Every command that goes through Docker, in the order you'd usually use them. The migrate commands here use the `migrate/migrate` image instead of the WSL install, so you don't need `migrate` installed locally. They join the compose network (`payments-system_default`) and reach Postgres by its service name, `postgres`, instead of `localhost`.

```bash
# Start the stack and check it
docker compose up -d
docker compose ps --format "table {{.Service}}\t{{.State}}\t{{.Status}}"
docker compose exec redis redis-cli ping
docker compose exec postgres pg_isready -U payments

# Migrations (Docker alternative to the local migrate CLI)
docker run --rm -v "$PWD/db/wallet/migrations:/migrations" --network payments-system_default \
  migrate/migrate:v4.18.1 -path /migrations \
  -database "postgres://payments:payments@postgres:5432/wallet?sslmode=disable" up
docker run --rm -it -v "$PWD/db/wallet/migrations:/migrations" --network payments-system_default \
  migrate/migrate:v4.18.1 -path /migrations \
  -database "postgres://payments:payments@postgres:5432/wallet?sslmode=disable" down
docker run --rm -v "$PWD/db/wallet/migrations:/migrations" --network payments-system_default \
  migrate/migrate:v4.18.1 -path /migrations \
  -database "postgres://payments:payments@postgres:5432/wallet?sslmode=disable" force 0

# Inspect the wallet database
docker compose exec -T postgres psql -U payments -d wallet -c '\dt'
docker compose exec -T postgres psql -U payments -d wallet -c '\dv'
docker compose exec -T postgres psql -U payments -d wallet -c 'TABLE schema_migrations'
docker compose exec -T postgres psql -U payments -d wallet -c 'TABLE transactions'

# Stop the stack (add -v to also wipe the data)
docker compose down
```

The `down` migration uses `-it` because it waits for you to confirm. The NATS health check is left out because it's a plain `curl` to the published port, not a Docker command.
