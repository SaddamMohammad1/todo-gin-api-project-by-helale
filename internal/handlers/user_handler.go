package handlers

import (
	"net/http"
	"strings"

	"github.com/SaddamMohammad1/todo-rest-api-using-gin-part2/internal/models"
	"github.com/SaddamMohammad1/todo-rest-api-using-gin-part2/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func CreateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Variable to store JSON input
		var registerRequest RegisterRequest

		// Bind incoming JSON request into 'registerRequest'
		// Also validates required fields
		if err := c.ShouldBindJSON(&registerRequest); err != nil {

			// Return 400 Bad Request if JSON is invalid or missing fields
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "JSON request" + err.Error(),
			})
			return
		}

		if len(registerRequest.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Password must be atleast 6 characters long!",
			})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Falied to hash password" + err.Error(),
			})
			return
		}

		user := &models.User{
			Email:    registerRequest.Email,
			Password: string(hashedPassword),
		}

		// Save todo item to database using repository function
		createdUser, err := repository.CreateUser(pool, user)

		// If DB failed, return 500 Internal Server Error
		if err != nil {
			if strings.Contains(err.Error(), "dublicate") || strings.Contains(err.Error(), "unique") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Email already registered",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Return 201 Created with created todo object
		c.JSON(http.StatusCreated, createdUser)
	}
}
