package main

import (
	"todo-api/database"
	_ "todo-api/routers"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	// 1. Direct, high-priority manual header interceptor
	beego.InsertFilter("*", beego.BeforeStatic, func(ctx *context.Context) {
		// Allow any origin to communicate with this app
		ctx.Output.Header("Access-Control-Allow-Origin", "*")
		
		// Permit the methods used by the frontend dashboard engine
		ctx.Output.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		
		// Crucial: Explicitly approve Authorization and Content-Type headers
		ctx.Output.Header("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type, Accept, X-Requested-With")
		
		// 2. Intercept and terminate the browser preflight handshake immediately
		if ctx.Input.Method() == "OPTIONS" {
			ctx.Output.SetStatus(200)
			// Send an empty string payload to immediately close out the browser connection successfully
			_ = ctx.Output.Body([]byte("")) 
			return
		}
	})

	database.Init()
	
	beego.Run()
}
