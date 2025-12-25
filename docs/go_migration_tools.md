Here are some popular database migration tools for Go projects:

### 1. golang-migrate/migrate

This is one of the most popular and feature-rich migration tools in the Go ecosystem. It's very flexible and can be used as a CLI tool or as a library in your Go application.

-   **GitHub:** [https://github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate)
-   **Key Features:**
    -   Supports a wide range of databases, including PostgreSQL, MySQL, SQLite, and more.
    -   Can be used as a CLI or imported as a library.
    -   Migrations are written in plain SQL, which makes them easy to write and review.
    -   Manages migration state in the database itself.

**CLI Usage Example:**
```bash
# Create a new migration file
migrate create -ext sql -dir db/migrations -seq create_users_table

# Apply all up migrations
migrate -path db/migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" up

# Rollback the last migration
migrate -path db/migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" down 1
```

### 2. pressly/goose

Goose is another excellent choice for managing database migrations. It's simple, straightforward, and can also be used as a CLI or a library.

-   **GitHub:** [https://github.com/pressly/goose](https://github.com/pressly/goose)
-   **Key Features:**
    -   Supports Go-based migrations in addition to SQL migrations.
    -   Simple and easy-to-use CLI.
    -   Can be embedded into your Go application.
    -   Good documentation.

**CLI Usage Example:**
```bash
# Create a new SQL migration file
goose -dir "db/migrations" create create_users_table sql

# Apply all up migrations
goose -dir "db/migrations" postgres "user=... password=... dbname=... sslmode=disable" up

# Rollback the last migration
goose -dir "db/migrations" postgres "user=... password=... dbname=... sslmode=disable" down
```

### 3. sql-migrate

This tool is a bit simpler than the others and is designed to work well with libraries that use Go's standard `database/sql` package. It's a good choice if you want a lightweight solution.

-   **GitHub:** [https://github.com/rubenv/sql-migrate](https://github.com/rubenv/sql-migrate)
-   **Key Features:**
    -   Migrations are written in SQL.
    *   Works with any database that has a `database/sql` driver.
    -   Can be used as a CLI or a library.
    -   Simple and easy to get started with.

**CLI Usage Example:**
```bash
# Create a new migration file
sql-migrate new create_users_table

# Apply all up migrations
sql-migrate up

# Rollback the last migration
sql-migrate down
```

## Recommendation

For a project using GoFrame, **`golang-migrate/migrate`** is a very solid choice due to its popularity, extensive documentation, and wide range of supported databases. It's a reliable and well-tested tool that will fit well into a professional development workflow.

`pressly/goose` is also a great option, especially if you think you might want to write some migrations in Go instead of just SQL.

You can't go wrong with any of these, but **`golang-migrate/migrate`** is generally the most recommended for new projects.
