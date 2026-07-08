package services

import (
	"strings"
	"todo-api/models"
	"todo-api/repositories"
)

type TaskService interface {
	CreateTask(title, description string) (*models.Task, error)
	GetAllTasks(status string, page, limit int) ([]models.Task, error)
	GetTaskByID(id int) (*models.Task, error)
	DeleteTask(id int) (bool, error)
}

type taskService struct {
	repo repositories.TaskRepository
}

func NewTaskService(repo repositories.TaskRepository) TaskService {
	return &taskService{
		repo: repo,
	}
}

func (s *taskService) CreateTask(title, description string) (*models.Task, error) {
	task := &models.Task{
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		Status:      "pending",
	}

	err := s.repo.Create(task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *taskService) GetAllTasks(status string, page, limit int) ([]models.Task, error) {
	return s.repo.GetAll(strings.TrimSpace(status), page, limit)
}

func (s *taskService) GetTaskByID(id int) (*models.Task, error) {
	return s.repo.GetByID(id)
}

func (s *taskService) DeleteTask(id int) (bool, error) {
	return s.repo.Delete(id)
}