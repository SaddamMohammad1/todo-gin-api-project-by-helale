package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SaddamMohammad1/todo-rest-api-using-gin-part2/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == "" || tokenString == authHeader {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			ctx.Abort()
			return
		}

		userID, ok := claims["user_id"].(string)

		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			ctx.Abort()
			return
		}

		if exp, ok := claims["exp"].(float64); ok {
			expirationTime := time.Unix(int64(exp), 0)

			if time.Now().After(expirationTime) {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"error": "Token has expired",
				})
				ctx.Abort()
				return
			}
		}

		ctx.Set("user_id", userID) // Set user_id in context for handlers to use
		ctx.Next()
	}
}
