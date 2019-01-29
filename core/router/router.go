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
		"/":       {controllers.Home, GP},
		"/login":  {controllers.Login, GP},
		"/logout": {controllers.Logout, GP},
	}

	// /v1/user/create
	V1Router = map[string]HttpHandle{
		"/user/create": {controllers.CreateUser, POST},
		"/user/update": {controllers.UpdateUser, POST},
		"/user/delete": {controllers.DeleteUser, POST},
		"/user/take":   {controllers.TakeUser, GP},
		"/user/list":   {controllers.ListUser, GP},

		"/group/create": {controllers.CreateGroup, POST},
		"/group/update": {controllers.UpdateGroup, POST},
		"/group/delete": {controllers.DeleteGroup, POST},
		"/group/take":   {controllers.TakeGroup, GP},
		"/group/list":   {controllers.ListGroup, GP},

		"/resource/create": {controllers.CreateResource, POST},
		"/resource/update": {controllers.UpdateResource, POST},
		"/resource/delete": {controllers.DeleteResource, POST},
		"/resource/take":   {controllers.TakeResource, GP},
		"/resource/list":   {controllers.ListResource, GP},

		"/auth/update": {controllers.UpdateAuth, GP},

		"/node/create": {controllers.CreateNode, POST},
		"/node/update": {controllers.UpdateNode, POST},
		"/node/delete": {controllers.DeleteNode, POST},
		"/node/take":   {controllers.TakeNode, GP},
		"/node/list":   {controllers.ListNode, GP},

		"/content/create": {controllers.CreateContent, POST},
		"/content/update": {controllers.UpdateContent, POST},
		"/content/delete": {controllers.DeleteContent, POST},
		"/content/take":   {controllers.TakeContent, GP},
		"/content/list":   {controllers.ListContent, GP},

		"/comment/create": {controllers.CreateComment, POST},
		"/comment/update": {controllers.UpdateComment, POST},
		"/comment/delete": {controllers.DeleteComment, POST},
		"/comment/take":   {controllers.TakeComment, GP},
		"/comment/list":   {controllers.ListComment, GP},
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
