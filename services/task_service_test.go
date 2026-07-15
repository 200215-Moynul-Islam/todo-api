package services

import (
	"errors"
	"reflect"
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

func TestTaskService_GetAllTasks(t *testing.T) {
	errDatabase := errors.New("database error")

	expectedTasks := []models.Task{
		{
			ID:     1,
			Title:  "Task 1",
			Status: "pending",
		},
		{
			ID:     2,
			Title:  "Task 2",
			Status: "completed",
		},
	}

	tests := []struct {
		name       string
		status      string
		setupMock   func(*mocks.MockTaskRepository)
		wantTasks   []models.Task
		wantErr     error
	}{
		{
			name:   "repository error",
			status: " pending ",
			setupMock: func(repo *mocks.MockTaskRepository) {
				repo.EXPECT().
					GetAll(1, "pending", 1, 10).
					Return(nil, errDatabase)
			},
			wantErr: errDatabase,
		},
		{
			name:   "success",
			status: " pending ",
			setupMock: func(repo *mocks.MockTaskRepository) {
				repo.EXPECT().
					GetAll(1, "pending", 1, 10).
					Return(expectedTasks, nil)
			},
			wantTasks: expectedTasks,
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

			tasks, err := service.GetAllTasks(1, tt.status, 1, 10)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(tasks, tt.wantTasks) {
				t.Fatalf("expected %+v, got %+v", tt.wantTasks, tasks)
			}
		})
	}
}

func TestTaskService_GetTaskByID(t *testing.T) {
	expectedTask := &models.Task{
		ID:          1,
		Title:       "Learn Go",
		Description: "Practice testing",
		Status:      "pending",
		User: &models.User{
			ID: 1,
		},
	}

	tests := []struct {
		name      string
		taskID    int
		userID    int
		setupMock func(*mocks.MockTaskRepository)
		wantTask  *models.Task
		wantErr   error
	}{
		{
			name:   "repository error",
			taskID: 1,
			userID: 1,
			setupMock: func(repo *mocks.MockTaskRepository) {
				repo.EXPECT().
					GetByID(1).
					Return(nil, errDatabase)
			},
			wantErr: errDatabase,
		},
		{
			name:   "task not found",
			taskID: 1,
			userID: 1,
			setupMock: func(repo *mocks.MockTaskRepository) {
				repo.EXPECT().
					GetByID(1).
					Return(nil, nil)
			},
			wantTask: nil,
		},
		{
			name:   "task belongs to another user",
			taskID: 1,
			userID: 2,
			setupMock: func(repo *mocks.MockTaskRepository) {
				repo.EXPECT().
					GetByID(1).
					Return(expectedTask, nil)
			},
			wantTask: nil,
		},
		{
			name:   "success",
			taskID: 1,
			userID: 1,
			setupMock: func(repo *mocks.MockTaskRepository) {
				repo.EXPECT().
					GetByID(1).
					Return(expectedTask, nil)
			},
			wantTask: expectedTask,
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

			task, err := service.GetTaskByID(tt.userID, tt.taskID)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(task, tt.wantTask) {
				t.Fatalf("expected %+v, got %+v", tt.wantTask, task)
			}
		})
	}
}

func TestTaskService_UpdateTask(t *testing.T) {
	newTitle := "  Learn Go Testing  "
	newDescription := "  Updated description  "
	newStatus := " completed "

	tests := []struct {
		name      string
		userID    int
		taskID    int
		title     *string
		desc      *string
		status    *string
		setupMock func(*mocks.MockTaskRepository)
		wantTask  *models.Task
		wantErr   error
	}{
		{
			name:   "repository get error",
			userID: 1,
			taskID: 1,
			setupMock: func(repo *mocks.MockTaskRepository) {
				repo.EXPECT().
					GetByID(1).
					Return(nil, errDatabase)
			},
			wantErr: errDatabase,
		},
		{
			name:   "task not found",
			userID: 1,
			taskID: 1,
			setupMock: func(repo *mocks.MockTaskRepository) {
				repo.EXPECT().
					GetByID(1).
					Return(nil, nil)
			},
		},
		{
			name:   "task belongs to another user",
			userID: 2,
			taskID: 1,
			setupMock: func(repo *mocks.MockTaskRepository) {
				repo.EXPECT().
					GetByID(1).
					Return(&models.Task{
						ID:     1,
						Title:  "Old",
						Status: "pending",
						User:   &models.User{ID: 1},
					}, nil)
			},
		},
		{
			name:   "repository update error",
			userID: 1,
			taskID: 1,
			title:  &newTitle,
			setupMock: func(repo *mocks.MockTaskRepository) {
				task := &models.Task{
					ID:     1,
					Title:  "Old",
					Status: "pending",
					User:   &models.User{ID: 1},
				}

				repo.EXPECT().
					GetByID(1).
					Return(task, nil)

				repo.EXPECT().
					Update(gomock.Any()).
					Return(errDatabase)
			},
			wantErr: errDatabase,
		},
		{
			name:   "update all fields successfully",
			userID: 1,
			taskID: 1,
			title:  &newTitle,
			desc:   &newDescription,
			status: &newStatus,
			setupMock: func(repo *mocks.MockTaskRepository) {
				task := &models.Task{
					ID:          1,
					Title:       "Old",
					Description: "Old Description",
					Status:      "pending",
					User:        &models.User{ID: 1},
				}

				repo.EXPECT().
					GetByID(1).
					Return(task, nil)

				repo.EXPECT().
					Update(gomock.AssignableToTypeOf(&models.Task{})).
					DoAndReturn(func(task *models.Task) error {
						if task.Title != "Learn Go Testing" {
							t.Errorf("expected title %q, got %q", "Learn Go Testing", task.Title)
						}

						if task.Description != "Updated description" {
							t.Errorf("expected description %q, got %q", "Updated description", task.Description)
						}

						if task.Status != "completed" {
							t.Errorf("expected status %q, got %q", "completed", task.Status)
						}

						return nil
					})
			},
			wantTask: &models.Task{
				ID:          1,
				Title:       "Learn Go Testing",
				Description: "Updated description",
				Status:      "completed",
				User:        &models.User{ID: 1},
			},
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

			task, err := service.UpdateTask(
				tt.userID,
				tt.taskID,
				tt.title,
				tt.desc,
				tt.status,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(task, tt.wantTask) {
				t.Fatalf("expected %+v, got %+v", tt.wantTask, task)
			}
		})
	}
}
