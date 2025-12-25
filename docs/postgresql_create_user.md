To create a new PostgreSQL database user and grant them privileges to your `gva` database, follow these steps in your terminal:

### 1. Connect to PostgreSQL

First, connect to the `psql` shell as a superuser (like `postgres`):
```bash
psql -U postgres
```
You may be prompted for the password for the `postgres` user.

### 2. Create a New User

Now, use the `CREATE USER` command to create your new user. It's crucial to create the user with a password. Replace `your_new_user` and `your_strong_password` with your desired username and a secure password.

```sql
CREATE USER your_new_user WITH PASSWORD 'your_strong_password';
```

### 3. Grant Privileges to the User

Next, you need to grant the new user permissions to access and work with the `gva` database.

**Grant CONNECT privilege:** This allows the user to connect to the database.
```sql
GRANT CONNECT ON DATABASE gva TO your_new_user;
```

**Grant USAGE privilege on the schema:** This allows the user to see and use objects within the `public` schema (the default schema).
```sql
GRANT USAGE ON SCHEMA public TO your_new_user;
```

**Grant privileges on all tables in the schema:** This allows the user to perform `SELECT`, `INSERT`, `UPDATE`, and `DELETE` operations on all existing tables in the `public` schema.
```sql
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO your_new_user;
```

**Grant privileges on future tables:** This is important. It ensures that the user will automatically have privileges on new tables created in the `public` schema in the future.
```sql
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO your_new_user;
```

### 4. Verify and Exit

You can verify the user was created by listing the users:
```sql
\du
```
To exit `psql`:
```sql
\q
```

Now, `your_new_user` is ready to be used in your backend application's configuration file (`backend/manifest/config/config.yaml`).
