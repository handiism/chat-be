package main

import (
	"log"

	"github.com/gin-contrib/cors"
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
	chat := pubsub.NewHandler(service)

	app := gin.Default()
	app.Use(cors.Default())
	app.POST("/user/login", user.Login())
	app.POST("/user/register", user.Register())
	app.GET("/user/:id", user.Fetch())
	app.GET("/chat/:id", chat.WebSocket())

	if err := app.Run("127.0.0.1:8080"); err != nil {
		log.Fatal(err.Error())
	}
}
