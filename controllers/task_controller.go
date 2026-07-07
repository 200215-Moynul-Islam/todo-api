package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"todo-api/models"
)

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

	now := time.Now().UTC()

	task := models.Task{
		ID:          uuid.New(),
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		Status:      "pending",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	models.CreateTask(task)

	c.SendSuccess(http.StatusCreated, "Task created successfully", task)
}

func (c *TaskController) GetAll() {
	status := strings.TrimSpace(c.GetString("status"))

	page, err := c.GetInt("page", 1)
	if err != nil || page < 1 {
		c.SendError(http.StatusBadRequest, "Invalid page")
		return
	}

	limit, err := c.GetInt("limit", 10)
	if err != nil || limit < 1 {
		c.SendError(http.StatusBadRequest, "Invalid limit")
		return
	}

	tasks := models.GetTasks(status, page, limit)

	c.SendSuccess(http.StatusOK, "Tasks retrieved successfully", tasks)
}

func (c *TaskController) GetByID() {
	id := c.Ctx.Input.Param(":id")

	taskID, err := uuid.Parse(id)
	if err != nil {
		c.SendError(http.StatusBadRequest, "Invalid task id")
		return
	}

	task, exists := models.GetTaskByID(taskID)
	if !exists {
		c.SendError(http.StatusNotFound, "Task not found")
		return
	}

	c.SendSuccess(http.StatusOK, "Task retrieved successfully", task)
}

func (c *TaskController) Delete() {
	id := c.Ctx.Input.Param(":id")

	taskID, err := uuid.Parse(id)
	if err != nil {
		c.SendError(http.StatusBadRequest, "Invalid task id")
		return
	}

	if !models.DeleteTask(taskID) {
		c.SendError(http.StatusNotFound, "Task not found")
		return
	}

	c.SendSuccess(http.StatusOK, "Task deleted successfully", nil)
}