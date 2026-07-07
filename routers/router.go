package routers

import (
	"todo-api/controllers"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.Router("/health", &controllers.HealthController{})
}
