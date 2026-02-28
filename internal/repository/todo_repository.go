package repository

import (
	"context"
	"time"

	"github.com/SaddamMohammad1/todo-rest-api-using-gin-part2/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateTodo inserts a new todo into the database and returns the created record
func CreateTodo(pool *pgxpool.Pool, title string, completed bool) (*models.Todo, error) {
	// Create context with 5-second timeout
	// If query takes longer than 5 seconds → canceled automatically
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Important: release resources when done

	// SQL query using PostgreSQL's RETURNING (returns inserted row)
	var query string = `
		INSERT INTO todos (title, completed)
		VALUES ($1, $2)
		RETURNING id, title, completed, created_at, updated_at
	`

	// Variable to store returned todo data
	var todo models.Todo

	// Execute query + scan returned row into todo struct
	var err error = pool.QueryRow(ctx, query, title, completed).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	/*
		Note:
			Scan() is used to copy values from a database row INTO Go variables.
			It reads the columns returned by your SQL query and fills your struct fields.
	*/

	// If anything goes wrong (DB error, timeout, scan error)
	if err != nil {
		return nil, err
	}

	// Return the created todo
	return &todo, nil
}
