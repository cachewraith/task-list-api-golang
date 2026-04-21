# Docker Setup Guide

This guide explains how to use Docker with the Todo List API project.

## File Structure

```
project-root/
├── Dockerfile              # Multi-stage production build
├── Dockerfile.dev          # Development build with hot reload support
├── docker-compose.yml      # Local/development environment
├── docker-compose.prod.yml # Production environment
├── .dockerignore           # Files to exclude from Docker context
├── .env                    # Local environment (gitignored, exists)
├── .env.example            # Local environment template
├── .env.prod.example       # Production environment template
└── cmd/api/main.go         # Application entry point
```

## Environment Files Setup

### Port Configuration (Dev vs Prod)

| Environment | Default Port | Config File |
|-------------|--------------|-------------|
| Development | 8080 | `.env` |
| Production | 8090 | `.env.prod` |

Change `SERVER_PORT` in your env file if you need different ports.

### Local Development
```bash
cp .env.example .env
# Edit .env - set SERVER_PORT=8080 (or any port you prefer)
```

### Production
```bash
cp .env.prod.example .env.prod
# Edit .env.prod:
# - SERVER_PORT=8090 (different from dev)
# - API_BASE_URL=http://localhost:8090 (change when you get a domain)
```

### When You Get a Domain (Future)

1. Point domain to your server IP
2. Update `.env.prod`:
   ```
   API_BASE_URL=https://api.yourdomain.com
   ```
3. Add reverse proxy (nginx/traefik) in front of the API
4. Remove port mapping or bind to localhost only

## Commands

### Local Development

```bash
# Build and start all services (PostgreSQL + API)
docker compose up -d

# View logs
docker compose logs -f api
docker compose logs -f postgres

# Stop all services
docker compose down

# Stop and remove volumes (⚠️ deletes database data)
docker compose down -v

# Rebuild after code changes
docker compose up -d --build

# Rebuild specific service
docker compose up -d --build api

# Run migrations (SQL files from migrations/ folder)
docker compose exec postgres sh -c 'for f in /tmp/migrations/*.sql; do psql -U $$POSTGRES_USER -d $$POSTGRES_DB -f "$$f"; done'

# Copy migrations into container first, then run:
docker compose cp migrations/ postgres:/tmp/
docker compose exec postgres sh -c 'for f in /tmp/migrations/*.sql; do psql -U $$POSTGRES_USER -d $$POSTGRES_DB -f "$$f"; done'

# Access database
docker compose exec postgres psql -U todo_user -d todo_db
```

### Production

```bash
# Build and start production services
docker compose -f docker-compose.prod.yml --env-file .env.prod up -d

# Run migrations (only needed for fresh database or new migrations)
docker compose -f docker-compose.prod.yml --env-file .env.prod --profile migrate run --rm migrate

# View logs
docker compose -f docker-compose.prod.yml logs -f api

# Stop all services
docker compose -f docker-compose.prod.yml down

# Stop and remove volumes (⚠️ deletes production data)
docker compose -f docker-compose.prod.yml down -v

# Rebuild after code changes
docker compose -f docker-compose.prod.yml --env-file .env.prod up -d --build
```

## Daily Workflow (Recommended: Use Project Names)

Use `-p` flag to keep dev and prod completely isolated with different container names:

### Development
```bash
# Start dev environment
docker compose -p todo-dev up -d

# Stop dev
docker compose -p todo-dev down

# Rebuild dev
docker compose -p todo-dev up -d --build

# Dev logs
docker compose -p todo-dev logs -f api
```

### Production
```bash
# Start prod environment
docker compose -f docker-compose.prod.yml -p todo-prod --env-file .env.prod up -d

# Run migrations (first time or after new SQL files)
docker compose -f docker-compose.prod.yml -p todo-prod --env-file .env.prod --profile migrate run --rm migrate

# Stop prod
docker compose -p todo-prod -f docker-compose.prod.yml down

# Prod logs
docker compose -p todo-prod -f docker-compose.prod.yml logs -f api
```

### View Running Containers
```bash
# See all containers with their project
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
```

## Docker Desktop Usage

### For Local Development:
1. Open Docker Desktop
2. Open terminal in project root
3. Run: `docker compose -p todo-dev up -d`
4. Access API at: http://localhost:8080
5. View container status in Docker Desktop dashboard

### For Production:
1. Create `.env.prod` file with production values (use port 8090)
2. Run: `docker compose -f docker-compose.prod.yml -p todo-prod --env-file .env.prod up -d`
3. Run migrations: `docker compose -f docker-compose.prod.yml -p todo-prod --env-file .env.prod --profile migrate run --rm migrate`
4. Access API at: http://localhost:8090 (or your custom port/domain)

## Key Differences Between Environments

| Feature | Local (docker-compose.yml) | Production (docker-compose.prod.yml) |
|---------|---------------------------|-----------------------------------|
| Image | Based on `Dockerfile.dev` | Based on `Dockerfile` (multi-stage) |
| Source | Mounted as volume | Copied into image |
| Hot reload | Enabled | Disabled |
| Restart | `unless-stopped` | `always` |
| Security | Relaxed | Read-only root, no new privileges |
| Postgres | Exposed on port 5432 | Same (restrict for real prod) |
| Binary size | N/A (uses go run) | Minimal (scratch-based) |

## Volumes

### Local Development
- `postgres_data_dev` - PostgreSQL data persists between restarts
- `go_mod_cache` - Go modules cached for faster rebuilds

### Production
- `postgres_data_prod` - PostgreSQL data (back this up!)

## Troubleshooting

### Container won't start
```bash
# Check logs
docker compose logs api

# Check if PostgreSQL is healthy
docker compose ps
```

### Database connection issues
```bash
# Verify environment variables
docker compose exec api env | grep DB_

# Test PostgreSQL connection from API container
docker compose exec api nc -zv postgres 5432
```

### Rebuild everything from scratch
```bash
# Development
docker compose down -v
docker compose up -d --build

# Production
docker compose -f docker-compose.prod.yml down -v
docker compose -f docker-compose.prod.yml --env-file .env.prod up -d --build
```

## Security Notes for Real Production

1. **Never commit `.env.prod`** to version control
2. **Use strong passwords** for database
3. **Use strong JWT secrets** (min 32 characters, random)
4. **Restrict database port** - don't expose PostgreSQL publicly
5. **Use HTTPS** - put behind reverse proxy (nginx/traefik) with SSL
6. **Regular backups** - backup `postgres_data_prod` volume
7. **Consider using Docker secrets** or external secret manager
