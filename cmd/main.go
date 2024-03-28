package main

import (
	"bruno-task-manager/internal/endpoints"
	"bruno-task-manager/internal/handlers"
	"bruno-task-manager/internal/tools"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	tasks := handlers.NewTaskRoutes()
	err := tools.LoadTasksFromFile(&tasks.TaskResponse)
	if err != nil {
		log.Fatal(err)
		return
	}
	r := gin.Default()
	endpoints.InitEndpoints(r, tasks)
	r.Run(":8080")
}
