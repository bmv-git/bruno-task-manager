package entity

type Task struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty" binding:"required"`
	Description string `json:"description,omitempty"`
	Status      bool   `json:"status"`
	Priority    uint8  `json:"priority,omitempty"`
}

type TaskDTO struct {
	Tasks []Task `json:"tasks"`
	Total int    `json:"total"`
}
