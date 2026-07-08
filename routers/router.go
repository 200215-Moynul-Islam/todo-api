package routers

import (
	"todo-api/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/health", &controllers.HealthController{})

	beego.Router("/tasks", &controllers.TaskController{}, "get:GetAll;post:Create")
	beego.Router("/tasks/:id", &controllers.TaskController{}, "get:GetByID;put:Update;delete:Delete")
}
