package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"todo-api/mocks"
	"todo-api/models"
	"todo-api/services"
	"todo-api/utils"

	"github.com/beego/beego/v2/server/web/context"
	"go.uber.org/mock/gomock"
)

var errDatabase = errors.New("database error")
var errGemini = errors.New("gemini error")

func TestTaskController_GenerateDescription(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		authenticated bool
		setupMock     func(*mocks.MockTaskService)
		wantStatus    int
		wantMsg       string
		wantDesc      string
	}{
		{
			name:          "unauthorized",
			authenticated: false,
			wantStatus:    http.StatusUnauthorized,
			wantMsg:       utils.MsgUnauthorized,
		},
		{
			name:          "invalid request body",
			authenticated: true,
			body:          `{"title":`,
			wantStatus:    http.StatusBadRequest,
			wantMsg:       utils.MsgInvalidRequestBody,
		},
		{
			name:          "empty title",
			authenticated: true,
			body:          `{"title":"   "}`,
			wantStatus:    http.StatusBadRequest,
			wantMsg:       utils.MsgTitleRequired,
		},
		{
			name:          "service error",
			authenticated: true,
			body:          `{"title":"Learn Go"}`,
			setupMock: func(service *mocks.MockTaskService) {
				service.EXPECT().
					GenerateDescription("Learn Go").
					Return("", errGemini)
			},
			wantStatus: http.StatusBadGateway,
			wantMsg:    utils.MsgFailedToGenerateDescription,
		},
		{
			name:          "success",
			authenticated: true,
			body:          `{"title":"Learn Go"}`,
			setupMock: func(service *mocks.MockTaskService) {
				service.EXPECT().
					GenerateDescription("Learn Go").
					Return("Practice Go testing every day.", nil)
			},
			wantStatus: http.StatusOK,
			wantMsg:    utils.MsgDescriptionGenerated,
			wantDesc:   "Practice Go testing every day.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockTaskService(ctrl)

			restore := replaceTaskService(mockService)
			defer restore()

			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			ctx, rec := newTestContext(
				http.MethodPost,
				"/tasks",
				tt.body,
				tt.authenticated,
			)

			controller := &TaskController{}
			controller.Ctx = ctx

			controller.GenerateDescription()

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}

			var response struct {
				Success bool                          `json:"success"`
				Message string                        `json:"message"`
				Data    *GenerateDescriptionResponse  `json:"data"`
			}

			if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if response.Success != (tt.wantStatus < 400) {
				t.Fatalf("expected success %v, got %v", tt.wantStatus < 400, response.Success)
			}

			if response.Message != tt.wantMsg {
				t.Fatalf("expected message %q, got %q", tt.wantMsg, response.Message)
			}

			if tt.wantDesc != "" {
				if response.Data == nil {
					t.Fatal("expected response data")
				}

				if response.Data.Description != tt.wantDesc {
					t.Fatalf("expected description %q, got %q", tt.wantDesc, response.Data.Description)
				}
			} else if response.Data != nil {
				t.Fatal("expected response data to be nil")
			}
		})
	}
}

