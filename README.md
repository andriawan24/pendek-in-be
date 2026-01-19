# pendek-in

`pendek-in` is a high-performance, feature-rich URL shortener service built with Go. It provides a robust RESTful API for creating short links, tracking detailed click analytics, and managing user accounts with JWT and Google OAuth integration.

## 🚀 Features

-   **Link Shortening:** Create custom or randomly generated short codes for long URLs.
-   **Advanced Analytics:** Track clicks, browser information, and geolocation (Country-level).
-   **User Authentication:** Secure access using JWT (JSON Web Tokens) and Google OAuth 2.0.
-   **Profile Management:** User profiles with support for profile image uploads.
-   **Performance:** Optimized with Redis caching for fast redirections.
-   **API Documentation:** Interactive Swagger UI for easy API exploration.
-   **Database Safety:** Type-safe SQL queries generated via `sqlc` and versioned migrations with `goose`.

## 🛠 Tech Stack

-   **Backend:** [Go](https://go.dev/) (1.25+)
-   **Web Framework:** [Gin Gonic](https://gin-gonic.com/)
-   **Database:** [PostgreSQL](https://www.postgresql.org/)
-   **SQL Generator:** [sqlc](https://sqlc.dev/)
-   **Migrations:** [Goose](https://github.com/pressly/goose)
-   **Caching:** [Redis](https://redis.io/)
-   **API Docs:** [Swagger / Swag](https://github.com/swaggo/swag)
-   **Auth:** JWT & Google OAuth 2.0

## 🏁 Getting Started

### Prerequisites

-   Go 1.25 or later
-   PostgreSQL
-   Redis
-   `make` (optional, but recommended)

### Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/andriawan24/link-short.git
    cd link-short
    ```

2.  **Set up environment variables:**
    ```bash
    cp .env.example .env
    ```
    Edit `.env` and fill in your database, Redis, and OAuth credentials.

3.  **Install dependencies:**
    ```bash
    make deps
    ```

4.  **Run database migrations:**
    ```bash
    make migrate-up
    ```

### Running the Application

To start the server in development mode with hot reload (requires `air`):
```bash
make run-watch
```

Standard run:
```bash
make run
```

The server will be available at `http://localhost:8080` (or your configured `HTTP_PORT`).

## 📊 Database Management

This project uses `sqlc` for type-safe database access and `goose` for migrations.

-   **Generate Go code from SQL:** `make sqlc`
-   **Create a new migration:** `make migrate-create name=migration_name`
-   **Apply migrations:** `make migrate-up`
-   **Rollback migration:** `make migrate-down`

## 📖 API Documentation

Once the server is running, you can access the interactive Swagger documentation at:
`http://localhost:8080/swagger/index.html`

To update the documentation after changing code comments:
```bash
make swagger
```

## 🧪 Testing

Run the test suite:
```bash
make test
```

Run tests with coverage:
```bash
make coverage
```

## 📜 License

Distributed under the MIT License. See `LICENSE` (if applicable) for more information.
