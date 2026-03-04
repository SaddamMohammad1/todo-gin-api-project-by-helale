package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/SaddamMohammad1/todo-rest-api-using-gin-part2/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateTodo inserts a new todo into the database and returns the created record
func CreateTodo(pool *pgxpool.Pool, title string, completed bool, userID string) (*models.Todo, error) {
	// Create context with 5-second timeout
	// If query takes longer than 5 seconds → canceled automatically
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Important: release resources when done

	// SQL query using PostgreSQL's RETURNING (returns inserted row)
	var query string = `
		INSERT INTO todos (title, completed, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, title, completed, created_at, updated_at, user_id
	`

	// Variable to store returned todo data
	var todo models.Todo

	// Execute query + scan returned row into todo struct
	var err error = pool.QueryRow(ctx, query, title, completed, userID).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserID,
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

func GetAllTodos(pool *pgxpool.Pool, userID string) ([]models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, title, completed, created_at, updated_at, user_id
		FROM todos
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	var rows, err = pool.Query(ctx, query, userID)
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
			&todo.UserID,
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

func GetTodoById(pool *pgxpool.Pool, id int, userID string) (*models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, title, completed, created_at, updated_at, user_id
		FROM todos
		WHERE id = $1 AND user_id = $2
	`

	var todo models.Todo

	var err error = pool.QueryRow(ctx, query, id, userID).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserID,
	)

	if err != nil {
		return nil, err
	}

	return &todo, err
}

func UpdateTodo(pool *pgxpool.Pool, id int, title string, completed bool, userID string) (*models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		UPDATE todos
		SET title = $1, completed = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND user_id = $4
		RETURNING id, title, completed, created_at, updated_at, user_id
	`

	var todo models.Todo

	var err error = pool.QueryRow(ctx, query, title, completed, id, userID).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserID,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func DeleteTodo(pool *pgxpool.Pool, id int, userID string) error {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		DELETE FROM todos
		WHERE id = $1 AND user_id = $2
	`

	// Execute DELETE query with the provided id.
	var commandTag, err = pool.Exec(ctx, query, id, userID)

	if err != nil {
		return err
	}

	// commandTag.RowsAffected() tells how many rows were deleted.
	// If 0 → no todo exists with given ID.
	// Example: If id = 10 but no row with id 10, return not found error.
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("todo with id %d not found", id)
	}

	return nil
}