func TestTaskController_Create(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		authenticated bool
		setupMock     func(*mocks.MockTaskService)
		wantStatus    int
		wantMsg       string
		wantTask      *models.Task
	}{
		{
			name:          "unauthorized",
			authenticated: false,
			wantStatus:    http.StatusUnauthorized,
			wantMsg:       utils.MsgUnauthorized,
		},
		{
			name:          "invalid request body",
			authenticated: true,
			body:          `{"title":`,
			wantStatus:    http.StatusBadRequest,
			wantMsg:       utils.MsgInvalidRequestBody,
		},
		{
			name:          "empty title",
			authenticated: true,
			body:          `{"title":"   ","description":"Learn testing"}`,
			wantStatus:    http.StatusBadRequest,
			wantMsg:       utils.MsgTitleRequired,
		},
		{
			name:          "service error",
			authenticated: true,
			body:          `{"title":"Learn Go","description":"Practice testing"}`,
			setupMock: func(service *mocks.MockTaskService) {
				service.EXPECT().
					CreateTask(
						1,
						"Learn Go",
						"Practice testing",
					).
					Return(nil, errDatabase)
			},
			wantStatus: http.StatusInternalServerError,
			wantMsg:    utils.MsgFailedToCreateTask,
		},
		{
			name:          "success",
			authenticated: true,
			body:          `{"title":"Learn Go","description":"Practice testing"}`,
			setupMock: func(service *mocks.MockTaskService) {
				service.EXPECT().
					CreateTask(
						1,
						"Learn Go",
						"Practice testing",
					).
					Return(&models.Task{
						ID:          1,
						Title:       "Learn Go",
						Description: "Practice testing",
						Status:      "pending",
						User: &models.User{
							ID: 1,
						},
					}, nil)
			},
			wantStatus: http.StatusCreated,
			wantMsg:    utils.MsgTaskCreated,
			wantTask: &models.Task{
				ID:          1,
				Title:       "Learn Go",
				Description: "Practice testing",
				Status:      "pending",
				User: &models.User{
					ID: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockTaskService(ctrl)

			restore := replaceTaskService(mockService)
			defer restore()

			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			ctx, rec := newTestContext(
				http.MethodPost,
				"/tasks",
				tt.body,
				tt.authenticated,
			)

			controller := &TaskController{}
			controller.Ctx = ctx

			controller.Create()

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}

			var response struct {
				Success bool         `json:"success"`
				Message string       `json:"message"`
				Data    *models.Task `json:"data"`
			}

			if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if response.Success != (tt.wantStatus < 400) {
				t.Fatalf("expected success %v, got %v", tt.wantStatus < 400, response.Success)
			}

			if response.Message != tt.wantMsg {
				t.Fatalf("expected message %q, got %q", tt.wantMsg, response.Message)
			}

			if tt.wantTask != nil {
				if !reflect.DeepEqual(response.Data, tt.wantTask) {
					t.Fatalf("expected %+v, got %+v", tt.wantTask, response.Data)
				}
			} else if response.Data != nil {
				t.Fatal("expected response data to be nil")
			}
		})
	}
}

func TestTaskController_GetAll(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		authenticated bool
		setupMock     func(*mocks.MockTaskService)
		wantStatus    int
		wantMsg       string
		wantTasks     []models.Task
	}{
		{
			name:          "unauthorized",
			url:           "/tasks",
			authenticated: false,
			wantStatus:    http.StatusUnauthorized,
			wantMsg:       utils.MsgUnauthorized,
		},
		{
			name:          "invalid page",
			url:           "/tasks?page=abc",
			authenticated: true,
			wantStatus:    http.StatusBadRequest,
			wantMsg:       utils.MsgInvalidRequestBody,
		},
		{
			name:          "page less than one",
			url:           "/tasks?page=0",
			authenticated: true,
			wantStatus:    http.StatusBadRequest,
			wantMsg:       utils.MsgInvalidRequestBody,
		},
		{
			name:          "invalid limit",
			url:           "/tasks?limit=abc",
			authenticated: true,
			wantStatus:    http.StatusBadRequest,
			wantMsg:       utils.MsgInvalidRequestBody,
		},
		{
			name:          "limit less than one",
			url:           "/tasks?limit=0",
			authenticated: true,
			wantStatus:    http.StatusBadRequest,
			wantMsg:       utils.MsgInvalidRequestBody,
		},
		{
			name:          "service error",
			url:           "/tasks?status=pending&page=2&limit=5",
			authenticated: true,
			setupMock: func(service *mocks.MockTaskService) {
				service.EXPECT().
					GetAllTasks(1, "pending", 2, 5).
					Return(nil, errDatabase)
			},
			wantStatus: http.StatusInternalServerError,
			wantMsg:    utils.MsgFailedToRetrieveTasks,
		},
		{
			name:          "success with defaults",
			url:           "/tasks",
			authenticated: true,
			setupMock: func(service *mocks.MockTaskService) {
				service.EXPECT().
					GetAllTasks(1, "", 1, 10).
					Return([]models.Task{
						{
							ID:     1,
							Title:  "Task 1",
							Status: "pending",
						},
					}, nil)
			},
			wantStatus: http.StatusOK,
			wantMsg:    utils.MsgTasksRetrieved,
			wantTasks: []models.Task{
				{
					ID:     1,
					Title:  "Task 1",
					Status: "pending",
				},
			},
		},
		{
			name:          "success with query parameters",
			url:           "/tasks?status=completed&page=2&limit=5",
			authenticated: true,
			setupMock: func(service *mocks.MockTaskService) {
				service.EXPECT().
					GetAllTasks(1, "completed", 2, 5).
					Return([]models.Task{
						{
							ID:     2,
							Title:  "Task 2",
							Status: "completed",
						},
					}, nil)
			},
			wantStatus: http.StatusOK,
			wantMsg:    utils.MsgTasksRetrieved,
			wantTasks: []models.Task{
				{
					ID:     2,
					Title:  "Task 2",
					Status: "completed",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockTaskService(ctrl)

			restore := replaceTaskService(mockService)
			defer restore()

			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			ctx, rec := newTestContext(
				http.MethodGet,
				tt.url,
				"",
				tt.authenticated,
			)

			controller := &TaskController{}
			controller.Ctx = ctx

			controller.GetAll()

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}

			var response struct {
				Success bool            `json:"success"`
				Message string          `json:"message"`
				Data    []models.Task   `json:"data"`
			}

			if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if response.Success != (tt.wantStatus < 400) {
				t.Fatalf("expected success %v, got %v", tt.wantStatus < 400, response.Success)
			}

			if response.Message != tt.wantMsg {
				t.Fatalf("expected message %q, got %q", tt.wantMsg, response.Message)
			}

			if tt.wantTasks != nil {
				if !reflect.DeepEqual(response.Data, tt.wantTasks) {
					t.Fatalf("expected %+v, got %+v", tt.wantTasks, response.Data)
				}
			} else if response.Data != nil {
				t.Fatal("expected response data to be nil")
			}
		})
	}
}

