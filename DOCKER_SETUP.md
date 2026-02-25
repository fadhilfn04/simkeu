# SIMKEU - Multi-Service Docker Compose Setup

This is the unified docker-compose configuration for running all SIMKEU microservices with a single shared PostgreSQL database.

## Architecture

- **1 PostgreSQL Container** with 9 databases (one for each service)
- **9 Go Microservices** (auth, blockchain, debitur, log, master, payment, piutang, realisasi, tagihan)
- **Shared Network** (simkeu-network) for inter-service communication
- **Health checks** to ensure database is ready before services start

## Services and Ports

| Service | External Port | Internal Port | Database |
|---------|---------------|---------------|----------|
| Auth | 8080 | 8080 | simkeu_auth |
| Blockchain | 8081 | 8080 | simkeu_blockchain |
| Debitur | 8082 | 8080 | simkeu_debitur |
| Log | 8083 | 8080 | simkeu_log |
| Master | 8084 | 8080 | simkeu_master |
| Payment | 8085 | 8080 | simkeu_payment |
| Piutang | 8086 | 8080 | simkeu_piutang |
| Realisasi | 8087 | 8080 | simkeu_realisasi |
| Tagihan | 8088 | 8080 | simkeu_tagihan |

## Getting Started

### Start All Services

```bash
# From the root directory (/opt/simkeu)
docker-compose up --build
```

### Stop All Services

```bash
docker-compose down
```

### Remove All Data

```bash
docker-compose down -v
```

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f auth

# Specific service (last 100 lines)
docker-compose logs -f --tail=100 blockchain
```

## Database Configuration

- **Host**: `db` (from within containers)
- **Host**: `localhost` (from your machine)
- **Port**: 5432
- **Username**: simkeu
- **Password**: simkeu123

### Database Names

- `simkeu_auth`
- `simkeu_blockchain`
- `simkeu_debitur`
- `simkeu_log`
- `simkeu_master`
- `simkeu_payment`
- `simkeu_piutang`
- `simkeu_realisasi`
- `simkeu_tagihan`

### Connect to PostgreSQL from Host

```bash
psql -h localhost -U simkeu -d simkeu_auth
# Enter password: simkeu123
```

## Service Endpoints

### Health Checks

```bash
curl http://localhost:8080/health   # Auth
curl http://localhost:8081/health   # Blockchain
curl http://localhost:8082/health   # Debitur
# ... and so on
```

### Auth Service Examples

```bash
# Register
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

## Network Communication

Services can communicate with each other using the service name as hostname:

```
http://auth:8080/login
http://blockchain:8080/api/status
```

## Troubleshooting

### Services keep restarting

Check if the database is healthy:
```bash
docker-compose logs db
```

### Connection refused errors

Ensure all services wait for the database health check:
```bash
docker-compose logs auth
```

### Database initialization failed

The `init-databases.sh` script creates all databases on first run. If it fails:

1. Stop containers: `docker-compose down`
2. Remove volume: `docker volume rm simkeu_postgres_data`
3. Start again: `docker-compose up --build`

## Individual Service Development

If you need to run a single service separately, each service has its own `docker-compose.yml` in its directory. However, for the unified setup, use the root-level `docker-compose.yml`.
