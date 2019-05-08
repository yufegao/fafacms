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

// 路由，最后一个参数表示是否需要管理权限
var (
	HomeRouter = map[string]HttpHandle{
		// 前端路由
		"/":          {"Home", controllers.Home, GP, false},
		"/p":         {"List Peoples", controllers.Peoples, GP, false},         // 列出用户
		"/u/node":    {"List User Nodes", controllers.Nodes, GP, false},        // 列出某用户下的节点
		"/u/info":    {"List User Info", controllers.UserInfo, GP, false},      // 获取某用户信息
		"/u/count":   {"Count User Content", controllers.UserCount, GP, false}, // 统计某用户文章情况（某用户可留空）
		"/u/content": {"List User Content", controllers.Contents, GP, false},   // 列出某用户下文章
		"/c":         {"Get Content", controllers.Content, GP, false},          // 获取文章

		// 前端的用户授权路由，不需要登录即可操作
		"/login":           {"User Login", controllers.Login, GP, false},
		"/logout":          {"User Logout", controllers.Logout, GP, false},
		"/register":        {"User Register", controllers.RegisterUser, GP, false},
		"/activate":        {"User Verify Email To Activate", controllers.ActivateUser, GP, false},               // 用户自己激活
		"/activate/code":   {"User Resend Email Activate Code", controllers.ResendActivateCodeToUser, GP, false}, // 激活码过期重新获取
		"/password/forget": {"User Forget Password Gen Code", controllers.ForgetPasswordOfUser, GP, false},       // 忘记密码，验证码发往邮箱
		"/password/change": {"User Change Password", controllers.ChangePasswordOfUser, GP, false},                // 根据邮箱验证码修改密码
	}

	// /v1/user/create
	// need login group auth
	V1Router = map[string]HttpHandle{
		// 用户组操作
		"/group/create":        {"Create Group", controllers.CreateGroup, POST, true},
		"/group/update":        {"Update Group", controllers.UpdateGroup, POST, true},
		"/group/delete":        {"Delete Group", controllers.DeleteGroup, POST, true},
		"/group/take":          {"Take Group", controllers.TakeGroup, GP, true},
		"/group/list":          {"List Group", controllers.ListGroup, GP, true},
		"/group/user/list":     {"Group List User", controllers.ListGroupUser, GP, true},         // 超级管理员列出组下的用户
		"/group/resource/list": {"Group List Resource", controllers.ListGroupResource, GP, true}, // 超级管理员列出组下的资源

		// 用户操作
		"/user/list":         {"User List All", controllers.ListUser, GP, true},              // 超级管理员列出用户列表
		"/user/create":       {"User Create", controllers.CreateUser, GP, true},              // 超级管理员创建用户，默认激活
		"/user/assign":       {"User Assign Group", controllers.AssignGroupToUser, GP, true}, // 超级管理员给用户分配用户组
		"/user/info":         {"User Info Self", controllers.TakeUser, GP, false},            // 获取自己的信息
		"/user/update":       {"User Update Self", controllers.UpdateUser, GP, false},        // 更新自己的信息
		"/user/admin/update": {"User Update Admin", controllers.UpdateUserAdmin, GP, true},   // 管理员修改其他用户信息

		// 资源操作
		"/resource/list":   {"Resource List All", controllers.ListResource, GP, true},               // 列出资源
		"/resource/assign": {"Resource Assign Group", controllers.AssignGroupAndResource, GP, true}, // 资源分配给组

		// 文件操作
		"/file/upload":       {"File Upload", controllers.UploadFile, POST, false},
		"/file/list":         {"File List Self", controllers.ListFile, POST, false},
		"/file/update":       {"File Update Self", controllers.UpdateFile, POST, false},
		"/file/admin/list":   {"File List All", controllers.ListFileAdmin, POST, true},     // 管理员查看所有文件
		"/file/admin/update": {"File Update All", controllers.UpdateFileAdmin, POST, true}, // 管理员修改文件

		// 内容节点操作
		"/node/create":     {"Create Node Self", controllers.CreateNode, POST, false},
		"/node/update":     {"Update Node Self", controllers.UpdateNode, POST, false},
		"/node/delete":     {"Delete Node Self", controllers.DeleteNode, POST, false},
		"/node/take":       {"Take Node Self", controllers.TakeNode, GP, false},
		"/node/list":       {"List Node Self", controllers.ListNode, GP, false},
		"/node/admin/list": {"List Node All", controllers.ListNodeAdmin, GP, true}, // 管理员查看其他用户节点

		// 内容操作
		"/content/create":  {"Create Content Self", controllers.CreateContent, POST, false},   // 创建文章内容
		"/content/update":  {"Update Content Self", controllers.UpdateContent, POST, false},   // 更新内容，更新会写入预览
		"/content/publish": {"Publish Content Self", controllers.PublishContent, POST, false}, // 将预览刷进另外一个字段
		"/content/cancel":  {"Cancel Content Self", controllers.CancelContent, POST, false},   // 取消预览的内容，刷回来

		"/content/list":               {"List Content Self", controllers.ListContent, GP, false},                   // 列出文章
		"/content/admin/list":         {"List Content All", controllers.ListContentAdmin, GP, true},                // 管理员列出文章，什么类型都可以
		"/content/history/list":       {"List Content History Self", controllers.ListContentHistory, GP, false},    // 列出文章的历史记录
		"/content/history/admin/list": {"List Content History All", controllers.ListContentHistoryAdmin, GP, true}, // 管理员列出文章的历史纪录

		"/content/take":               {"Take Content Self", controllers.TakeContent, GP, false},                    // 获取文章内容
		"/content/admin/take":         {"Take Content Self", controllers.TakeContentAdmin, GP, true},                // 管理员获取文章内容
		"/content/history/take":       {"Take Content History Self", controllers.TakeContentHistory, GP, false},     // 获取文章历史内容
		"/content/history/admin/take": {"Take Content History Self", controllers.TakeContentHistoryAdmin, GP, true}, // 管理员获取文章历史内容

		"/content/admin/update": {"Update Content All", controllers.DeleteContentAdmin, POST, true},                  // 超级管理员修改文章，比如禁用或者逻辑删除/恢复文章
		"/content/rubbish":      {"Sent Content Self To Rubbish", controllers.DeleteContent, POST, false},            // 一般回收站
		"/content/redo":         {"Sent Rubbish Content Self To Origin", controllers.DeleteContentRedo, POST, false}, // 一般回收站恢复
		"/content/delete":       {"Delete Content Self Logic", controllers.DeleteContent2, POST, false},              // 逻辑删除文章

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
