# Kafgres

Minimal microservice that cyclically reads rows from PostgreSQL and publishes them to Kafka. This program is used in interviews to solve practical troubleshooting tasks. Directory deploy/ at the interview differs from the current one and contains errors that need to be fixed.

## How it works

1. Configuration is loaded from environment variables.
2. Connections to PostgreSQL and Kafka are established.
3. A periodic worker runs which:
   - reads up to 100 rows from the `test_data` table;
   - writes the data to Kafka as JSON;
   - updates health state and logs success/errors.
4. HTTP server exposes `/health`.

Key implementation points:
- Service bootstrap: `internal/app/kafgres.go`.
- Read/write cycle and logging: `internal/pkg/worker/worker.go`.
- Health endpoint: `internal/pkg/health/health.go`.

## Local run

`deploy/docker-compose.yml` starts Postgres, Kafka (KRaft), init jobs, and the service.

```bash
docker-compose -f deploy/docker-compose.yml up -d
```

## Configuration

All settings are environment-based (defaults exist):

| Variable | Default | Description |
| --- | --- | --- |
| `POSTGRES_HOST` | `localhost` | PostgreSQL host |
| `POSTGRES_PORT` | `5432` | PostgreSQL port |
| `POSTGRES_USER` | `postgres` | PostgreSQL user |
| `POSTGRES_PASSWORD` | `password` | PostgreSQL password |
| `POSTGRES_DB` | `postgres` | PostgreSQL database |
| `POSTGRES_TABLE` | `test_data` | PostgreSQL table to read from |
| `KAFKA_BROKERS` | `localhost:9092` | Kafka brokers (comma-separated) |
| `KAFKA_TOPIC` | `test-topic` | Kafka topic |
| `HTTP_PORT` | `8080` | HTTP server port |
| `DEFAULT_POLL_INTERVAL` | `5s` | Poll interval (Go duration, e.g. `10s`, `1m`) |

### Health check

```bash
curl http://localhost:8080/health
```

Expected responses:
- `success` — last read/write cycle succeeded;
- `failure` — read/write failed or database returned no rows.

### Kafka check

Consume messages from Kafka:

```bash
docker compose -f deploy/docker-compose.yml exec kafka kafka-console-consumer \
  --bootstrap-server kafka:29092 \
  --topic test-topic \
  --from-beginning
```

### Database schema

Expected table `test_data`:

```sql
CREATE TABLE test_data (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL
);
```
