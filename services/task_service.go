package services

import (
	"strings"
	"todo-api/models"
	"todo-api/repositories"
)

type TaskService interface {
	CreateTask(userID int, title, description string) (*models.Task, error)
	GetAllTasks(userID int, status string, page, limit int) ([]models.Task, error)
	GetTaskByID(userID int, id int) (*models.Task, error)
	UpdateTask(userID int, id int, title, description, status *string) (*models.Task, error)
	DeleteTask(userID int, id int) (bool, error)
}

type taskService struct {
	repo repositories.TaskRepository
}

func NewTaskService(repo repositories.TaskRepository) TaskService {
	return &taskService{
		repo: repo,
	}
}

func (s *taskService) CreateTask(userID int, title, description string) (*models.Task, error) {
	task := &models.Task{
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		Status:      "pending",
		User:        &models.User{ID: userID},
	}

	err := s.repo.Create(task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *taskService) GetAllTasks(userID int, status string, page, limit int) ([]models.Task, error) {
	return s.repo.GetAll(userID, strings.TrimSpace(status), page, limit)
}

func (s *taskService) GetTaskByID(userID int, id int) (*models.Task, error) {
	task, err := s.repo.GetByID(id)
	if err != nil || task == nil {
		return nil, err
	}
	// Check ownership
	if task.User == nil || task.User.ID != userID {
		return nil, nil // Return nil if it doesn't belong to the user
	}
	return task, nil
}

func (s *taskService) UpdateTask(userID int, id int, title, description, status *string) (*models.Task, error) {
	task, err := s.repo.GetByID(id)
	if err != nil || task == nil {
		return nil, err
	}
	// Check ownership
	if task.User == nil || task.User.ID != userID {
		return nil, nil
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

func (s *taskService) DeleteTask(userID int, id int) (bool, error) {
	task, err := s.repo.GetByID(id)
	if err != nil || task == nil {
		return false, err
	}
	// Check ownership
	if task.User == nil || task.User.ID != userID {
		return false, nil
	}

	return s.repo.Delete(id)
}
