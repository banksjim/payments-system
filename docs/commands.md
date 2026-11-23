# Commands

## Docker

Creates and starts the compose project's containers in the background (`-d`), pulling images and creating networks and volumes if they don't exist yet. If they're already running, it leaves them alone. Run it again after `docker compose down` to bring the stack back up.

```bash
docker compose up -d
```

Shows status of this project's compose services (service, state, uptime/health).

```bash  
docker compose ps --format "table {{.Service}}\t{{.State}}\t{{.Status}}"  
```

Checks that the **NATS** monitoring endpoint on port 8222 whether the server is healthy, returning `{"status":"ok"}` if it is.

```bash
curl -s localhost:8222/healthz
```

Checks that the **Postgres server** in the compose `postgres` service is accepting connections.

```bash
docker compose exec postgres pg_isready -U payments
```

Checks that the **Redis server** in the compose redis service responds, returning `PONG` if it's up.

```bash
docker compose exec redis redis-cli ping
```

Stops and removes the containers and the default network for this compose project. Named volumes such as Postgres data are kept unless you add `-v` (see below). Doesn't touch the code or images.

```bash
docker compose down
```

Stops and removes the containers and network for this compose project, and also deletes its named volumes. This wipes the Postgres and other stored data, so it can't be undone. Use it only when you want a clean slate.

```bash
docker compose down -v
```
