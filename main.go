package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/handiism/chat/pubsub"
	"github.com/handiism/chat/repo/psql"
	"github.com/handiism/chat/user"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(err.Error())
	}

	pool := psql.NewPool()
	repo := psql.NewRepo(pool)
	service := user.NewService(repo)
	user := user.NewHandler(service)
	chat := pubsub.NewHandler()

	app := gin.Default()
	app.Use(CORSMiddleware())
	app.POST("/user/login", user.Login())
	app.POST("/user/register", user.Register())
	app.GET("/user/:id", user.Fetch())
	app.GET("/chat", chat.WebSocket())

	if err := app.Run("0.0.0.0:8080"); err != nil {
		log.Fatal(err.Error())
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
