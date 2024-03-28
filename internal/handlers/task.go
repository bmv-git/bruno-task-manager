package handlers

import (
	"bruno-task-manager/internal/entity"
	"bruno-task-manager/internal/service"
	"bruno-task-manager/internal/tools"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type TaskRoutes struct {
	TaskResponse entity.TaskDTO
}

func NewTaskRoutes() *TaskRoutes {
	return &TaskRoutes{
		TaskResponse: entity.TaskDTO{Tasks: make([]entity.Task, 0, 10), Total: 0},
	}
}

func (t *TaskRoutes) CreateTask(c *gin.Context) {
	var task entity.Task
	err := c.BindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	task.ID = uuid.New().String()
	t.TaskResponse.Tasks = append(t.TaskResponse.Tasks, task)
	c.JSON(http.StatusOK, gin.H{"message": "task created"})
	err = tools.SaveTasksToFile(t.TaskResponse.Tasks)
	if err != nil {
		c.JSON(http.StatusMultiStatus, gin.H{"error": err.Error()})
	}
} // чтение JSON из body

func (t *TaskRoutes) GetAllTasks(c *gin.Context) {
	statusStr, existsStatus := c.GetQuery("status")
	priorityStr, existsPriority := c.GetQuery("priority")
	if existsStatus && existsPriority {
		c.Header("Cache-Control", "public, max-age=3600")
		c.JSON(http.StatusOK, t.TaskResponse.Tasks)
		return
	}

	status, err := strconv.ParseBool(statusStr) // преобразование строки вида bool в тип bool
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	priority, err := strconv.Atoi(priorityStr) // преобразование строки вида int в тип int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filterTasks := make([]entity.Task, 0, len(t.TaskResponse.Tasks)) // временный срез для выборки задач из среза t.Tasks
	for _, task := range t.TaskResponse.Tasks {
		if task.Status == status && task.Priority == uint8(priority) {
			filterTasks = append(filterTasks, task)
		}
	}
	c.JSON(http.StatusOK, filterTasks)
} // чтение из запросов query

func (t *TaskRoutes) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	task, idx, ok := service.GetTaskByID(id, t.TaskResponse)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	err := c.BindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t.TaskResponse.Tasks[idx] = task
	c.JSON(http.StatusOK, task)
	err = tools.LoadTasksFromFile(&t.TaskResponse)
	if err != nil {
		c.JSON(http.StatusMultiStatus, gin.H{"error": err.Error()})
	}
} // чтение параметра

func (t *TaskRoutes) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	ok := service.DeleteTaskByID(id, t.TaskResponse)
	if ok {
		c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
		return
	}
	err := tools.SaveTasksToFile(t.TaskResponse.Tasks)
	if err != nil {
		c.JSON(http.StatusMultiStatus, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task not found"})
} // чтение параметра

func (t *TaskRoutes) ListTasks(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	iL := (page - 1) * 10
	if iL >= len(t.TaskResponse.Tasks) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "на этой странице нет задач"})
		return
	}
	iH := page * 10
	if iH > len(t.TaskResponse.Tasks) {
		iH = len(t.TaskResponse.Tasks)
	}
	slice := t.TaskResponse.Tasks[iL:iH]
	c.JSON(http.StatusOK, slice)
	response := make(map[string]interface{})
	response["tasks"] = slice
	response["total"] = len(t.TaskResponse.Tasks)
	c.JSON(http.StatusOK, response)
} // чтение из запроса query
