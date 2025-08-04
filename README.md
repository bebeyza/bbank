
# 🏦 bbank – Backend Service

A Go‑based backend service for a simple banking system. The project includes user management, transaction processing, balance tracking, and an HTTP API layer.

---

## 🚀 Overview

- **Language**: Go (Golang)
- **Database**: PostgreSQL (via GORM)
- **Goal**: Demonstrate clean architecture, concurrency control, secure account operations, and API design in Go.
- **Current status**: Core functionality is implemented ▸ user registration/login, transaction handling, and balance updates.  

---

## 🔧 Project Structure

```
cmd/
  └─ server/                  # Main entry point
internal/
  ├─ config/                 # Env‑based configuration
  ├─ domain/                 # Models (User, Transaction, Balance)
  ├─ repository/             # DB access via GORM
  ├─ service/                # Business logic
  ├─ api/                    # HTTP handlers and router
  └─ worker/                 # Transaction queue & concurrency
pkg/
  └─ logging/                # zerolog or zap integration
```

---

## ✔️ Implemented Components

### 1. Project Setup
- Go modules for dependency management.
- Config loaded from environment variables (PORT, DB_CONN, JWT_SECRET, etc.).
- Logging framework integrated (e.g. `zerolog` or `zap` — note: found conflicts between GORM logger and zap/zerolog which are not yet resolved).
- Graceful shutdown with context cancellation.

### 2. Database & Models
- PostgreSQL schema for `users`, `transactions`, `balances`, `audit_logs`.
- GORM migrations set up using `gorm.io/driver/postgres`.
- Domain structs:
  - `User`, with validation and password hashing.
  - `Transaction`, with state management and rollback.
  - `Balance`, updated in thread-safe manner using `sync.RWMutex`.

### 3. Concurrency & Queue
- Worker pool to process transactions.
- Request queue with Go channels.
- Atomic counters for statistics and batch processing support.

### 4. Core Services
- **UserService**: Register/login, password hashing, role‑based access.
- **TransactionService**: Credit, debit, wallet transfer, rollback support.
- **BalanceService**: Safe updates, historical balance tracking, optimized calculations.

### 5. API & Middleware
- HTTP server with custom router and middleware.
- **Implemented**: Auth endpoints, basic user/transaction/balance API routes.

---

## 🛠️ Local Setup & Running

1. Clone the repo:  
   ```bash
   git clone https://github.com/bebeyza/bbank.git
   cd bbank
   ```

2. Configure environment variables in `.env` (or export manually):
   ```
   PORT=8080
   DB_DSN=postgres://user:pass@localhost:5432/bbank
   JWT_SECRET=your_jwt_secret
   ```

3. Run migrations (if using a script or GORM auto-migrate).  
   Example in Go code:
   ```go
   db.AutoMigrate(&User{}, &Transaction{}, &Balance{}, &AuditLog{})
   ```

4. Start the server:
   ```bash
   go run cmd/server/main.go
   ```

5. Access API on `http://localhost:8080`  
   Test routes:
   - `POST /api/v1/auth/register`
   - `POST /api/v1/auth/login`
   - `POST /api/v1/transactions/credit` etc.

---

## 📦 Docker & DevOps (TODO)

- Multi‑stage `Dockerfile` to build and run the application.
- `docker-compose.yml` with:
  - App service
  - PostgreSQL

---
