package middlewares

import (
	"github.com/vnaki/ris/types"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
)

func Cors(types.Engine) iris.Handler {
	return cors.New(cors.Options{
		AllowedMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowedOrigins:     []string{"*"},
		AllowedHeaders:     []string{"*"},
		AllowCredentials:   true,
		ExposedHeaders:     []string{"*"},
		OptionsPassthrough: false,
	})
}
