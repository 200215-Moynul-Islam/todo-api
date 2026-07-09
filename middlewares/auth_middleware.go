package middlewares

import (
	"net/http"
	"strings"
	"todo-api/utils"

	beegoCtx "github.com/beego/beego/v2/server/web/context"
)

func AuthFilter(ctx *beegoCtx.Context) {
	authHeader := ctx.Input.Header("Authorization")
	if authHeader == "" {
		utils.SendJSONResponse(ctx, http.StatusUnauthorized, false, "Authorization header is required", nil)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		utils.SendJSONResponse(ctx, http.StatusUnauthorized, false, "Authorization header must be Bearer {token}", nil)
		return
	}

	tokenString := parts[1]
	userID, err := utils.ValidateToken(tokenString)
	if err != nil {
		utils.SendJSONResponse(ctx, http.StatusUnauthorized, false, "Invalid or expired token", nil)
		return
	}

	// Store the user ID in the context so controllers can access it
	ctx.Input.SetData("userID", userID)
}