func TestTaskController_GetByID(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		authenticated bool
		setupMock     func(*mocks.MockTaskService)
		wantStatus    int
		wantMsg       string
		wantTask      *models.Task
	}{
		{
			name:          "unauthorized",
			authenticated: false,
			wantStatus:    http.StatusUnauthorized,
			wantMsg:       utils.MsgUnauthorized,
		},
		{
			name:          "invalid id",
			id:            "abc",
			authenticated: true,
			wantStatus:    http.StatusBadRequest,
			wantMsg:       utils.MsgInvalidTaskID,
		},
		{
			name:          "id less than one",
			id:            "0",
			authenticated: true,
			wantStatus:    http.StatusBadRequest,
			wantMsg:       utils.MsgInvalidTaskID,
		},
		{
			name:          "service error",
			id:            "1",
			authenticated: true,
			setupMock: func(service *mocks.MockTaskService) {
				service.EXPECT().
					GetTaskByID(1, 1).
					Return(nil, errDatabase)
			},
			wantStatus: http.StatusInternalServerError,
			wantMsg:    utils.MsgFailedToRetrieveTask,
		},
		{
			name:          "task not found",
			id:            "1",
			authenticated: true,
			setupMock: func(service *mocks.MockTaskService) {
				service.EXPECT().
					GetTaskByID(1, 1).
					Return(nil, nil)
			},
			wantStatus: http.StatusNotFound,
			wantMsg:    utils.MsgTaskNotFound,
		},
		{
			name:          "success",
			id:            "1",
			authenticated: true,
			setupMock: func(service *mocks.MockTaskService) {
				service.EXPECT().
					GetTaskByID(1, 1).
					Return(&models.Task{
						ID:          1,
						Title:       "Learn Go",
						Description: "Practice testing",
						Status:      "pending",
						User: &models.User{
							ID: 1,
						},
					}, nil)
			},
			wantStatus: http.StatusOK,
			wantMsg:    utils.MsgTaskRetrieved,
			wantTask: &models.Task{
				ID:          1,
				Title:       "Learn Go",
				Description: "Practice testing",
				Status:      "pending",
				User: &models.User{
					ID: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockTaskService(ctrl)

			restore := replaceTaskService(mockService)
			defer restore()

			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			ctx, rec := newTestContext(
				http.MethodGet,
				"/tasks/"+tt.id,
				"",
				tt.authenticated,
			)

			ctx.Input.SetParam(":id", tt.id)

			controller := &TaskController{}
			controller.Ctx = ctx

			controller.GetByID()

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}

			var response struct {
				Success bool         `json:"success"`
				Message string       `json:"message"`
				Data    *models.Task `json:"data"`
			}

			if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if response.Success != (tt.wantStatus < 400) {
				t.Fatalf("expected success %v, got %v", tt.wantStatus < 400, response.Success)
			}

			if response.Message != tt.wantMsg {
				t.Fatalf("expected message %q, got %q", tt.wantMsg, response.Message)
			}

			if tt.wantTask != nil {
				if !reflect.DeepEqual(response.Data, tt.wantTask) {
					t.Fatalf("expected %+v, got %+v", tt.wantTask, response.Data)
				}
			} else if response.Data != nil {
				t.Fatal("expected response data to be nil")
			}
		})
	}
}

