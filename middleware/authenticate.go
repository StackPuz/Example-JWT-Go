package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	jwt.StandardClaims
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/" || c.Request.URL.Path == "/login" {
			c.Next()
			return
		}
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("jwt_secret")), nil
		})
		if err != nil || !token.Valid {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(*Claims); ok {
			c.Set("user", claims)
		} else {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}
		c.Next()
	}
}
