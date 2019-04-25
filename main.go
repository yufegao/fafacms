/*
    2019-4-24：

	程序主入口
	花花CMS是一个内容管理系统，代码尽可能地补充必要注释，方便后人协作
**/
package main

import (
	"flag"
	"github.com/hunterhug/fafacms/core/config"
	"github.com/hunterhug/fafacms/core/controllers"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"github.com/hunterhug/fafacms/core/router"
	"github.com/hunterhug/fafacms/core/server"
	"github.com/hunterhug/fafacms/core/util/mail"
)

var (
	// 全局配置文件路径
	configFile string

	// 是否创建数据库表
	createTable bool

	// 开发时每次都发邮件的形式不好，可以先调试模式
	mailDebug bool

	// 跳过授权，某些超级管理接口需要绑定组和路由，可以先开调试模式
	canSkipAuth bool

	// 分布式Session开关，可以先开调试模式，存于内存中
	sessionUseRedis bool
)

// 初始化时解析命令行，辅助程序
// 这些调试参数不置于文件配置中
func init() {
	// 默认读取本路径下 ./config.json 配置
	flag.StringVar(&configFile, "config", "./config.json", "config file")

	// 正式部署时，请全部设置为 false
	flag.BoolVar(&createTable, "init_db", true, "create db table")
	flag.BoolVar(&mailDebug, "email_debug", true, "Email debug")
	flag.BoolVar(&canSkipAuth, "auth_skip_debug", true, "Auth skip debug")
	flag.BoolVar(&sessionUseRedis, "use_session_redis", false, "Use Redis Session")
	flag.Parse()
}

// 入口
// 欢迎查看优美代码，我是花花
func main() {

	// 将调试参数跨包注入
	mail.Debug = mailDebug
	controllers.AuthDebug = canSkipAuth

	var err error

	// 初始化全局配置
	err = server.InitConfig(configFile)
	if err != nil {
		panic(err)
	}

	// 初始化日志
	flog.InitLog(config.FafaConfig.DefaultConfig.LogPath)

	// 如果全局调试，那么所有DEBUG以上级别日志将会打印
	// 实际情况下，最好设置为 true，
	if config.FafaConfig.DefaultConfig.LogDebug {
		flog.SetLogLevel("DEBUG")
	}

	flog.Log.Notice("Hi! FaFa CMS!")
	flog.Log.Debugf("Hi! Config is %#v", config.FafaConfig)

	// 初始化数据库连接
	err = server.InitRdb(config.FafaConfig.DbConfig)
	if err != nil {
		panic(err)
	}

	// 初始化网站Session存储
	if sessionUseRedis {
		err = server.InitSession(config.FafaConfig.SessionConfig)
		if err != nil {
			panic(err)
		}
	} else {
		server.InitMemorySession()
	}

	// 创建数据库表，需要先手动创建DB
	if createTable {
		server.CreateTable([]interface{}{
			model.User{},
			model.Group{},
			model.Resource{},
			model.GroupResource{},
			model.Content{},
			model.ContentNode{},
			model.ContentSupport{},
			model.ContentCal{},
			model.Comment{},
			model.CommentSupport{},
			model.CommentCal{},
			model.Log{},
			model.File{},
		})
	}

	// Server Run
	engine := server.Server()

	// Storage static API
	engine.Static("/storage", config.FafaConfig.DefaultConfig.StoragePath)
	engine.Static("/storage_x", config.FafaConfig.DefaultConfig.StoragePath+"_x")

	// Web welcome home!
	router.SetRouter(engine)

	// V1 API, will may be change to V2...
	v1 := engine.Group("/v1")
	v1.Use(controllers.AuthFilter)

	// Router Set
	router.SetAPIRouter(v1, router.V1Router)

	flog.Log.Noticef("Server run in %s", config.FafaConfig.DefaultConfig.WebPort)
	err = engine.Run(config.FafaConfig.DefaultConfig.WebPort)
	if err != nil {
		panic(err)
	}
}
