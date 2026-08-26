# libraryAPI

A small REST API for managing a library: users, books, and book loans.

> This project was built as a take-home assignment for a **Backend Developer (Go)**
> job interview. It aims to show clean, idiomatic Go, correct use of SQL and
> transactions, and a sensible local development setup — not to be a
> feature-complete product.

## Tech stack

- **Go 1.25**
- **[Gin](https://github.com/gin-gonic/gin)** — HTTP router
- **`database/sql`** + **[lib/pq](https://github.com/lib/pq)** — PostgreSQL driver (no ORM)
- **PostgreSQL 16**
- **[golang-migrate](https://github.com/golang-migrate/migrate)** — schema migrations (embedded, run on startup)
- **Docker Compose** — local database

## Project layout

```
.
├── main.go            # entrypoint, routes
├── db/                # DB connection + migration runner
├── models/            # domain types (User, Book, BookLoan)
├── handlers/          # HTTP handlers (users, books, loans)
└── devops/            # docker-compose, SQL migrations, DB README
```

## Getting started

### 1. Start PostgreSQL

```bash
cd devops
docker compose up -d
```

See [`devops/README.md`](devops/README.md) for details.

### 2. Configure the connection (optional)

The app reads `DATABASE_URL`, falling back to the local docker database:

```
postgresql://library:library@localhost:5432/librarydb?sslmode=disable
```

Copy `.env.example` to `.env` and adjust if needed.

### 3. Run the API

```bash
go run .
```

On startup the app connects to the database and applies any pending migrations
automatically. The server listens on `:8000`.

## API

| Method | Path                  | Description                          |
| ------ | --------------------- | ------------------------------------ |
| GET    | `/users`              | List all users                       |
| POST   | `/users`              | Create a user                        |
| GET    | `/users/:id`          | Get a single user                    |
| GET    | `/users/:id/loans`    | List a user's loans                  |
| GET    | `/books`              | List all books                       |
| POST   | `/books`              | Create a book                        |
| GET    | `/books/available`    | List books with available copies     |
| GET    | `/loans`              | List all loans (with user + book)    |
| POST   | `/loans/new`          | Borrow a book                        |
| POST   | `/loans/return/:id`   | Return a borrowed book               |

### Examples

Create a user:

```bash
curl -X POST localhost:8000/users \
  -H 'Content-Type: application/json' \
  -d '{"first_name":"Ada","last_name":"Lovelace"}'
```

Create a book (`total` = number of copies):

```bash
curl -X POST localhost:8000/books \
  -H 'Content-Type: application/json' \
  -d '{"title":"The Go Programming Language","total":3}'
```

Borrow a book:

```bash
curl -X POST localhost:8000/loans/new \
  -H 'Content-Type: application/json' \
  -d '{"user_id":"<user-uuid>","book_id":"<book-uuid>"}'
```

Return a book:

```bash
curl -X POST localhost:8000/loans/return/<loan-uuid>
```

## Design notes

- **Transactions** guard the loan flow: borrowing decrements a book's available
  count and inserts the loan atomically; returning closes the loan and increments
  the count atomically. A crash mid-operation can't desync inventory.
- **Availability** is protected at two levels: the `available > 0` guard when
  borrowing, and a `CHECK (available >= 0 AND available <= total)` constraint in
  the schema.
- **UUIDs** from the client are parsed defensively (`400` on malformed input,
  never a panic).
- **Migrations** are embedded in the binary and applied on startup, so the app is
  self-contained — no separate migration step or CLI required to run it.

## Possible improvements

Given more time: request-scoped logging, a config struct instead of ad-hoc env
reads, pagination on list endpoints, integration tests against a throwaway
Postgres container, and graceful shutdown.
