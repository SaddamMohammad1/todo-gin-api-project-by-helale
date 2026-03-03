package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/SaddamMohammad1/todo-rest-api-using-gin-part2/internal/config"
	"github.com/SaddamMohammad1/todo-rest-api-using-gin-part2/internal/models"
	"github.com/SaddamMohammad1/todo-rest-api-using-gin-part2/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
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

		// Save user item to database using repository function
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

		// Return 201 Created with created user object
		c.JSON(http.StatusCreated, createdUser)
	}
}

func LoginHandler(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loginRequest LoginRequest

		if err := ctx.BindJSON(&loginRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		user, err := repository.GetUserByEmail(pool, loginRequest.Email)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
			})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
			})
			return
		}

		// map[string]any{}
		claims := jwt.MapClaims{
			"user_id": user.ID,
			"email":   user.Email,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString([]byte(cfg.JWTSecret))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate token: " + err.Error(),
			})
		}

		ctx.JSON(http.StatusOK, LoginResponse{
			Token: tokenString,
		})
	}
}

// Hanlder for testing middleware
func TestProtectedHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("user_id")

		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Protected route accessed successfully!",
			"user_id": userID,
		})
	}
}
