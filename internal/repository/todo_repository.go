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

func GetAllTodos(pool *pgxpool.Pool) ([]models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, title, completed, created_at, updated_at
		FROM todos
		ORDER BY created_at DESC
	`

	var rows, err = pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Slice to store all todos
	var todos []models.Todo = []models.Todo{}

	/*
		Next() :
			rows contains multiple records from the database.
			But Go does not give you all rows at once.
			Instead, you must iterate over them one by one — like reading lines from a file.
	*/
	// Loop through result rows
	for rows.Next() {
		var todo models.Todo

		// Scan row data into todo struct fields
		err = rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Add todo to slice
		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return all todos
	return todos, nil
}

func GetTodoById(pool *pgxpool.Pool, id int) (*models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, title, completed, created_at, updated_at
		FROM todos
		WHERE id = $1
	`

	var todo models.Todo

	var err error = pool.QueryRow(ctx, query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &todo, err
}
