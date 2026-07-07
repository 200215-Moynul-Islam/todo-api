package controllers

import (
	"github.com/beego/beego/v2/server/web"
)

type HealthController struct {
	web.Controller
}

func (c *HealthController) Get() {
	c.Data["json"] = map[string]any{
		"success": true,
		"message": "Service is healthy",
	}
	c.ServeJSON()
}