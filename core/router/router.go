package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hunterhug/fafa/core/controllers"
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

var Router = map[string]HttpHandle{
	//登陆相关
	"/login": {controllers.Login, GP},
}

func SetRouter(router *gin.Engine) {
	for url, app := range Router {
		for _, method := range app.Method {
			router.Handle(method, url, app.Func)
		}
	}
}
