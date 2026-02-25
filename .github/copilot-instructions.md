# Copilot Instructions for SIMKEU Codebase

## Project Overview
SIMKEU is a microservices-based financial management system with:
- **Go services** (auth, blockchain, debitur, log, master, payment, piutang, realisasi, tagihan)
- **Python FastAPI** for realisasi service
- **Next.js frontend**
- PostgreSQL databases per service
- Docker Compose for local development, Kubernetes for deployment

## Architecture Patterns

### Layered Service Structure
Each Go service follows a consistent structure:
```
service/
  ├── cmd/main.go              # Entry point
  ├── internal/
  │   ├── database/db.go       # PostgreSQL connection (retry logic)
  │   ├── handler/             # HTTP handlers (Gin framework)
  │   ├── service/             # Business logic
  │   ├── repository/          # Data access layer
  │   └── middleware/          # Auth/logging middleware
  ├── Dockerfile               # Multi-stage Go build
  ├── docker-compose.yml       # Service + PostgreSQL
  └── go.mod
```

**Key Pattern**: Strict separation between HTTP handlers (Gin), business logic (service layer), and database access (repository pattern with `*sql.DB`).

### Database & Connection
- PostgreSQL 15 with connection retry logic (10 attempts, 3-second intervals)
- Connection string built from env vars: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- Each service has its own database in docker-compose (e.g., `simkeu_auth`)
- Raw SQL queries with parameterized inputs (`$1`, `$2` notation)

### Authentication
The auth service implements:
- User registration/login with bcrypt password hashing
- JWT token generation (HS256) with 24-hour expiry
- JWT secret via `JWT_SECRET` environment variable
- Claims include `user_id`, `email`, `exp`

Expected pattern for other services: decode JWT from `Authorization` header in middleware, extract user_id for request context.

## Development Workflows

### Local Development with Docker Compose
```bash
cd services/<service-name>
docker-compose up --build
```
Service runs on port 8080, PostgreSQL on 5432. Database credentials are hardcoded in compose (dev-only).

### Building Individual Services
```bash
cd services/<service-name>
go build -o app ./cmd
```

### Environment-Specific Configuration
- **Local**: docker-compose.yml hardcoded values
- **Staging/Production**: Environment variables must be set before deployment
- Critical vars: `DB_*`, `JWT_SECRET`, `PORT` (default 8080)

## Code Conventions

### Go Service Style
- Import auth module as `simkeu/service-auth` (see go.mod)
- Gin request binding: `ShouldBindJSON()` for validation
- Error responses: `gin.H{"error": "message"}` with appropriate HTTP status
- Success responses: `gin.H{"message": "..."}` or data JSON

### Handler Pattern (from auth_service example)
```go
func (h *AuthHandler) Register(c *gin.Context) {
  var input struct {
    Email string `json:"email"`
    Password string `json:"password"`
  }
  if err := c.ShouldBindJSON(&input); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
    return
  }
  // Process with service layer
}
```

### Repository Pattern
- Methods receive `*sql.DB` from package receiver
- Named parameterized queries with `$1, $2` syntax (PostgreSQL)
- Return multiple values for partial queries: `(id int, value string, err error)`
- Example: `r.DB.QueryRow().Scan(&id, &password)`

### Service/Business Logic
- Encrypt sensitive data (bcrypt for passwords)
- Generate tokens with expiry
- Validate business rules before repository calls
- Return meaningful errors that handlers can translate to HTTP status codes

## Integration Points & Dependencies

### Microservice Communication
- Services communicate via HTTP (likely RESTful APIs)
- Auth token validation required across services
- Each service owns its PostgreSQL database (no shared databases)

### External Dependencies
- `github.com/gin-gonic/gin` - HTTP framework
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/golang-jwt/jwt/v5` - JWT token handling
- `golang.org/x/crypto/bcrypt` - Password hashing

### Deployment
- Dockerfile uses `golang:1.25-alpine` base image
- Build: `go mod tidy` + `go build -o app ./cmd`
- Expose port 8080
- k8s/ directory contains Kubernetes manifests (check for ConfigMaps, Secrets for env vars)

## Common Tasks

### Adding a New Endpoint
1. Create handler in `internal/handler/service_handler.go`
2. Call existing service methods or create new ones in `internal/service/service.go`
3. Use repository for DB queries in `internal/repository/`
4. Register route in `cmd/main.go` (Gin router)
5. Update docker-compose environment if new secrets needed

### Adding Database Migrations
- No migration tool configured yet; SQL schemas assumed to exist
- When creating new tables, document DDL and manual setup for k8s environments

### Debugging Service Issues
- Docker Compose logs: `docker-compose logs -f auth`
- Database retry logic waits 30 seconds total (10 × 3s) before failing
- Check env var names match exactly (case-sensitive in Linux)

## Red Flags & Common Mistakes
- Don't hardcode secrets in code; use environment variables
- Database connection strings must use `sslmode=disable` for local dev
- JWT expiry is in Unix timestamp seconds (not milliseconds)
- Repository methods should not know about HTTP status codes (keep concerns separated)
- Each service's go.mod module name must be unique (e.g., `simkeu/service-auth`, not generic)
