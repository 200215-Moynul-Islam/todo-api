package utils

import (
	beegoCtx "github.com/beego/beego/v2/server/web/context"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    any         `json:"data,omitempty"`
}

// SendJSONResponse writes a standardized APIResponse directly to Beego context output.
func SendJSONResponse(ctx *beegoCtx.Context, status int, success bool, message string, data any) {
	ctx.Output.SetStatus(status)
	response := APIResponse{
		Success: success,
		Message: message,
		Data:    data,
	}
	ctx.Output.JSON(response, false, false)
}
