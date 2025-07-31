# Postgres and pgx

There are two ways to connect to postgres.

- connect through pgx directly
- the `database/sql` interface with a `pgx` adapter (current)

## Notes

The `database/sql` is a connection pool
