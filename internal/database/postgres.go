package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect function is used to connect to the PostgreSQL database
// using pgxpool (connection pooling).
// It returns a *pgxpool.Pool when connection is successful.
func Connect(databaseURL string) (*pgxpool.Pool, error) {
	// Context is required for DB operations such as opening the pool, pinging the database, etc.
	var ctx context.Context = context.Background()

	// config will store parsed DB details like host, port, user, password.
	// Example databaseURL: "postgres://postgres:1234@localhost:5432/mydb"
	var config *pgxpool.Config
	var err error
	config, err = pgxpool.ParseConfig(databaseURL)

	if err != nil {
		log.Printf("Unable to parse Database_url: %v", err)
		return nil, err
	}

	// Connection pool maintains multiple reusable connections.
	var pool *pgxpool.Pool
	pool, err = pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		log.Printf("Unable to create connection pool: %v", err)
		return nil, err
	}

	// Ping checks whether DB is really reachable.
	err = pool.Ping(ctx)

	if err != nil {
		log.Printf("Unable to ping database: %v", err)
		pool.Close()
		return nil, err
	}

	log.Println("Successfully conntected to PostgreSQL database")

	// Return a fully working database pool
	return pool, nil
}
