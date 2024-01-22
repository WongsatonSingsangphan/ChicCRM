package middlewares

import (
	"chicCRM/pkg/auth"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var jwtSecret = []byte("thenilalive")
		providedToken := c.Request.Header.Get("Authorization")

		if providedToken == "" {
			c.JSON(401, gin.H{"status": "Error", "message": "Missing token"})
			c.Abort()
			return
		}

		// Extract the token from the "Bearer <token>" format
		providedToken = strings.TrimPrefix(providedToken, "Bearer ")

		if auth.IsTokenBlacklisted(db, providedToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Error", "message": "THE TOKEN ALREADY USED : ONE TIME TOKEN"})
			c.Abort()
			return
		}

		// Verify the token
		token, err := jwt.Parse(providedToken, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"status": "Error", "message": "Invalid Token"})
			c.Abort()
			return
		}

		// Set the claims in the context
		c.Set("claims", token.Claims)
		c.Next()
	}
}

func AuthMiddlewareResetPassword(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var jwtSecret = []byte("thenilalive")
		providedToken := c.Request.Header.Get("Authorization")

		if providedToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Error", "message": "Missing token"})
			c.Abort()
			return
		}

		// Extract the token from the "Bearer <token>" format
		providedToken = strings.TrimPrefix(providedToken, "Bearer ")

		// Check if the token matches what's in the database
		if !auth.IsTokenMatched(db, providedToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Error", "message": "Token does not match our records"})
			c.Abort()
			return
		}

		// Verify the token
		// token, err := jwt.Parse(providedToken, func(token *jwt.Token) (interface{}, error) {
		token, _ := jwt.Parse(providedToken, func(token *jwt.Token) (interface{}, error) {

			// Ensure the token algorithm is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		// if err != nil || !token.Valid {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"status": "Error", "message": "Invalid token"})
		// 	c.Abort()
		// 	return
		// }

		// Token is valid and matches the database, proceed with the request
		c.Set("claims", token.Claims)
		c.Next()
	}
}
