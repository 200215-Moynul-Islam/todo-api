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
	UpdateTask(id int, title, description, status *string) (*models.Task, error)
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

func (s *taskService) UpdateTask(id int, title, description, status *string) (*models.Task, error) {
	task, err := s.repo.GetByID(id)
	if err != nil || task == nil {
		return task, err
	}

	if title != nil {
		task.Title = strings.TrimSpace(*title)
	}
	if description != nil {
		task.Description = strings.TrimSpace(*description)
	}
	if status != nil {
		task.Status = strings.TrimSpace(*status)
	}

	if err := s.repo.Update(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *taskService) DeleteTask(id int) (bool, error) {
	return s.repo.Delete(id)
}
