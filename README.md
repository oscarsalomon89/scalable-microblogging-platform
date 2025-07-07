# Scalable Microblogging Platform

A scalable, modular microblogging application built with Go, following the Hexagonal (Ports & Adapters) architecture. The project is designed for maintainability, testability, and easy extension, using PostgreSQL for persistence and Redis for caching.

---

## Features

- User creation
- Tweet creation and timeline retrieval
- Modular, clean architecture (Hexagonal)
- RESTful API
- Database migrations managed via [golang-migrate](https://github.com/golang-migrate/migrate)
- Dockerized for easy local development
- Environment-based configuration

---

## Tech Stack

- **Backend:** Go (Golang)
- **Database:** PostgreSQL
- **Cache:** Redis
- **Migrations:** migrate/migrate
- **Containerization:** Docker & Docker Compose

---

## Architecture

This project follows Hexagonal Architecture (Ports & Adapters) combined with package-oriented design, leveraging idiomatic Go practices to create clear, maintainable, and cohesive packages. The business logic is decoupled from infrastructure and delivery mechanisms, promoting modularity and testability. Key modules include:

- **cmd/api/**: Application entry point and HTTP server
- **internal/adapters/http/**: HTTP handlers for users and tweets
- **internal/adapters/postgres/**: PostgreSQL repositories
- **internal/adapters/redis/**: Redis repositories
- **internal/platform/pg/**: Database connection and migrations
- **internal/application/**: Business logic for users and tweets
- **pkg/**: Shared utilities (validation, error handling, etc.)
- **migrations/**: SQL migration files

---

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- (Optional) [Go](https://golang.org/) for local development without Docker

### Clone the Repository

```sh
git clone git@github.com:oscarsalomon89/scalable-microblogging-platform.git
cd scalable-microblogging-platform
```

### Environment Variables

Copy the example environment file and adjust as needed:

```sh
cp .env.example .env
```

### Build and Run with Docker Compose

1. Start all services (API, DB, Redis):

```sh
make compose-up
```

This will build and start all services defined in docker-compose.yml, including the API, PostgreSQL, and Redis.

API available at: http://localhost:8080
PostgreSQL: localhost:5432 (user: admin, pass: admin, db: twitterdb)
Redis: localhost:6379

2. View service logs:

```sh
make compose-logs
```

3. Stop all services:

```sh
make compose-down
```

4. Run tests:

```sh
make test
```

5. (Optional) Use alternate development environment
   If you have a
   docker-compose.dev.yml
   for a different development stack:

```sh
make compose-dev-up
# To stop dev services:
make compose-dev-down
```

6. Running Migrations Manually

To run migrations separately:

```sh
make migrate-up    # Apply all pending migrations
make migrate-down  # Revert the last migration
```

7. (Optional) Force migration

```sh
make migrate-force
```

8. Get migration version:

```sh
make migrate-version
```

---

## Project Structure

```
├── cmd/api/                # Main application entrypoint
├── internal/
│   ├── adapters/
│   │   ├── http/           # HTTP handlers
│   │   ├── postgres/       # PostgreSQL repositories
│   │   └── redis/          # Redis repositories
│   ├── application/        # Business logic
│   └── platform/           # Internal infrastructure details
├── pkg/                    # Shared utilities
├── migrations/             # SQL migration files
├── Dockerfile              # App Dockerfile
├── docker-compose.yml      # Docker Compose setup
├── .env.example            # Example environment variables
```

> **Note:**  
> For a more detailed explanation of the project structure, architectural decisions, and package responsibilities, please refer to the [project wiki](https://github.com/oscarsalomon89/scalable-microblogging-platform/wiki#-arquitectura).

---

## API Endpoints

Main endpoints:

- `POST /api/v1/users` - Register user
- `POST /api/v1/users/follow` - Follow a user
- `POST /api/v1/users/unfollow` - Unfollow a user
- `POST /api/v1/tweets` - Create tweet
- `GET /api/v1/tweets/timeline` - List tweets

> **Note:**  
> At this time, Swagger or OpenAPI documentation is not included due to project time constraints. However, you can find more detailed information about request/response formats and additional endpoints in the [project wiki](https://github.com/oscarsalomon89/scalable-microblogging-platform/wiki#-casos-de-uso).

---

## Further Documentation

Project Wiki: [Project Wiki](https://github.com/oscarsalomon89/scalable-microblogging-platform/wiki) — Architecture guides, technical decisions, usage scenarios, and more.

Design assumptions and decisions: [Design assumptions and decisions](docs/assumptions.md).

---

## Contributing

Contributions are welcome! Please fork this repository and submit a pull request.
