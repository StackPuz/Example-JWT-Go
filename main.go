package main

import (
	"app/middleware"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	router := gin.Default()
	router.Use(middleware.Authenticate())
	router.LoadHTMLFiles("public/index.html", "public/login.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	router.GET("/user", func(c *gin.Context) {
		user, _ := c.Get("user")
		claims := user.(*middleware.Claims)
		c.JSON(http.StatusOK, gin.H{"name": claims.Name})
	})

	router.POST("/login", func(c *gin.Context) {
		var login map[string]string
		c.BindJSON(&login)
		if login["name"] == "admin" && login["password"] == "1234" {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, &middleware.Claims{
				Id:   1,
				Name: login["name"],
				StandardClaims: jwt.StandardClaims{
					IssuedAt:  time.Now().Unix(),
					ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
				},
			})
			tokenString, _ := token.SignedString([]byte(os.Getenv("jwt_secret")))
			c.JSON(http.StatusOK, gin.H{"token": tokenString})
		} else {
			c.Status(http.StatusBadRequest)
		}
	})
	router.Run()
}
