package endpoints

import (
	"bruno-task-manager/internal/handlers"
	"github.com/gin-gonic/gin"
)

func InitEndpoints(r *gin.Engine, tr *handlers.TaskRoutes) {
	r.GET("/all", tr.GetAllTasks)
	r.POST("/task", tr.CreateTask)
	r.PUT("/task/:id", tr.UpdateTask)
	r.DELETE("/tasks/:id", tr.DeleteTask)
	r.GET("/tasks", tr.ListTasks)
}
