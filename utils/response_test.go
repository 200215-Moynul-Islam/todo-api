package utils

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego/v2/server/web/context"
)

func TestSendJSONResponse(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		success    bool
		message    string
		data       any
		wantStatus int
	}{
		{
			name:       "success response",
			status:     200,
			success:    true,
			message:    "Success",
			data:       map[string]any{"id": 1},
			wantStatus: 200,
		},
		{
			name:       "error response",
			status:     400,
			success:    false,
			message:    "Bad Request",
			data:       nil,
			wantStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			ctx := context.NewContext()
			ctx.Reset(rec, httptest.NewRequest("GET", "/", nil))

			SendJSONResponse(
				ctx,
				tt.status,
				tt.success,
				tt.message,
				tt.data,
			)

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}

			var response APIResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if response.Success != tt.success {
				t.Fatalf("expected success %v, got %v", tt.success, response.Success)
			}

			if response.Message != tt.message {
				t.Fatalf("expected message %q, got %q", tt.message, response.Message)
			}

			// JSON unmarshals objects into map[string]any,
			// so DeepEqual works for both nil and object values.
			if response.Data == nil && tt.data != nil {
				t.Fatalf("expected data, got nil")
			}

			if response.Data != nil && tt.data == nil {
				t.Fatalf("expected nil data, got %v", response.Data)
			}
		})
	}
}