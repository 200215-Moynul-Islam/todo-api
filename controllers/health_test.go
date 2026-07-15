package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo-api/utils"

	"github.com/beego/beego/v2/server/web/context"
)

type responseBody struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func TestHealthController_Get(t *testing.T) {
	// Arrange
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	ctx := context.NewContext()
	ctx.Reset(rec, req)

	controller := &HealthController{}
	controller.Ctx = ctx

	// Act
	controller.Get()

	// Assert
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response responseBody
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !response.Success {
		t.Fatal("expected success to be true")
	}

	if response.Message != utils.MsgServerRunning {
		t.Fatalf(
			"expected message %q, got %q",
			utils.MsgServerRunning,
			response.Message,
		)
	}
}
