To change the password for the PostgreSQL user `gva`, follow these steps in your terminal:

### 1. Connect to PostgreSQL

First, connect to the `psql` shell as a superuser (like `postgres`):
```bash
psql -U postgres
```
You may be prompted for the password for the `postgres` user.

### 2. Change the User's Password

Once connected to `psql`, use the `ALTER USER` command to set a new password for `gva`. Replace `your_new_strong_password` with the new password you want to set. Make sure this new password matches the one you put in `backend/manifest/config/config.yaml`.

```sql
ALTER USER gva WITH PASSWORD 'your_new_strong_password';
```

### 3. Verify and Exit

To exit `psql`:
```sql
\q
```

After changing the password in PostgreSQL, remember to update the `link` in your `backend/manifest/config/config.yaml` file to use this new password. Then, I can try running the migration again for you.
