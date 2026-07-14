package services

import (
	"errors"
	"testing"

	"todo-api/mocks"
	"todo-api/models"

	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_RegisterUser_Validation(t *testing.T) {
	tests := []struct {
		name     string
		userName string
		email    string
		password string
		wantErr  error
	}{
		{
			name:     "empty name",
			userName: "",
			email:    "john@example.com",
			password: "password",
			wantErr:  ErrNameRequired,
		},
		{
			name:     "empty email",
			userName: "John",
			email:    "",
			password: "password",
			wantErr:  ErrEmailRequired,
		},
		{
			name:     "empty password",
			userName: "John",
			email:    "john@example.com",
			password: "",
			wantErr:  ErrPasswordRequired,
		},
		{
			name:     "short password",
			userName: "John",
			email:    "john@example.com",
			password: "12345",
			wantErr:  ErrPasswordTooShort,
		},
		{
			name:     "invalid email",
			userName: "John",
			email:    "invalid-email",
			password: "password",
			wantErr:  ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockUserRepository(ctrl)

			service := NewUserService(mockRepo)

			user, err := service.RegisterUser(
				tt.userName,
				tt.email,
				tt.password,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if user != nil {
				t.Fatal("expected nil user")
			}
		})
	}
}

func TestUserService_RegisterUser_GetByEmailError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	expectedErr := errors.New("database error")

	mockRepo.EXPECT().
		GetByEmail("john@example.com").
		Return(nil, expectedErr)

	user, err := service.RegisterUser(
		"John",
		"john@example.com",
		"password123",
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if user != nil {
		t.Fatal("expected nil user")
	}
}

func TestUserService_RegisterUser_EmailAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	existingUser := &models.User{
		ID:    1,
		Name:  "John",
		Email: "john@example.com",
	}

	mockRepo.EXPECT().
		GetByEmail("john@example.com").
		Return(existingUser, nil)

	user, err := service.RegisterUser(
		"John",
		"john@example.com",
		"password123",
	)

	if !errors.Is(err, ErrEmailExists) {
		t.Fatalf("expected error %v, got %v", ErrEmailExists, err)
	}

	if user != nil {
		t.Fatal("expected nil user")
	}
}

func TestUserService_RegisterUser_CreateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	expectedErr := errors.New("failed to create user")

	mockRepo.EXPECT().
		GetByEmail("john@example.com").
		Return(nil, nil)

	mockRepo.EXPECT().
		Create(gomock.Any()).
		Return(expectedErr)

	user, err := service.RegisterUser(
		"John",
		"john@example.com",
		"password123",
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if user != nil {
		t.Fatal("expected nil user")
	}
}

func TestUserService_RegisterUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	inputName := "  John Doe  "
	inputEmail := "  John@Example.COM  "
	inputPassword := "password123"

	mockRepo.EXPECT().
		GetByEmail("john@example.com").
		Return(nil, nil)

	mockRepo.EXPECT().
		Create(gomock.AssignableToTypeOf(&models.User{})).
		DoAndReturn(func(user *models.User) error {

			if user.Name != "John Doe" {
				t.Errorf("expected name %q, got %q", "John Doe", user.Name)
			}

			if user.Email != "john@example.com" {
				t.Errorf("expected email %q, got %q", "john@example.com", user.Email)
			}

			if user.Password == inputPassword {
				t.Error("password should be hashed")
			}

			if err := bcrypt.CompareHashAndPassword(
				[]byte(user.Password),
				[]byte(inputPassword),
			); err != nil {
				t.Error("password hash does not match input password")
			}

			return nil
		})

	user, err := service.RegisterUser(
		inputName,
		inputEmail,
		inputPassword,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("expected user, got nil")
	}

	if user.Name != "John Doe" {
		t.Errorf("expected name %q, got %q", "John Doe", user.Name)
	}

	if user.Email != "john@example.com" {
		t.Errorf("expected email %q, got %q", "john@example.com", user.Email)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(inputPassword),
	); err != nil {
		t.Error("returned password hash is invalid")
	}
}