func TestTaskController_Update(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		body          string
		authenticated bool
		setupMock     func(*mocks.MockTaskService)
		wantStatus    int
		wantMsg       string
		wantTask      *models.Task
	}{
		{
			name:          "unauthorized",
			authenticated: false,
			wantStatus:    http.StatusUnauthorized,
			wantMsg:       utils.MsgUnauthorized,
		},
		{
			name:          "invalid id",
			id:            "abc",
			authenticated: true,
			wantStatus:    http.StatusBadRequest,
			wantMsg:       utils.MsgInvalidTaskID,
		},
		{
			name:          "id less than one",
			id:            "0",
			authenticated: true,
			wantStatus:    http.StatusBadRequest,
			wantMsg:       utils.MsgInvalidTaskID,
		},
		{
			name:          "invalid request body",
			id:            "1",
			authenticated: true,
			body:          `{"title":`,
			wantStatus:    http.StatusBadRequest,
			wantMsg:       utils.MsgInvalidRequestBody,
		},
		{
			name:          "no fields provided",
			id:            "1",
			authenticated: true,
			body:          `{}`,
			wantStatus:    http.StatusBadRequest,
			wantMsg:       utils.MsgAtLeastOneFieldRequired,
		},
		{
			name:          "empty title",
			id:            "1",
			authenticated: true,
			body:          `{"title":"   "}`,
			wantStatus:    http.StatusBadRequest,
			wantMsg:       utils.MsgTitleRequired,
		},
		{
			name:          "empty status",
			id:            "1",
			authenticated: true,
			body:          `{"status":"   "}`,
			wantStatus:    http.StatusBadRequest,
			wantMsg:       utils.MsgStatusRequired,
		},
		{
			name:          "service error",
			id:            "1",
			authenticated: true,
			body:          `{"title":"Updated Title"}`,
			setupMock: func(service *mocks.MockTaskService) {
				title := "Updated Title"

				service.EXPECT().
					UpdateTask(1, 1, &title, nil, nil).
					Return(nil, errDatabase)
			},
			wantStatus: http.StatusInternalServerError,
			wantMsg:    utils.MsgFailedToUpdateTask,
		},
		{
			name:          "task not found",
			id:            "1",
			authenticated: true,
			body:          `{"title":"Updated Title"}`,
			setupMock: func(service *mocks.MockTaskService) {
				title := "Updated Title"

				service.EXPECT().
					UpdateTask(1, 1, &title, nil, nil).
					Return(nil, nil)
			},
			wantStatus: http.StatusNotFound,
			wantMsg:    utils.MsgTaskNotFound,
		},
		{
			name:          "success",
			id:            "1",
			authenticated: true,
			body:          `{"title":"Updated Title","status":"completed"}`,
			setupMock: func(service *mocks.MockTaskService) {
				title := "Updated Title"
				status := "completed"

				service.EXPECT().
					UpdateTask(1, 1, &title, nil, &status).
					Return(&models.Task{
						ID:     1,
						Title:  "Updated Title",
						Status: "completed",
						User: &models.User{
							ID: 1,
						},
					}, nil)
			},
			wantStatus: http.StatusOK,
			wantMsg:    utils.MsgTaskUpdated,
			wantTask: &models.Task{
				ID:     1,
				Title:  "Updated Title",
				Status: "completed",
				User: &models.User{
					ID: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockTaskService(ctrl)

			restore := replaceTaskService(mockService)
			defer restore()

			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			ctx, rec := newTestContext(
				http.MethodPut,
				"/tasks/"+tt.id,
				tt.body,
				tt.authenticated,
			)

			ctx.Input.SetParam(":id", tt.id)

			controller := &TaskController{}
			controller.Ctx = ctx

			controller.Update()

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}

			var response struct {
				Success bool         `json:"success"`
				Message string       `json:"message"`
				Data    *models.Task `json:"data"`
			}

			if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if response.Success != (tt.wantStatus < 400) {
				t.Fatalf("expected success %v, got %v", tt.wantStatus < 400, response.Success)
			}

			if response.Message != tt.wantMsg {
				t.Fatalf("expected message %q, got %q", tt.wantMsg, response.Message)
			}

			if tt.wantTask != nil {
				if !reflect.DeepEqual(response.Data, tt.wantTask) {
					t.Fatalf("expected %+v, got %+v", tt.wantTask, response.Data)
				}
			} else if response.Data != nil {
				t.Fatal("expected response data to be nil")
			}
		})
	}
}

