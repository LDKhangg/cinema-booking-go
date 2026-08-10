# Cinema Booking Go (Backend API)

A robust backend REST API for a Cinema Booking system built with **Go** (Golang), **PostgreSQL**, and clean architecture patterns.

## 🚀 Features

- **Movie Management**: Manage movies, release/close dates, publishing status, and retrieval of published movies.
- **Theater & Room Management**: Manage theaters, screening rooms, and automated seat matrix generation (standard & VIP seats).
- **Showtime Management**: Schedule screenings for movies in specific rooms with automated conflict/overlap detection and pricing.
- **Clean Architecture**: Organized into distinct layers (`entity`, `repository`, `service`, `handler`).
- **Standard Library Routing**: Uses Go 1.22+ native `net/http` ServeMux with method and path routing (`GET /movies`, `POST /theaters`, etc.).

---

## 🛠️ Tech Stack

- **Language:** Go (1.22+)
- **Database:** PostgreSQL (with `database/sql` & `lib/pq` / `pgx`)
- **Containerization:** Docker Compose

---

## 📂 Project Structure

```text
├── cmd/
│   └── main.go           # Application entrypoint
├── internal/
│   ├── movie/            # Movie domain (entity, repo, service, handler)
│   ├── theater/          # Theater & Room domain
│   ├── showtime/         # Showtime domain
│   └── server/           # HTTP router setup
├── pkg/
│   └── database/         # Database connection setup
├── compose.yml           # PostgreSQL docker-compose service
└── go.mod
```

---

## 🏃 Getting Started

### 1. Prerequisites
- Go 1.22+ installed
- Docker & Docker Compose (for PostgreSQL)

### 2. Start PostgreSQL via Docker Compose
```bash
docker compose up -d
```

### 3. Run the Application
```bash
go run cmd/main.go
```

The server will start on port `8080`.
