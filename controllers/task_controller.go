package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"todo-api/repositories"
	"todo-api/services"
)

var taskService = services.NewTaskService(repositories.NewTaskRepository())

type TaskController struct {
	BaseController
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (c *TaskController) Create() {
	var req CreateTaskRequest

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.SendError(http.StatusBadRequest, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		c.SendError(http.StatusBadRequest, "Title is required")
		return
	}

	task, err := taskService.CreateTask(req.Title, req.Description)
	if err != nil {
		c.SendError(http.StatusInternalServerError, "Failed to create task")
		return
	}

	c.SendSuccess(http.StatusCreated, "Task created successfully", task)
}

func (c *TaskController) GetAll() {
	status := strings.TrimSpace(c.GetString("status"))

	page, err := c.GetInt("page", 1)
	if err != nil || page < 1 {
		c.SendError(http.StatusBadRequest, "Invalid page parameters")
		return
	}

	limit, err := c.GetInt("limit", 10)
	if err != nil || limit < 1 {
		c.SendError(http.StatusBadRequest, "Invalid limit parameters")
		return
	}

	tasks, err := taskService.GetAllTasks(status, page, limit)
	if err != nil {
		c.SendError(http.StatusInternalServerError, "Failed to retrieve tasks")
		return
	}

	c.SendSuccess(http.StatusOK, "Tasks retrieved successfully", tasks)
}

func (c *TaskController) GetByID() {
	id, err := c.GetInt(":id")
	if err != nil || id <= 0 {
		c.SendError(http.StatusBadRequest, "Invalid task id format")
		return
	}

	task, err := taskService.GetTaskByID(id)
	if err != nil {
		c.SendError(http.StatusInternalServerError, "Failed to retrieve task")
		return
	}

	if task == nil {
		c.SendError(http.StatusNotFound, "Task not found")
		return
	}

	c.SendSuccess(http.StatusOK, "Task retrieved successfully", task)
}

func (c *TaskController) Delete() {
	id, err := c.GetInt(":id")
	if err != nil || id <= 0 {
		c.SendError(http.StatusBadRequest, "Invalid task id format")
		return
	}

	ok, err := taskService.DeleteTask(id)
	if err != nil {
		c.SendError(http.StatusInternalServerError, "Failed to delete task")
		return
	}

	if !ok {
		c.SendError(http.StatusNotFound, "Task not found")
		return
	}

	c.SendSuccess(http.StatusOK, "Task deleted successfully", nil)
}