package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego/v2/server/web/context"
)

func TestBaseController_GetUserID(t *testing.T) {
	tests := []struct {
		name       string
		input      any
		wantUserID int
		wantOK     bool
	}{
		{
			name:       "user ID exists",
			input:      1,
			wantUserID: 1,
			wantOK:     true,
		},
		{
			name:       "user ID not found",
			input:      nil,
			wantUserID: 0,
			wantOK:     false,
		},
		{
			name:       "user ID has wrong type",
			input:      "1",
			wantUserID: 0,
			wantOK:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			ctx := context.NewContext()
			ctx.Reset(rec, req)

			if tt.input != nil {
				ctx.Input.SetData("userID", tt.input)
			}

			controller := &BaseController{}
			controller.Ctx = ctx

			userID, ok := controller.GetUserID()

			if ok != tt.wantOK {
				t.Fatalf("expected ok=%v, got %v", tt.wantOK, ok)
			}

			if userID != tt.wantUserID {
				t.Fatalf("expected user ID %d, got %d", tt.wantUserID, userID)
			}
		})
	}
}
