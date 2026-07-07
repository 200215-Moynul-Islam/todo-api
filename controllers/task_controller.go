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