package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"todo-api/repositories"
	"todo-api/services"
	"todo-api/utils"
)

var taskService = services.NewTaskService(repositories.NewTaskRepository())

type TaskController struct {
	BaseController
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
}

func (c *TaskController) Create() {
	userID, ok := c.GetUserID()
	if !ok {
		utils.SendJSONResponse(c.Ctx, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	var req CreateTaskRequest

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, "Invalid request body", nil)
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, "Title is required", nil)
		return
	}

	task, err := taskService.CreateTask(userID, req.Title, req.Description)
	if err != nil {
		utils.SendJSONResponse(c.Ctx, http.StatusInternalServerError, false, "Failed to create task", nil)
		return
	}

	utils.SendJSONResponse(c.Ctx, http.StatusCreated, true, "Task created successfully", task)
}

func (c *TaskController) GetAll() {
	userID, ok := c.GetUserID()
	if !ok {
		utils.SendJSONResponse(c.Ctx, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	status := strings.TrimSpace(c.GetString("status"))

	page, err := c.GetInt("page", 1)
	if err != nil || page < 1 {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, "Invalid page parameters", nil)
		return
	}

	limit, err := c.GetInt("limit", 10)
	if err != nil || limit < 1 {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, "Invalid limit parameters", nil)
		return
	}

	tasks, err := taskService.GetAllTasks(userID, status, page, limit)
	if err != nil {
		utils.SendJSONResponse(c.Ctx, http.StatusInternalServerError, false, "Failed to retrieve tasks", nil)
		return
	}

	utils.SendJSONResponse(c.Ctx, http.StatusOK, true, "Tasks retrieved successfully", tasks)
}

func (c *TaskController) GetByID() {
	userID, ok := c.GetUserID()
	if !ok {
		utils.SendJSONResponse(c.Ctx, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	id, err := c.GetInt(":id")
	if err != nil || id <= 0 {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, "Invalid task id format", nil)
		return
	}

	task, err := taskService.GetTaskByID(userID, id)
	if err != nil {
		utils.SendJSONResponse(c.Ctx, http.StatusInternalServerError, false, "Failed to retrieve task", nil)
		return
	}

	if task == nil {
		utils.SendJSONResponse(c.Ctx, http.StatusNotFound, false, "Task not found", nil)
		return
	}

	utils.SendJSONResponse(c.Ctx, http.StatusOK, true, "Task retrieved successfully", task)
}

func (c *TaskController) Update() {
	userID, ok := c.GetUserID()
	if !ok {
		utils.SendJSONResponse(c.Ctx, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	id, err := c.GetInt(":id")
	if err != nil || id <= 0 {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, "Invalid task id format", nil)
		return
	}

	var req UpdateTaskRequest

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, "Invalid request body", nil)
		return
	}

	if req.Title == nil && req.Description == nil && req.Status == nil {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, "At least one task field is required", nil)
		return
	}

	if req.Title != nil && strings.TrimSpace(*req.Title) == "" {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, "Title is required", nil)
		return
	}

	if req.Status != nil && strings.TrimSpace(*req.Status) == "" {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, "Status is required", nil)
		return
	}

	task, err := taskService.UpdateTask(userID, id, req.Title, req.Description, req.Status)
	if err != nil {
		utils.SendJSONResponse(c.Ctx, http.StatusInternalServerError, false, "Failed to update task", nil)
		return
	}

	if task == nil {
		utils.SendJSONResponse(c.Ctx, http.StatusNotFound, false, "Task not found", nil)
		return
	}

	utils.SendJSONResponse(c.Ctx, http.StatusOK, true, "Task updated successfully", task)
}

func (c *TaskController) Delete() {
	userID, ok := c.GetUserID()
	if !ok {
		utils.SendJSONResponse(c.Ctx, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	id, err := c.GetInt(":id")
	if err != nil || id <= 0 {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, "Invalid task id format", nil)
		return
	}

	ok, err = taskService.DeleteTask(userID, id)
	if err != nil {
		utils.SendJSONResponse(c.Ctx, http.StatusInternalServerError, false, "Failed to delete task", nil)
		return
	}

	if !ok {
		utils.SendJSONResponse(c.Ctx, http.StatusNotFound, false, "Task not found", nil)
		return
	}

	utils.SendJSONResponse(c.Ctx, http.StatusOK, true, "Task deleted successfully", nil)
}
