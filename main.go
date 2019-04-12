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
	configFile  string
	createTable bool
	mailDebug   bool
)

func init() {
	flag.StringVar(&configFile, "config", "./config.json", "config file")
	flag.BoolVar(&createTable, "db", true, "create db table")
	flag.BoolVar(&mailDebug, "eb", true, "Email debug")
	flag.Parse()
}

func main() {
	mail.Debug = mailDebug

	var err error

	// Init Config
	err = server.InitConfig(configFile)
	if err != nil {
		panic(err)
	}

	// Init Log
	flog.InitLog(config.FafaConfig.DefaultConfig.LogPath)
	if config.FafaConfig.DefaultConfig.Debug {
		flog.SetLogLevel("DEBUG")
	}

	flog.Log.Notice("Hi! FaFa CMS!")
	flog.Log.Debugf("Hi! Config is %#v", config.FafaConfig)

	// Init Db
	err = server.InitRdb(config.FafaConfig.DbConfig)
	if err != nil {
		panic(err)
	}

	err = server.InitSession(config.FafaConfig.SessionConfig)
	if err != nil {
		panic(err)
	}

	// Create Table, Here to init db
	if createTable {
		server.CreateTable([]interface{}{
			model.User{},
			model.Group{},
			model.Resource{},
			model.GroupResource{},
			model.Content{},
			model.ContentNode{},
			model.Comment{},
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
