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
		//首页
		"/": {controllers.Home, GP},
	}

	V1Router = map[string]HttpHandle{
		//登陆相关
		"/login": {controllers.Login, GP},
	}
)

func SetRouter(router *gin.Engine) {
	for url, app := range HomeRouter {
		for _, method := range app.Method {
			router.Handle(method, url, app.Func)
		}
	}
}

func SetAPIRouter(router *gin.RouterGroup) {
	for url, app := range V1Router {
		for _, method := range app.Method {
			router.Handle(method, url, app.Func)
		}
	}
}
