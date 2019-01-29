package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hunterhug/fafacms/core/controllers"
)

type HttpHandle struct {
	Func   gin.HandlerFunc
	Method []string
}

var (
	POST = []string{"POST"}
	GET  = []string{"GET"}
	GP   = []string{"POST", "GET"}
)

var (
	HomeRouter = map[string]HttpHandle{
		"/":      {controllers.Home, GP},
		"/login": {controllers.Login, GP},
	}

	// /v1/user/create
	V1Router = map[string]HttpHandle{
		"/user/create": {controllers.CreateUser, POST},
		"/user/update": {controllers.UpdateUser, POST},
		"/user/delete": {controllers.DeleteUser, POST},
		"/user/take":   {controllers.TakeUser, GP},
	}

	// /b/upload
	BaseRouter = map[string]HttpHandle{
		"/upload": {controllers.Upload, POST},
	}
)

// home end.
func SetRouter(router *gin.Engine) {
	for url, app := range HomeRouter {
		for _, method := range app.Method {
			router.Handle(method, url, app.Func)
		}
	}
}

func SetAPIRouter(router *gin.RouterGroup, handles map[string]HttpHandle) {
	for url, app := range handles {
		for _, method := range app.Method {
			router.Handle(method, url, app.Func)
		}
	}
}
