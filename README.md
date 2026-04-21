# 🚀 Todo List API

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)](https://go.dev)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-4169E1?style=flat-square&logo=postgresql)](https://postgresql.org)
[![Architecture](https://img.shields.io/badge/Architecture-Clean%20Architecture-FF6B6B?style=flat-square)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)

> **Production-Ready REST API** demonstrating enterprise-grade Go development with Clean Architecture, JWT authentication, real-time monitoring, and comprehensive observability.

## ✨ Why This Project Stands Out

This isn't just another Todo API. It's a **reference implementation** showcasing how to build production-ready Go services with:

- 🏗️ **Clean Architecture** - Proper separation of concerns, testable, maintainable
- 🔐 **Enterprise Security** - JWT auth, bcrypt hashing, input validation
- 📊 **Observability** - In-memory logging with Telegram alerts for errors
- ⚡ **Performance** - Connection pooling, async notifications, optimized queries
- 🧪 **Quality** - Interface-based design, dependency injection, error handling

---

## 🛠️ Tech Stack

| Category | Technology | Purpose |
|----------|------------|---------|
| **Language** | Go 1.22+ | High-performance, concurrent backend |
| **Router** | go-chi/chi | Lightweight, idiomatic HTTP routing |
| **Database** | PostgreSQL 14+ | Reliable, ACID-compliant data store |
| **Driver** | pgx | High-performance PostgreSQL driver |
| **Auth** | golang-jwt/jwt | Industry-standard JWT implementation |
| **Security** | bcrypt | Password hashing (adaptive, slow) |
| **Alerts** | Telegram Bot API | Real-time error notifications |

---

## 🏛️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP REQUEST                             │
└──────────────────────┬──────────────────────────────────────┘
                       │
           ┌───────────▼───────────┐
           │   Handler Layer       │  ← Request validation, JSON serialization
           │   (chi.Router)        │     Response formatting
           └───────────┬───────────┘
                       │
           ┌───────────▼───────────┐
           │   Service Layer       │  ← Business logic, use cases
           │   (Interfaces)        │     Authorization rules
           └───────────┬───────────┘
                       │
           ┌───────────▼───────────┐
           │   Repository Layer    │  ← Data access abstraction
           │   (pgx/PostgreSQL)    │     Transaction management
           └───────────┬───────────┘
                       │
           ┌───────────▼───────────┐
           │   Domain Layer        │  ← Entities, value objects
           │   (Pure Go structs)   │     Domain errors
           └───────────────────────┘
```

### 🎯 Design Patterns Implemented

- **Dependency Injection** - No global state, testable components
- **Repository Pattern** - Data access abstraction, swapable implementations
- **Service Layer** - Business logic isolated from transport concerns
- **Middleware Chain** - Cross-cutting concerns (auth, logging, recovery)
- **Singleton Pattern** - Logger, Notifier with thread-safe initialization

---

## 🚀 Features

### Core API
- ✅ **Full CRUD** - Create, Read, Update, Delete todos
- ✅ **Bulk Operations** - Create multiple todos in one request
- ✅ **User Management** - Registration, login, profile
- ✅ **Data Ownership** - Users only see their own todos (row-level security)

### Security
- 🔐 **JWT Authentication** - Stateless, scalable auth with 24h expiry
- 🔒 **Password Security** - bcrypt hashing with salt (cost: 10)
- 🛡️ **Input Validation** - Domain-level validation with clear error messages
- 🔍 **Authorization** - Middleware validates token on every protected route

### Observability & Monitoring
- 📊 **In-Memory Logging** - View logs via API endpoint
- 🚨 **Telegram Alerts** - Instant notifications for ERROR/WARN logs
- 📈 **Structured Logging** - Context, function name, timestamps
- 🔍 **Health Checks** - `/health` endpoint for monitoring

### Developer Experience
- 📚 **API Documentation** - Complete Postman collection (see ENDPOINTS.md)
- 🧪 **Makefile Commands** - `make run`, `make migrate`, `make test`
- 🔄 **Hot Reload** - Easy development workflow
- 📦 **Migration Scripts** - Version-controlled database changes

---

## 📁 Project Structure

```
todo-list-api-golang/
├── 📂 cmd/api/
│   └── main.go                    # Application entry point, dependency wiring
│
├── 📂 internal/                   # Private application code
│   ├── 📂 config/                 # Environment configuration
│   │   └── config.go
│   ├── 📂 domain/                 # Core business entities
│   │   ├── todo.go               # Todo entity with validation
│   │   ├── user.go               # User entity with password hashing
│   │   └── errors.go             # Domain error types
│   ├── 📂 handler/                # HTTP handlers (transport layer)
│   │   ├── auth_handler.go       # Register, login, logout
│   │   ├── todo_handler.go       # CRUD operations
│   │   └── log_handler.go        # View/clear application logs
│   ├── 📂 middleware/             # HTTP middleware
│   │   └── auth.go               # JWT validation, context injection
│   ├── 📂 repository/             # Data access interfaces
│   │   ├── todo_repository.go    # Todo repository interface
│   │   ├── user_repository.go    # User repository interface
│   │   └── 📂 postgres/           # PostgreSQL implementations
│   │       ├── todo_repository.go
│   │       └── user_repository.go
│   └── 📂 service/                # Business logic
│       ├── auth_service.go       # JWT generation, validation
│       └── todo_service.go       # Todo use cases
│
├── 📂 pkg/                        # Public shared packages
│   ├── 📂 logger/                 # In-memory logging with notifier support
│   │   └── logger.go
│   ├── 📂 notifier/               # Telegram alerting
│   │   └── telegram.go
│   ├── 📂 response/               # Standardized API responses
│   │   └── api_response.go
│   └── 📂 errors/                 # HTTP error mapping
│       └── http_errors.go
│
├── 📂 migrations/                 # SQL migrations
│   ├── 001_create_todos_table.sql
│   ├── 002_create_users_table.sql
│   └── 003_add_user_id_to_todos.sql
│
├── 📂 scripts/
│   └── migrate.sh                # One-command migration runner
│
├── .env.example                  # Configuration template
├── Makefile                      # Common development tasks
├── go.mod                        # Go module definition
└── ENDPOINTS.md                  # Complete API documentation
```

---

## 📊 API Response Format

### Success Response (200 OK, 201 Created)
```json
{
  "success": true,
  "message": "",
  "data": { ... },
  "meta": { "total": 5 },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Error Response (400, 401, 404, 500)
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "title is required"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

## 🚦 Quick Start

### Prerequisites
- Go 1.22+
- PostgreSQL 14+
- Make (optional)

### 1. Clone & Install Dependencies
```bash
git clone https://github.com/yourusername/todo-list-api-golang.git
cd todo-list-api-golang
make deps        # or: go mod tidy
```

### 2. Configure Environment
```bash
cp .env.example .env
# Edit .env with your database credentials
```

### 3. Run Database Migrations
```bash
make migrate     # One command, all migrations
```

### 4. Start the Server
```bash
make run         # Server starts on :8080
```

### 5. Test It
```bash
# Health check
curl http://localhost:8080/health

# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com","password":"secret123"}'
```

---

## 📚 API Endpoints

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/health` | ❌ | Health check |
| `POST` | `/api/v1/auth/register` | ❌ | Register new user |
| `POST` | `/api/v1/auth/login` | ❌ | Login, get JWT |
| `POST` | `/api/v1/auth/logout` | ✅ | Logout |
| `GET` | `/api/v1/todos` | ✅ | List todos (paginated) |
| `POST` | `/api/v1/todos` | ✅ | Create todo |
| `POST` | `/api/v1/todos/bulk` | ✅ | Create multiple |
| `GET` | `/api/v1/todos/:id` | ✅ | Get single todo |
| `PUT` | `/api/v1/todos/:id` | ✅ | Update todo |
| `DELETE` | `/api/v1/todos/:id` | ✅ | Delete todo |
| `GET` | `/api/v1/logs` | ✅ | View logs |
| `DELETE` | `/api/v1/logs` | ✅ | Clear logs |

**Full documentation:** See [ENDPOINTS.md](ENDPOINTS.md) for detailed examples.

---

## 🛡️ Security Highlights

- ✅ **No plaintext secrets** - JWT secret from environment
- ✅ **Password hashing** - bcrypt (adaptive cost factor)
- ✅ **Input sanitization** - Validation at domain layer
- ✅ **SQL injection safe** - Parameterized queries (pgx)
- ✅ **No CORS issues** - Configurable middleware
- ✅ **Request timeouts** - Context-based cancellation

---

## 📈 Performance Features

- **Connection Pooling** - pgx built-in pool (max 10 connections)
- **Async Notifications** - Telegram alerts don't block requests
- **Efficient Queries** - Indexed user_id, completed columns
- **JSON Encoding** - Streaming encoder for large responses
- **Graceful Shutdown** - In-flight requests complete before exit

---

## 🧪 Testing Approach

```bash
# Run all tests
make test

# Run specific package
go test ./internal/service/...

# With coverage
go test -cover ./...
```

**Testability Features:**
- Interface-based repositories (mockable)
- Dependency injection (no globals)
- Domain logic isolated from HTTP/DB

---

## 🚨 Telegram Alerts Setup (Optional)

Get instant notifications when errors occur:

1. Create bot with [@BotFather](https://t.me/botfather)
2. Add bot to your group
3. Configure `.env`:
```bash
TELEGRAM_BOT_TOKEN=your_token_here
TELEGRAM_GROUP_ID=-1001234567890
TELEGRAM_THREAD_ID=305  # Optional
```

**Alert Preview:**
🚨 **ERROR** - AUTH_SERVICE.Login

Time: 2024-01-15 14:30:25  
Message: Failed login attempt for user: john@example.com

---

## 💼 For HR / Hiring Managers

### What This Project Demonstrates

| Skill | Evidence |
|-------|----------|
| **Clean Architecture** | Clear separation: handler → service → repository → domain |
| **Go Proficiency** | Interfaces, goroutines, context, error handling |
| **API Design** | RESTful, consistent responses, proper HTTP codes |
| **Security Awareness** | JWT, bcrypt, input validation, no secrets in code |
| **Database Skills** | PostgreSQL, migrations, indexing, transactions |
| **DevOps Mindset** | Makefile, health checks, graceful shutdown |
| **Observability** | Structured logging, external alerting |
| **Code Quality** | No globals, dependency injection, testable design |

### Production Readiness Checklist

- ✅ Container-ready (12-factor app)
- ✅ Environment-based config
- ✅ Health check endpoint
- ✅ Graceful shutdown
- ✅ Structured logging
- ✅ Error monitoring (Telegram)
- ✅ Database migrations
- ✅ Input validation
- ✅ Authentication & authorization

---

## 📝 Development Commands

```bash
make deps       # Download dependencies
make run        # Start development server
make build      # Build binary
make test       # Run tests
make migrate    # Run database migrations
make clean      # Clean build artifacts
```

---

## 🤝 Contributing

Contributions welcome! Please follow:
1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing`)
5. Open Pull Request

---

## 📄 License

MIT License - see [LICENSE](LICENSE) file

---

## 🙏 Acknowledgments

- [go-chi/chi](https://github.com/go-chi/chi) - Lightweight router
- [golang-jwt/jwt](https://github.com/golang-jwt/jwt) - JWT implementation
- [jackc/pgx](https://github.com/jackc/pgx) - PostgreSQL driver

---

**⭐ Star this repo if you found it helpful!**
