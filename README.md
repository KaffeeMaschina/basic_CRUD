# basic_CRUD
This is my basic CRUD REST API written with Golang. As database is used PostgreSQL.

To start the application you need:
- to create .env file in the root of the project.
```
BIND_ADDR=
LOG_LEVEL=
PG_DATABASE_NAME=
PG_USER=
PG_PASSWORD=
PG_HOST=
PG_PORT=
MIGRATION_DIR=
```
- start a docker container with PostgreSQL, using command "docker compose up"
- "make install-deps" to get "goose" locally, it is a utility for migrations
- "make local-migration-up" to get migrations to database
- "make" to build an application
- "./apiserver" to start