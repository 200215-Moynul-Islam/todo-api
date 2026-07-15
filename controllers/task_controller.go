package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"todo-api/clients"
	"todo-api/repositories"
	"todo-api/services"
	"todo-api/utils"
)

var taskService services.TaskService = services.NewTaskService(
	repositories.NewTaskRepository(),
	clients.NewGeminiClient(),
)

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

type GenerateDescriptionRequest struct {
	Title string `json:"title"`
}

type GenerateDescriptionResponse struct {
	Description string `json:"description"`
}

// GenerateDescription creates a short, AI-generated description without
// persisting a task. The caller can review it before creating or updating one.
func (c *TaskController) GenerateDescription() {
	if _, ok := c.GetUserID(); !ok {
		utils.SendJSONResponse(c.Ctx, http.StatusUnauthorized, false, utils.MsgUnauthorized, nil)
		return
	}

	var req GenerateDescriptionRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, utils.MsgInvalidRequestBody, nil)
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, utils.MsgTitleRequired, nil)
		return
	}

	description, err := taskService.GenerateDescription(title)
	if err != nil {
		utils.SendJSONResponse(c.Ctx, http.StatusBadGateway, false, utils.MsgFailedToGenerateDescription, nil)
		return
	}

	utils.SendJSONResponse(c.Ctx, http.StatusOK, true, utils.MsgDescriptionGenerated, GenerateDescriptionResponse{Description: description})
}

func (c *TaskController) Create() {
	userID, ok := c.GetUserID()
	if !ok {
		utils.SendJSONResponse(c.Ctx, http.StatusUnauthorized, false, utils.MsgUnauthorized, nil)
		return
	}

	var req CreateTaskRequest

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, utils.MsgInvalidRequestBody, nil)
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, utils.MsgTitleRequired, nil)
		return
	}

	task, err := taskService.CreateTask(userID, req.Title, req.Description)
	if err != nil {
		utils.SendJSONResponse(c.Ctx, http.StatusInternalServerError, false, utils.MsgFailedToCreateTask, nil)
		return
	}

	utils.SendJSONResponse(c.Ctx, http.StatusCreated, true, utils.MsgTaskCreated, task)
}

func (c *TaskController) GetAll() {
	userID, ok := c.GetUserID()
	if !ok {
		utils.SendJSONResponse(c.Ctx, http.StatusUnauthorized, false, utils.MsgUnauthorized, nil)
		return
	}

	status := strings.TrimSpace(c.GetString("status"))

	page, err := c.GetInt("page", 1)
	if err != nil || page < 1 {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, utils.MsgInvalidRequestBody, nil)
		return
	}

	limit, err := c.GetInt("limit", 10)
	if err != nil || limit < 1 {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, utils.MsgInvalidRequestBody, nil)
		return
	}

	tasks, err := taskService.GetAllTasks(userID, status, page, limit)
	if err != nil {
		utils.SendJSONResponse(c.Ctx, http.StatusInternalServerError, false, utils.MsgFailedToRetrieveTasks, nil)
		return
	}

	utils.SendJSONResponse(c.Ctx, http.StatusOK, true, utils.MsgTasksRetrieved, tasks)
}

func (c *TaskController) GetByID() {
	userID, ok := c.GetUserID()
	if !ok {
		utils.SendJSONResponse(c.Ctx, http.StatusUnauthorized, false, utils.MsgUnauthorized, nil)
		return
	}

	id, err := c.GetInt(":id")
	if err != nil || id <= 0 {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, utils.MsgInvalidTaskID, nil)
		return
	}

	task, err := taskService.GetTaskByID(userID, id)
	if err != nil {
		utils.SendJSONResponse(c.Ctx, http.StatusInternalServerError, false, utils.MsgFailedToRetrieveTask, nil)
		return
	}

	if task == nil {
		utils.SendJSONResponse(c.Ctx, http.StatusNotFound, false, utils.MsgTaskNotFound, nil)
		return
	}

	utils.SendJSONResponse(c.Ctx, http.StatusOK, true, utils.MsgTaskRetrieved, task)
}

func (c *TaskController) Update() {
	userID, ok := c.GetUserID()
	if !ok {
		utils.SendJSONResponse(c.Ctx, http.StatusUnauthorized, false, utils.MsgUnauthorized, nil)
		return
	}

	id, err := c.GetInt(":id")
	if err != nil || id <= 0 {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, utils.MsgInvalidTaskID, nil)
		return
	}

	var req UpdateTaskRequest

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, utils.MsgInvalidRequestBody, nil)
		return
	}

	if req.Title == nil && req.Description == nil && req.Status == nil {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, utils.MsgAtLeastOneFieldRequired, nil)
		return
	}

	if req.Title != nil && strings.TrimSpace(*req.Title) == "" {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, utils.MsgTitleRequired, nil)
		return
	}

	if req.Status != nil && strings.TrimSpace(*req.Status) == "" {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, utils.MsgStatusRequired, nil)
		return
	}

	task, err := taskService.UpdateTask(userID, id, req.Title, req.Description, req.Status)
	if err != nil {
		utils.SendJSONResponse(c.Ctx, http.StatusInternalServerError, false, utils.MsgFailedToUpdateTask, nil)
		return
	}

	if task == nil {
		utils.SendJSONResponse(c.Ctx, http.StatusNotFound, false, utils.MsgTaskNotFound, nil)
		return
	}

	utils.SendJSONResponse(c.Ctx, http.StatusOK, true, utils.MsgTaskUpdated, task)
}

func (c *TaskController) Delete() {
	userID, ok := c.GetUserID()
	if !ok {
		utils.SendJSONResponse(c.Ctx, http.StatusUnauthorized, false, utils.MsgUnauthorized, nil)
		return
	}

	id, err := c.GetInt(":id")
	if err != nil || id <= 0 {
		utils.SendJSONResponse(c.Ctx, http.StatusBadRequest, false, utils.MsgInvalidTaskID, nil)
		return
	}

	ok, err = taskService.DeleteTask(userID, id)
	if err != nil {
		utils.SendJSONResponse(c.Ctx, http.StatusInternalServerError, false, utils.MsgFailedToDeleteTask, nil)
		return
	}

	if !ok {
		utils.SendJSONResponse(c.Ctx, http.StatusNotFound, false, utils.MsgTaskNotFound, nil)
		return
	}

	utils.SendJSONResponse(c.Ctx, http.StatusOK, true, utils.MsgTaskDeleted, nil)
}
