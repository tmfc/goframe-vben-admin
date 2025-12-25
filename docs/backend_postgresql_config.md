I have configured the backend to use PostgreSQL. Here's a summary of the changes:

1.  **Created `backend/manifest/config/config.yaml`:**
    This file now contains the database configuration for your project.

    ```yaml
    # Database configuration for PostgreSQL
    database:
      default:
        link: "pgsql:postgres:123456@tcp(127.0.0.1:5432)/goframe_vben_admin?sslmode=disable"
        type: "pgsql"
        debug: true # Set to false in production
    ```

    **IMPORTANT:** You need to edit this file and replace `postgres:123456` with your actual PostgreSQL username and password.

2.  **Added PostgreSQL Driver:**
    I have added the necessary PostgreSQL driver to `backend/main.go` and updated the `go.mod` and `go.sum` files by running `go mod tidy`.

Your backend is now configured to use the PostgreSQL database named `goframe_vben_admin` on your local machine.

You can now start your backend server (`cd backend && gf run main.go`), and it will attempt to connect to the PostgreSQL database.
