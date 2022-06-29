package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

type GithubRepository struct {
	Name string `json:"name"`
	blah string `json:"blah"`
}

type GithubSender struct {
	AvatarUrl string `json:"avatar_url"`
}

type GithubPusher struct {
	Name string `json:"name"`
}

type GithubWorkflowJob struct {
	Conclusion string `json:"conclusion"`
}

type GithubWebhookEvent struct {
	Repository GithubRepository  `json:"repository"`
	Sender     GithubSender      `json:"sender"`
	Pusher     GithubPusher      `json:"pusher"`
	Action     string            `json:"action"`
	Workflow   GithubWorkflowJob `json:"workflow_job"`
}

type GithubEvents struct {
	Repository     string `json:"repository"`
	Action         string `json:"action"`
	Sender         string `json:"user"`
	Avatar         string `json:"avatar"`
	WorkflowStatus string `json:"workflow_status"`
}

func main() {
	Router = gin.Default()
	api := Router.Group("/data")
	{
		api.POST("/events", persistEvents)
		api.GET("/events", getEvents)
	}
	Router.Run(":8080")
}

func getEvents(ctx *gin.Context) {
	file, err := os.Open("data.csv")

	if err != nil {
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()

	var res []GithubEvents

	for scanner.Scan() {
		t := scanner.Text()
		entry := strings.Split(t, ",")

		res = append(res, GithubEvents{
			Repository:     entry[0],
			Action:         entry[1],
			Sender:         entry[2],
			Avatar:         entry[3],
			WorkflowStatus: entry[4],
		})

		ctx.JSON(200, res)
	}
}

func persistEvents(ctx *gin.Context) {
	var event GithubWebhookEvent

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
	entry += event.Sender.AvatarUrl + ","
	entry += event.Workflow.Conclusion

	f.WriteString(entry + "\n")

	ctx.JSON(200, event)
}

func filterWorkflowJobs(event GithubWebhookEvent) bool {
	if event.Action != "completed" {
		return false
	}

	return true
}
