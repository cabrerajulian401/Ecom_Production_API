# Ecommerce API

A RESTful API for an ecommerce platform built with Go, following Clean Architecture principles.

## Description

This is a production-grade ecommerce REST API demonstrating modern Go backend development practices. The API provides endpoints for managing products and processing customer orders with full transactional support. Built with Clean Architecture principles, the codebase maintains clear separation of concerns: HTTP handlers act as controllers that parse requests and format responses, services encapsulate business logic and validation, and repositories handle all database operations through type-safe SQL queries. This layered approach ensures the core business logic remains decoupled from external frameworks, making the codebase highly testable, maintainable, and adaptable to changing requirements. The API leverages PostgreSQL for reliable data persistence, sqlc for compile-time verified SQL queries that eliminate runtime errors, and database transactions to guarantee data integrity during complex operations like order placement. Middleware handles cross-cutting concerns including request logging, panic recovery, timeouts, and request tracing—all essential features for production observability and reliability.

## Tech Stack

- **Go 1.25** - Backend language
- **Chi** - Lightweight HTTP router
- **PostgreSQL 16** - Database
- **sqlc** - Type-safe SQL code generation
- **Docker** - Containerized PostgreSQL

## Project Structure

```
├── cmd/                          # Application entry points
│   ├── main.go                   # Main entry point
│   └── api_interface.go          # HTTP server & routing
├── internal/
│   ├── adapters/postgresql/      # Database layer
│   │   ├── migrations/           # SQL migrations
│   │   └── sqlc/                 # Generated query code
│   ├── products/                 # Products domain
│   │   ├── handlers.go           # HTTP handlers
│   │   └── service.go            # Business logic
│   ├── orders/                   # Orders domain
│   │   ├── handlers.go
│   │   ├── service.go
│   │   └── payload.go            # Request/response types
│   ├── writer/                   # HTTP response helpers
│   └── env/                      # Environment config
├── docker-compose.yaml
├── sqlc.yaml
└── go.mod
```

## Getting Started

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- [sqlc](https://sqlc.dev/) (for code generation)
- [goose](https://github.com/pressly/goose) (for migrations)

### Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/cabrerajulian401/ecom.git
   cd ecom
   ```

2. **Start PostgreSQL**
   ```bash
   docker-compose up -d
   ```

3. **Run migrations**
   ```bash
   goose -dir internal/adapters/postgresql/migrations postgres \
     "host=localhost user=postgres password=postgres dbname=ecom sslmode=disable" up
   ```

4. **Generate sqlc code** (if modifying queries)
   ```bash
   sqlc generate
   ```

5. **Run the server**
   ```bash
   go run cmd/*.go
   ```

   Server starts at `http://localhost:8080`

## API Endpoints

| Method | Endpoint    | Description        |
|--------|-------------|--------------------|
| GET    | /health     | Health check       |
| GET    | /products   | List all products  |
| POST   | /order      | Create a new order |

### Examples

**List Products**
```bash
curl http://localhost:8080/products
```

**Create Order**
```bash
curl -X POST http://localhost:8080/order \
  -H "Content-Type: application/json" \
  -d '{
    "customerId": 1,
    "items": [
      {"productId": 1, "quantity": 2},
      {"productId": 2, "quantity": 1}
    ]
  }'
```

## Environment Variables

| Variable       | Default                                                        | Description          |
|----------------|----------------------------------------------------------------|----------------------|
| GOOSE_DBSTRING | host=localhost user=postgres password=postgres dbname=ecom ... | Database connection  |

## Architecture

This project follows **Clean Architecture** principles:

- **Handlers** (Controllers) - Parse HTTP requests, call services
- **Services** (Use Cases) - Business logic, orchestration
- **Repository** (Gateways) - Database access via sqlc

Dependencies flow inward: `Handlers → Services → Repository`

## License

MIT

