package middlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"todo-api/utils"

	beegoCtx "github.com/beego/beego/v2/server/web/context"
)

func TestAuthFilter(t *testing.T) {
	validToken, err := utils.GenerateToken(1)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	tests := []struct {
		name         string
		authHeader   string
		wantStatus   int
		wantMessage  string
		expectUserID bool
	}{
		{
			name:        "missing authorization header",
			authHeader:  "",
			wantStatus:  http.StatusUnauthorized,
			wantMessage: "Authorization header is required",
		},
		{
			name:        "invalid authorization format",
			authHeader:  "invalid",
			wantStatus:  http.StatusUnauthorized,
			wantMessage: "Authorization header must be Bearer {token}",
		},
		{
			name:        "wrong authorization scheme",
			authHeader:  "Basic token",
			wantStatus:  http.StatusUnauthorized,
			wantMessage: "Authorization header must be Bearer {token}",
		},
		{
			name:        "invalid token",
			authHeader:  "Bearer invalid-token",
			wantStatus:  http.StatusUnauthorized,
			wantMessage: "Invalid or expired token",
		},
		{
			name:         "valid token",
			authHeader:   "Bearer " + validToken,
			wantStatus:   http.StatusOK,
			expectUserID: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, rec := newMiddlewareContext(tt.authHeader)

			AuthFilter(ctx)

			if tt.expectUserID {
				userID, ok := ctx.Input.GetData("userID").(int)
				if !ok {
					t.Fatal("expected userID in context")
				}

				if userID != 1 {
					t.Fatalf("expected userID %d, got %d", 1, userID)
				}

				if rec.Code != http.StatusOK {
					t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
				}

				if strings.TrimSpace(rec.Body.String()) != "" {
					t.Fatalf("expected empty response body, got %q", rec.Body.String())
				}

				return
			}

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}

			var response struct {
				Success bool   `json:"success"`
				Message string `json:"message"`
			}

			if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if response.Success {
				t.Fatal("expected success to be false")
			}

			if response.Message != tt.wantMessage {
				t.Fatalf("expected message %q, got %q", tt.wantMessage, response.Message)
			}
		})
	}
}

func newMiddlewareContext(authHeader string) (*beegoCtx.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	rec := httptest.NewRecorder()

	ctx := beegoCtx.NewContext()
	ctx.Reset(rec, req)

	return ctx, rec
}
