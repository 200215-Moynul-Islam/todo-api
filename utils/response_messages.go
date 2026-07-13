package utils

const (
	// Common messages
	MsgUnauthorized       = "Unauthorized"
	MsgInvalidRequestBody = "Invalid request body"
	MsgServerRunning      = "Server is running"

	// Task messages
	MsgTitleRequired               = "Title is required"
	MsgStatusRequired              = "Status is required"
	MsgAtLeastOneFieldRequired     = "At least one field (title, description, status) is required"
	MsgTaskNotFound                = "Task not found"
	MsgInvalidTaskID               = "Invalid task id format"
	MsgTaskCreated                 = "Task created successfully"
	MsgFailedToCreateTask          = "Failed to create task"
	MsgTaskUpdated                 = "Task updated successfully"
	MsgFailedToUpdateTask          = "Failed to update task"
	MsgTaskDeleted                 = "Task deleted successfully"
	MsgFailedToDeleteTask          = "Failed to delete task"
	MsgTaskRetrieved               = "Task retrieved successfully"
	MsgFailedToRetrieveTask        = "Failed to retrieve task"
	MsgTasksRetrieved              = "Tasks retrieved successfully"
	MsgFailedToRetrieveTasks       = "Failed to retrieve tasks"
	MsgDescriptionGenerated        = "Description generated successfully"
	MsgFailedToGenerateDescription = "Failed to generate description"

	// Auth messages
	MsgUserRegistered       = "User registered successfully"
	MsgFailedToRegisterUser = "Failed to register user"
	MsgLoginSuccessful      = "Login successful"
	MsgFailedToLogin        = "Failed to log in"
)
