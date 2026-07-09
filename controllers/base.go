package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type BaseController struct {
	beego.Controller
}

// GetUserID retrieves the authenticated user ID from the request context
func (c *BaseController) GetUserID() (int, bool) {
	val := c.Ctx.Input.GetData("userID")
	userID, ok := val.(int)
	return userID, ok
}