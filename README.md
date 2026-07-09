# Todo API

A RESTful To-Do API built with Go and the Beego v2 framework, featuring JWT-based authentication, PostgreSQL persistence, and a clean Controller → Service → Repository architecture.

## Features

- User registration and login with JWT authentication
- Full CRUD for tasks, scoped per authenticated user
- Filtering by status and pagination on task listing
- Standardized JSON response format across all endpoints
- PostgreSQL with versioned migrations
- Dockerized setup with automatic migration on startup

## Tech Stack

- **Language:** Go 1.26
- **Framework:** Beego v2
- **ORM:** Beego ORM (`beego/beego/v2/client/orm`)
- **Database:** PostgreSQL 17
- **Auth:** JWT (`golang-jwt/jwt/v5`)
- **Password hashing:** bcrypt
- **Migrations:** golang-migrate
- **CLI Tooling:** Bee CLI
- **Containerization:** Docker / Docker Compose

## Project Structure

```
.
├── conf/               # App configuration
├── controllers/        # HTTP handlers
├── database/           # DB init and connection
├── middlewares/        # Auth filter
├── migrations/         # SQL migrations
├── models/             # ORM models
├── repositories/       # Data access layer
├── routers/            # Route definitions
├── services/           # Business logic
└── utils/              # JWT and response helpers
```

## Key Dependencies

| Package                        | Purpose                       |
| ------------------------------ | ----------------------------- |
| `github.com/beego/beego/v2`    | Web framework, router, ORM    |
| `github.com/lib/pq`            | PostgreSQL driver             |
| `github.com/golang-jwt/jwt/v5` | JWT generation and validation |
| `golang.org/x/crypto/bcrypt`   | Password hashing              |

## Prerequisites

- Go 1.26+
- Docker and Docker Compose
- PostgreSQL 17 (if running without Docker)

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/200215-Moynul-Islam/todo-api.git
cd todo-api
```

### 2. Configure environment

Copy the sample config and update values as needed:

```bash
cp conf/app.conf.sample conf/app.conf
```

Set a strong `JWT_SECRET` and adjust `POSTGRES_*` values to match your environment.

### 3. Run with Docker

```bash
docker-compose up --build
```

This starts PostgreSQL, runs pending migrations, and starts the API on `http://localhost:8080`.

### 4. Run locally (without Docker)

Ensure PostgreSQL is running and reachable per your `app.conf`, apply migrations, then:

```bash
go mod download
bee run
```

`bee run` watches for file changes and rebuilds/restarts automatically. Install it first if needed:

```bash
go install github.com/beego/bee/v2@latest
```

## API Endpoints

| Method | Endpoint         | Auth Required | Description                  |
| ------ | ---------------- | ------------- | ---------------------------- |
| GET    | `/health`        | No            | Health check                 |
| POST   | `/auth/register` | No            | Register a new user          |
| POST   | `/auth/login`    | No            | Log in and receive a JWT     |
| GET    | `/tasks`         | Yes           | List tasks (filter/paginate) |
| POST   | `/tasks`         | Yes           | Create a task                |
| GET    | `/tasks/:id`     | Yes           | Get a task by ID             |
| PUT    | `/tasks/:id`     | Yes           | Update a task                |
| DELETE | `/tasks/:id`     | Yes           | Delete a task                |

Authenticated requests require an `Authorization: Bearer <token>` header.

### Sample Requests

**Register**

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe","email":"jane@example.com","password":"secret123"}'
```

**Login**

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"jane@example.com","password":"secret123"}'
```

**Create task**

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy groceries","description":"Milk, eggs, bread"}'
```

**List tasks**

```bash
curl "http://localhost:8080/tasks?status=pending&page=1&limit=10" \
  -H "Authorization: Bearer <token>"
```

**Update task**

```bash
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"status":"done"}'
```

**Delete task**

```bash
curl -X DELETE http://localhost:8080/tasks/1 \
  -H "Authorization: Bearer <token>"
```

## Response Format

All responses follow a consistent structure:

```json
{
  "success": true,
  "message": "Task created successfully",
  "data": {}
}
```

## Environment Variables

| Variable            | Description                    |
| ------------------- | ------------------------------ |
| `POSTGRES_DB`       | Database name                  |
| `POSTGRES_USER`     | Database user                  |
| `POSTGRES_PASSWORD` | Database password              |
| `POSTGRES_HOST`     | Database host                  |
| `POSTGRES_PORT`     | Database port                  |
| `POSTGRES_SSLMODE`  | SSL mode for the DB connection |
| `JWT_SECRET`        | Secret key used to sign JWTs   |

## Migrations

Migrations live in `migrations/` and run automatically via the `migrate` service in `docker-compose.yml`. For manual/local use, install the CLI first:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Then run:

```bash
migrate -path migrations -database "postgres://todo_user:your_password@localhost:5432/todo_db?sslmode=disable" up
```
