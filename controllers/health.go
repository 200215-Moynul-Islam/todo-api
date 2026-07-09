package controllers

import (
	"net/http"
	"todo-api/utils"
)

type HealthController struct {
	BaseController
}

func (c *HealthController) Get() {
	utils.SendJSONResponse(c.Ctx, http.StatusOK, true, "Server is running", nil)
}