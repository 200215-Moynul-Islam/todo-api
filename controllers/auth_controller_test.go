package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/beego/beego/v2/server/web/context"
	"go.uber.org/mock/gomock"

	"todo-api/models"
	"todo-api/mocks"
	"todo-api/services"
	"todo-api/utils"
)

type registerResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    *models.User `json:"data"`
}

func TestAuthController_Register(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		setupMock  func(*mocks.MockUserService)
		wantStatus int
		wantMsg    string
		wantUser   bool
	}{
		{
			name:       "invalid request body",
			body:       `{"name":`,
			setupMock:  nil,
			wantStatus: http.StatusBadRequest,
			wantMsg:    utils.MsgInvalidRequestBody,
			wantUser:   false,
		},
		{
			name: "name required",
			body: `{
				"name":"",
				"email":"john@example.com",
				"password":"password123"
			}`,
			setupMock: func(service *mocks.MockUserService) {
				service.EXPECT().
					RegisterUser("", "john@example.com", "password123").
					Return(nil, services.ErrNameRequired)
			},
			wantStatus: http.StatusBadRequest,
			wantMsg:    services.ErrNameRequired.Error(),
			wantUser:   false,
		},
		{
			name: "email required",
			body: `{
				"name":"John",
				"email":"",
				"password":"password123"
			}`,
			setupMock: func(service *mocks.MockUserService) {
				service.EXPECT().
					RegisterUser("John", "", "password123").
					Return(nil, services.ErrEmailRequired)
			},
			wantStatus: http.StatusBadRequest,
			wantMsg:    services.ErrEmailRequired.Error(),
			wantUser:   false,
		},
		{
			name: "password required",
			body: `{
				"name":"John",
				"email":"john@example.com",
				"password":""
			}`,
			setupMock: func(service *mocks.MockUserService) {
				service.EXPECT().
					RegisterUser("John", "john@example.com", "").
					Return(nil, services.ErrPasswordRequired)
			},
			wantStatus: http.StatusBadRequest,
			wantMsg:    services.ErrPasswordRequired.Error(),
			wantUser:   false,
		},
		{
			name: "password too short",
			body: `{
				"name":"John",
				"email":"john@example.com",
				"password":"12345"
			}`,
			setupMock: func(service *mocks.MockUserService) {
				service.EXPECT().
					RegisterUser("John", "john@example.com", "12345").
					Return(nil, services.ErrPasswordTooShort)
			},
			wantStatus: http.StatusBadRequest,
			wantMsg:    services.ErrPasswordTooShort.Error(),
			wantUser:   false,
		},
		{
			name: "invalid email",
			body: `{
				"name":"John",
				"email":"invalid-email",
				"password":"password123"
			}`,
			setupMock: func(service *mocks.MockUserService) {
				service.EXPECT().
					RegisterUser("John", "invalid-email", "password123").
					Return(nil, services.ErrInvalidEmail)
			},
			wantStatus: http.StatusBadRequest,
			wantMsg:    services.ErrInvalidEmail.Error(),
			wantUser:   false,
		},
		{
			name: "email already exists",
			body: `{
				"name":"John",
				"email":"john@example.com",
				"password":"password123"
			}`,
			setupMock: func(service *mocks.MockUserService) {
				service.EXPECT().
					RegisterUser("John", "john@example.com", "password123").
					Return(nil, services.ErrEmailExists)
			},
			wantStatus: http.StatusConflict,
			wantMsg:    services.ErrEmailExists.Error(),
			wantUser:   false,
		},
		{
			name: "internal server error",
			body: `{
				"name":"John",
				"email":"john@example.com",
				"password":"password123"
			}`,
			setupMock: func(service *mocks.MockUserService) {
				service.EXPECT().
					RegisterUser("John", "john@example.com", "password123").
					Return(nil, errors.New("database error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantMsg:    utils.MsgFailedToRegisterUser,
			wantUser:   false,
		},
		{
			name: "success",
			body: `{
				"name":"John Doe",
				"email":"john@example.com",
				"password":"password123"
			}`,
			setupMock: func(service *mocks.MockUserService) {
				service.EXPECT().
					RegisterUser("John Doe", "john@example.com", "password123").
					Return(&models.User{
						ID:    1,
						Name:  "John Doe",
						Email: "john@example.com",
					}, nil)
			},
			wantStatus: http.StatusCreated,
			wantMsg:    utils.MsgUserRegistered,
			wantUser:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockUserService(ctrl)

			originalService := userService
			userService = mockService
			defer func() {
				userService = originalService
			}()

			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			req := httptest.NewRequest(
				http.MethodPost,
				"/register",
				strings.NewReader(tt.body),
			)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			ctx := context.NewContext()
			ctx.Reset(rec, req)
			ctx.Input.RequestBody = []byte(tt.body)

			controller := &AuthController{}
			controller.Ctx = ctx

			controller.Register()

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}

			var response registerResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if response.Success != (tt.wantStatus < 400) {
				t.Fatalf("unexpected success value: %v", response.Success)
			}

			if response.Message != tt.wantMsg {
				t.Fatalf("expected message %q, got %q", tt.wantMsg, response.Message)
			}

			if tt.wantUser {
				if response.Data == nil {
					t.Fatal("expected user in response")
				}

				if response.Data.Name != "John Doe" {
					t.Fatalf("expected user name %q, got %q", "John Doe", response.Data.Name)
				}

				if response.Data.Email != "john@example.com" {
					t.Fatalf("expected email %q, got %q", "john@example.com", response.Data.Email)
				}
			} else if response.Data != nil {
				t.Fatal("expected response data to be nil")
			}
		})
	}
}

