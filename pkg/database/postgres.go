// backend/pkg/database/postgres.go
package database

import (
	"database/sql"
	"time"

	// We use the pgx driver which is highly recommended for PostgreSQL.
	// The blank import is the standard way to register a driver.
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Connect establishes a connection pool to the PostgreSQL database.
func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify the connection is alive.
	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Configure the connection pool for better performance.
	db.SetMaxOpenConns(25)                 // Max number of open connections to the database.
	db.SetMaxIdleConns(25)                 // Max number of connections in the idle connection pool.
	db.SetConnMaxLifetime(5 * time.Minute) // Max amount of time a connection may be reused.

	return db, nil
}
