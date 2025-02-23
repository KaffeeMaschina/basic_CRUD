package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	defaultHost         = "localhost"
	sslmodeDisable      = "disable"
	poolMaxConn         = "10"
	poolMaxConnLifetime = "1h30m"
)

type Database struct {
	Conn *pgxpool.Pool
}

// NewDatabase creates new instance of database. It can be easily exchanged with another database.
func NewDatabase(username, password, port, database string) (*Database, error) {
	conn, err := NewPostgresConnection(username, password, port, database)
	if err != nil {
		return nil, err
	}
	db := &Database{Conn: conn}
	return db, nil
}

func NewPostgresConnection(username, password, port, database string) (*pgxpool.Pool, error) {
	const errorLocation = "internal.app.storage.postgresDB.go"
	// Create URL
	dbUrl := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s pool_max_conns=%s pool_max_conn_lifetime=%s",
		username, password, defaultHost, port, database, sslmodeDisable, poolMaxConn, poolMaxConnLifetime)

	// New Pool
	db, err := pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errorLocation, err)
	}

	// Check if connection is ok
	err = db.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("ping error %s: %w", errorLocation, err)
	}

	return db, nil
}
