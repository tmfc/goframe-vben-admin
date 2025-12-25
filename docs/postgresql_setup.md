To create a new PostgreSQL database, follow these steps:

1.  **Open your terminal or command prompt.**

2.  **Connect to the PostgreSQL server:**
    You can connect as the default `postgres` superuser (or your own user if you have one set up with create database privileges).
    ```bash
    psql -U postgres
    ```
    (You might be prompted for the `postgres` user's password.)

3.  **Create the new database:**
    Once connected to `psql`, use the `CREATE DATABASE` command. You can name your database as per your project's needs. For this project, a good name would be `goframe_vben_admin`.
    ```sql
    CREATE DATABASE goframe_vben_admin;
    ```

4.  **Verify the database creation (Optional):**
    You can list all databases to confirm yours was created:
    ```sql
    \l
    ```
    Then, to exit `psql`:
    ```sql
    \q
    ```

Now your PostgreSQL database `goframe_vben_admin` is created and ready for use.