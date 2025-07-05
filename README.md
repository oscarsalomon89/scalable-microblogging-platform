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

This will start the API, PostgreSQL, Redis, and run DB migrations automatically:

```sh
docker-compose up --build
```

- API available at: `http://localhost:8080`
- PostgreSQL: `localhost:5432` (user: admin, pass: admin, db: twitterdb)
- Redis: `localhost:6379`

### Running Migrations Manually

To run migrations separately:

```sh
docker-compose run --rm migrate
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

---

## API Endpoints

Main endpoints:

- `POST /v1/users` - Register user
- `POST /v1/users/follow` - Follow a user
- `POST /v1/tweets` - Create tweet
- `GET /v1/tweets` - List tweets

(See source code for full details)

---

## Contributing

Contributions are welcome! Please fork this repository and submit a pull request.

---

## License

This project is licensed under the MIT License.
