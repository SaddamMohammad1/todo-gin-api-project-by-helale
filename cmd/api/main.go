package main

import (
	"log"

	"github.com/SaddamMohammad1/todo-rest-api-using-gin-part2/internal/config"
	"github.com/SaddamMohammad1/todo-rest-api-using-gin-part2/internal/database"
	"github.com/SaddamMohammad1/todo-rest-api-using-gin-part2/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	var cfg *config.Config
	var err error
	cfg, err = config.Load()

	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	var pool *pgxpool.Pool
	pool, err = database.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	defer pool.Close()

	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil) // Hide proxy message two line in terminal
	router.GET("/", func(c *gin.Context) {
		// map[string]any{}
		c.JSON(200, gin.H{
			"message":  "Todo API is running!",
			"status":   "success",
			"database": "connected",
		})
	})

	router.POST("/todos", handlers.CreateTodoHandler(pool))
	router.GET("/todos", handlers.GetAllTodosHandler(pool))
	router.GET("/todos/:id", handlers.GetTodoByIdHandler(pool))

	router.Run(":" + cfg.Port)
}