func TestAuthController_Login(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		setupMock  func(*mocks.MockUserService)
		wantStatus int
		wantMsg    string
		wantToken  string
	}{
		{
			name:       "invalid request body",
			body:       `{"email":`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    utils.MsgInvalidRequestBody,
		},
		{
			name: "email required",
			body: `{
				"email":"",
				"password":"password123"
			}`,
			setupMock: func(service *mocks.MockUserService) {
				service.EXPECT().
					LoginUser("", "password123").
					Return("", services.ErrEmailRequired)
			},
			wantStatus: http.StatusBadRequest,
			wantMsg:    services.ErrEmailRequired.Error(),
		},
		{
			name: "password required",
			body: `{
				"email":"john@example.com",
				"password":""
			}`,
			setupMock: func(service *mocks.MockUserService) {
				service.EXPECT().
					LoginUser("john@example.com", "").
					Return("", services.ErrPasswordRequired)
			},
			wantStatus: http.StatusBadRequest,
			wantMsg:    services.ErrPasswordRequired.Error(),
		},
		{
			name: "invalid credentials",
			body: `{
				"email":"john@example.com",
				"password":"wrong-password"
			}`,
			setupMock: func(service *mocks.MockUserService) {
				service.EXPECT().
					LoginUser("john@example.com", "wrong-password").
					Return("", services.ErrInvalidCredentials)
			},
			wantStatus: http.StatusUnauthorized,
			wantMsg:    services.ErrInvalidCredentials.Error(),
		},
		{
			name: "internal server error",
			body: `{
				"email":"john@example.com",
				"password":"password123"
			}`,
			setupMock: func(service *mocks.MockUserService) {
				service.EXPECT().
					LoginUser("john@example.com", "password123").
					Return("", errors.New("database error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantMsg:    utils.MsgFailedToLogin,
		},
		{
			name: "success",
			body: `{
				"email":"john@example.com",
				"password":"password123"
			}`,
			setupMock: func(service *mocks.MockUserService) {
				service.EXPECT().
					LoginUser("john@example.com", "password123").
					Return("jwt-token", nil)
			},
			wantStatus: http.StatusOK,
			wantMsg:    utils.MsgLoginSuccessful,
			wantToken:  "jwt-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockUserService(ctrl)

			originalService := userService
			userService = mockService
			defer func() {
				userService = originalService
			}()

			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			req := httptest.NewRequest(
				http.MethodPost,
				"/login",
				strings.NewReader(tt.body),
			)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			ctx := context.NewContext()
			ctx.Reset(rec, req)
			ctx.Input.RequestBody = []byte(tt.body)

			controller := &AuthController{}
			controller.Ctx = ctx

			controller.Login()

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}

			var response struct {
				Success bool `json:"success"`
				Message string `json:"message"`
				Data    *LoginResponse `json:"data"`
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

			if tt.wantToken != "" {
				if response.Data == nil {
					t.Fatal("expected login response")
				}

				if response.Data.Token != tt.wantToken {
					t.Fatalf("expected token %q, got %q", tt.wantToken, response.Data.Token)
				}
			} else if response.Data != nil {
				t.Fatal("expected response data to be nil")
			}
		})
	}
}
