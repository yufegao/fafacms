package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hunterhug/fafacms/core/controllers"
)

type HttpHandle struct {
	Name   string
	Func   gin.HandlerFunc
	Method []string
	Admin  bool
}

var (
	POST = []string{"POST"}
	GET  = []string{"GET"}
	GP   = []string{"POST", "GET"}
)

var (
	HomeRouter = map[string]HttpHandle{
		"/":                       {"Home", controllers.Home, GP, false},
		"/u/:name":                {"user home page", controllers.Home, GP, false},
		"/u/:name/:node":          {"user node page", controllers.Home, GP, false},
		"/u/:name/:node/:content": {"user content page", controllers.Home, GP, false},
		"/login":                  {"User Login", controllers.Login, GP, false},
		"/logout":                 {"User Logout", controllers.Logout, GP, false},
		"/register":               {"User Register", controllers.RegisterUser, GP, false},
		"/activate":               {"User Verify Email To Activate", controllers.ActivateUser, GP, false},
		"/activate/code":          {"User Resend Email Activate Code", controllers.ResendActivateCodeToUser, GP, false},
		"/password/forget":        {"User Forget Password Gen Code", controllers.ForgetPasswordOfUser, GP, false},
		"/password/change":        {"User Change Password", controllers.ChangePasswordOfUser, GP, false},
	}

	// /v1/user/create
	// need login group auth
	V1Router = map[string]HttpHandle{

		"/group/create": {"Create Group", controllers.CreateGroup, POST, true},
		"/group/update": {"Update Group", controllers.UpdateGroup, POST, true},
		"/group/delete": {"Delete Group", controllers.DeleteGroup, POST, true},
		"/group/take":   {"Take Group", controllers.TakeGroup, GP, true},
		"/group/list":   {"List Group", controllers.ListGroup, GP, true},

		"/user/list":   {"User List All", controllers.ListUser, GP, true},
		"/user/assign": {"User Assign Group", controllers.AssignGroupToUser, GP, true},
		"/user/info":   {"User Info Self", controllers.TakeUser, GP, false},
		"/user/update": {"User Update Self", controllers.UpdateUser, GP, false},

		"/resource/list":   {"Resource List All", controllers.ListResource, GP, true},
		"/resource/assign": {"Resource Assign Group", controllers.AssignGroupAndResource, GP, true},


		"/file/upload":       {"File Upload", controllers.UploadFile, POST, false},
		"/file/list":         {"File List Self", controllers.ListFile, POST, false},
		"/file/update":       {"File Update Self", controllers.UpdateFile, POST, false},
		"/file/admin/list":   {"File List All", controllers.ListFileAdmin, POST, true},
		"/file/admin/update": {"File Update All", controllers.UpdateFileAdmin, POST, true},

		"/node/create": {"Create Node Self", controllers.CreateNode, POST, false},
		"/node/update": {"Update Node Self", controllers.UpdateNode, POST, false},
		"/node/delete": {"Delete Node Self", controllers.DeleteNode, POST, false},
		"/node/take":   {"Take Node Self", controllers.TakeNode, GP, false},
		"/node/list":   {"List Node Self", controllers.ListNode, GP, false},
		//
		//"/content/create": {controllers.CreateContent, POST},
		//"/content/update": {controllers.UpdateContent, POST},
		//"/content/delete": {controllers.DeleteContent, POST},
		//"/content/take":   {controllers.TakeContent, GP},
		//"/content/list":   {controllers.ListContent, GP},
		//
		//"/comment/create": {controllers.CreateComment, POST},
		//"/comment/update": {controllers.UpdateComment, POST},
		//"/comment/delete": {controllers.DeleteComment, POST},
		//"/comment/take":   {controllers.TakeComment, GP},
		//"/comment/list":   {controllers.ListComment, GP},
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