func TestTaskController_Delete(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		authenticated bool
		setupMock     func(*mocks.MockTaskService)
		wantStatus    int
		wantMsg       string
	}{
		{
			name:          "unauthorized",
			authenticated: false,
			wantStatus:    http.StatusUnauthorized,
			wantMsg:       utils.MsgUnauthorized,
		},
		{
			name:          "invalid id",
			id:            "abc",
			authenticated: true,
			wantStatus:    http.StatusBadRequest,
			wantMsg:       utils.MsgInvalidTaskID,
		},
		{
			name:          "id less than one",
			id:            "0",
			authenticated: true,
			wantStatus:    http.StatusBadRequest,
			wantMsg:       utils.MsgInvalidTaskID,
		},
		{
			name:          "service error",
			id:            "1",
			authenticated: true,
			setupMock: func(service *mocks.MockTaskService) {
				service.EXPECT().
					DeleteTask(1, 1).
					Return(false, errDatabase)
			},
			wantStatus: http.StatusInternalServerError,
			wantMsg:    utils.MsgFailedToDeleteTask,
		},
		{
			name:          "task not found",
			id:            "1",
			authenticated: true,
			setupMock: func(service *mocks.MockTaskService) {
				service.EXPECT().
					DeleteTask(1, 1).
					Return(false, nil)
			},
			wantStatus: http.StatusNotFound,
			wantMsg:    utils.MsgTaskNotFound,
		},
		{
			name:          "success",
			id:            "1",
			authenticated: true,
			setupMock: func(service *mocks.MockTaskService) {
				service.EXPECT().
					DeleteTask(1, 1).
					Return(true, nil)
			},
			wantStatus: http.StatusOK,
			wantMsg:    utils.MsgTaskDeleted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockTaskService(ctrl)

			restore := replaceTaskService(mockService)
			defer restore()

			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			ctx, rec := newTestContext(
				http.MethodDelete,
				"/tasks/"+tt.id,
				"",
				tt.authenticated,
			)

			ctx.Input.SetParam(":id", tt.id)

			controller := &TaskController{}
			controller.Ctx = ctx

			controller.Delete()

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}

			var response struct {
				Success bool        `json:"success"`
				Message string      `json:"message"`
				Data    any `json:"data"`
			}

			if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if response.Success != (tt.wantStatus < 400) {
				t.Fatalf("expected success %v, got %v", tt.wantStatus < 400, response.Success)
			}

			if response.Message != tt.wantMsg {
				t.Fatalf("expected message %q, got %q", tt.wantMsg, response.Message)
			}

			if response.Data != nil {
				t.Fatal("expected response data to be nil")
			}
		})
	}
}


// Helper functions for testing
func replaceTaskService(mockTaskService services.TaskService) func() {
    originalTaskService := taskService
    taskService = mockTaskService

    return func() {
        taskService = originalTaskService
    }
}

func newTestContext(method, url, body string, authenticated bool) (*context.Context, *httptest.ResponseRecorder) {
    req := httptest.NewRequest(method, url, strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    rec := httptest.NewRecorder()

    ctx := context.NewContext()
    ctx.Reset(rec, req)
    ctx.Input.RequestBody = []byte(body)

    if authenticated {
        ctx.Input.SetData("userID", 1)
    }

    return ctx, rec
}
