package main

import (
	"github.com/vnaki/ris"
	"github.com/vnaki/ris/examples/routes"
	"github.com/vnaki/ris/middlewares"
	"github.com/vnaki/ris/plugins"
)

func main()  {
	e := ris.New()

	// post max memory
	e.SetPostMemory(20 << 20)

	e.RouteMiddleware(middlewares.Cors)

	e.Plugin("logger", plugins.LoggerPlugin)
	//e.Plugin("data", plugins.MysqlPlugin)
	e.Plugin("data", plugins.SqlitePlugin)

	// default module
	e.Module("/", routes.ApiRoute)

	if err := e.Run("./config/app.yaml"); err != nil {
		panic(err)
	}
}
