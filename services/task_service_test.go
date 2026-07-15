package services

import (
	"errors"
	"testing"
	"todo-api/mocks"
	"todo-api/models"

	"go.uber.org/mock/gomock"
)

var errDatabase = errors.New("database error")

func TestTaskService_CreateTask(t *testing.T) {
	tests := []struct {
		name      string
		userID    int
		title     string
		desc      string
		setupMock func(*mocks.MockTaskRepository)
		wantErr   error
	}{
		{
			name:   "repository create error",
			userID: 1,
			title:  "  Learn Go  ",
			desc:   "  Practice testing  ",
			setupMock: func(repo *mocks.MockTaskRepository) {
				repo.EXPECT().
					Create(gomock.Any()).
					Return(errDatabase)
			},
			wantErr: errDatabase,
		},
		{
			name:   "success",
			userID: 1,
			title:  "  Learn Go  ",
			desc:   "  Practice testing  ",
			setupMock: func(repo *mocks.MockTaskRepository) {
				repo.EXPECT().
					Create(gomock.AssignableToTypeOf(&models.Task{})).
					DoAndReturn(func(task *models.Task) error {

						if task.Title != "Learn Go" {
							t.Errorf("expected title %q, got %q", "Learn Go", task.Title)
						}

						if task.Description != "Practice testing" {
							t.Errorf("expected description %q, got %q", "Practice testing", task.Description)
						}

						if task.Status != "pending" {
							t.Errorf("expected status %q, got %q", "pending", task.Status)
						}

						if task.User == nil {
							t.Fatal("expected user")
						}

						if task.User.ID != 1 {
							t.Errorf("expected user ID %d, got %d", 1, task.User.ID)
						}

						return nil
					})
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockTaskRepository(ctrl)
			mockGemini := mocks.NewMockGeminiClient(ctrl)

			tt.setupMock(mockRepo)

			service := NewTaskService(mockRepo, mockGemini)

			task, err := service.CreateTask(
				tt.userID,
				tt.title,
				tt.desc,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if tt.wantErr != nil {
				if task != nil {
					t.Fatal("expected nil task")
				}
				return
			}

			if task == nil {
				t.Fatal("expected task")
			}
		})
	}
}
