package routers

import (
	"todo-api/controllers"
	"todo-api/middlewares"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/health", &controllers.HealthController{})

	beego.Router("/auth/register", &controllers.AuthController{}, "post:Register")
	beego.Router("/auth/login", &controllers.AuthController{}, "post:Login")

	// Apply JWT authentication filter to tasks endpoints
	beego.InsertFilter("/tasks", beego.BeforeRouter, middlewares.AuthFilter)
	beego.InsertFilter("/tasks/*", beego.BeforeRouter, middlewares.AuthFilter)

	beego.Router("/tasks", &controllers.TaskController{}, "get:GetAll;post:Create")
	beego.Router("/tasks/generate-description", &controllers.TaskController{}, "post:GenerateDescription")
	beego.Router("/tasks/:id", &controllers.TaskController{}, "get:GetByID;put:Update;delete:Delete")
}
