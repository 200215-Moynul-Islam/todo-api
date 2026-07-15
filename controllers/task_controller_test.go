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
	"todo-api/utils"

	"github.com/beego/beego/v2/server/web/context"
	"go.uber.org/mock/gomock"
)

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

			originalService := taskService
			taskService = mockService
			defer func() {
				taskService = originalService
			}()

			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			req := httptest.NewRequest(
				http.MethodPost,
				"/tasks/generate-description",
				strings.NewReader(tt.body),
			)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			ctx := context.NewContext()
			ctx.Reset(rec, req)
			ctx.Input.RequestBody = []byte(tt.body)

			if tt.authenticated {
				ctx.Input.SetData("userID", 1)
			}

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
					Return(nil, errors.New("database error"))
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

			originalService := taskService
			taskService = mockService
			defer func() {
				taskService = originalService
			}()

			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			req := httptest.NewRequest(
				http.MethodPost,
				"/tasks",
				strings.NewReader(tt.body),
			)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			ctx := context.NewContext()
			ctx.Reset(rec, req)
			ctx.Input.RequestBody = []byte(tt.body)

			if tt.authenticated {
				ctx.Input.SetData("userID", 1)
			}

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
