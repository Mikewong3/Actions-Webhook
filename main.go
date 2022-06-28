package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

type GithubRepository struct {
	Name string `json:"name"`
}

type GithubSender struct {
	AvatarUrl string `json:"avatar_url"`
}

type GithubPusher struct {
	Name string `json:"name"`
}

type GithubEvent struct {
	Repository GithubRepository `json:"repository"`
	Sender     GithubSender     `json:"sender"`
	Pusher     GithubPusher     `json:"pusher"`
}

func main() {
	Router = gin.Default()
	api := Router.Group("/data")
	{
		api.POST("/events", persistEvents)
	}
	Router.Run(":8080")
}

func persistEvents(ctx *gin.Context) {
	var event GithubEvent

	if err := ctx.BindJSON(&event); err != nil {
		return
	}

	fmt.Printf("ctx.GetHeader(\"X-GitHub-Event\"): %v\n", ctx.GetHeader("X-GitHub-Event"))

	fmt.Printf("event: %+v\n", event)

	f, err := os.OpenFile("data.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		return
	}

	var entry string
	entry += ctx.GetHeader("X-GitHub-Event") + ","
	entry += event.Repository.Name + ","
	entry += event.Pusher.Name + ","
	entry += event.Sender.AvatarUrl

	f.WriteString(entry + "\n")

	ctx.JSON(200, event)
}
