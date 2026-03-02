package handlers

import (
	"net/http"

	"github.com/SaddamMohammad1/todo-rest-api-using-gin-part2/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateTodoInput struct {
	Title     string `json:"title" binding:"required"` // "title" must be provided
	Completed bool   `json:"completed"`                // Optional field
}

// Handler function for creating a new todo item
func CreateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {

	return func(c *gin.Context) {
		// Variable to store JSON input
		var input CreateTodoInput

		// Bind incoming JSON request into 'input'
		// Also validates required fields
		if err := c.ShouldBindJSON(&input); err != nil {

			// Return 400 Bad Request if JSON is invalid or missing fields
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Save todo item to database using repository function
		todo, err := repository.CreateTodo(pool, input.Title, input.Completed)

		// If DB failed, return 500 Internal Server Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Return 201 Created with created todo object
		c.JSON(http.StatusCreated, todo)
	}
}

func GetAllTodosHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		todos, err := repository.GetAllTodos(pool)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, todos)
	}
}
