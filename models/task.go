package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var tasks = make(map[uuid.UUID]Task)

func CreateTask(task Task) {
	tasks[task.ID] = task
}

func GetTasks(status string, page, limit int) []Task {
	var filteredTasks []Task

	for _, task := range tasks {
		if status == "" || task.Status == status {
			filteredTasks = append(filteredTasks, task)
		}
	}

	start := (page - 1) * limit
	if start >= len(filteredTasks) {
		return []Task{}
	}

	// Calculate end index.
	end := min(start+limit, len(filteredTasks))

	return filteredTasks[start:end]
}

func GetTaskByID(id uuid.UUID) (Task, bool) {
	task, exists := tasks[id]
	return task, exists
}

func DeleteTask(id uuid.UUID) bool {
	if _, exists := tasks[id]; !exists {
		return false
	}

	delete(tasks, id)
	return true
}