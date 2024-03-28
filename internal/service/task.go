package service

import "bruno-task-manager/internal/entity"

func GetTaskByID(id string, t entity.TaskDTO) (entity.Task, int, bool) {
	for idx, task := range t.Tasks {
		if task.ID == id {
			return task, idx, true
		}
	}
	return entity.Task{}, -1, false
}

func DeleteTaskByID(id string, t entity.TaskDTO) bool {
	_, idx, ok := GetTaskByID(id, t)
	if ok {
		t.Tasks = append(t.Tasks[:idx], t.Tasks[idx+1:]...)
		t.Total--
		return true
	}
	return false
}
