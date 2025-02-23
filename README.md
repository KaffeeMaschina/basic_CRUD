# basic_CRUD
This is my basic CRUD REST API written with Golang. As database is used PostgreSQL.

To start the application you need:
- to create .env file in the root of the project.
example:
```
BIND_ADDR=:8080
LOG_LEVEL=debug
PG_DATABASE_NAME=test_postgres
PG_USER=t_user
PG_PASSWORD=t_password
PG_HOST=localhost
PG_PORT=5432
MIGRATION_DIR=./migrations
```
- start a docker container with PostgreSQL, using command "docker compose up"
- "make install-deps" to get "goose" locally, it is a utility for migrations
- "make local-migration-up" to get migrations to database
- "make" to build an application
- "./apiserver" to start