package routes

import (
	"github.com/kataras/iris/v12/mvc"
	"github.com/vnaki/ris/examples/controllers"
)

func ApiRoute(m *mvc.Application) {
	m.Party("/api").Handle(controllers.NewApiController())
}
