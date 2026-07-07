package controllers

import "net/http"

type HealthController struct {
	BaseController
}

func (c *HealthController) Get() {
	c.SendSuccess(http.StatusOK, "Server is running", nil)
}