package main

import (
	"flag"
	"github.com/hunterhug/fafacms/core/config"
	"github.com/hunterhug/fafacms/core/controllers"
	"github.com/hunterhug/fafacms/core/model"
	"github.com/hunterhug/fafacms/core/router"
	"github.com/hunterhug/fafacms/core/server"
)

var (
	configFile  string
	createTable bool
)

func init() {
	flag.StringVar(&configFile, "config", "./config.json", "config file")
	flag.BoolVar(&createTable, "db", true, "create db table")
	flag.Parse()
}

func main() {
	var err error

	// Init Config
	err = server.InitConfig(configFile)
	if err != nil {
		panic(err)
	}

	// Init Log
	config.InitLog(config.FafaConfig.DefaultConfig.LogPath)
	if config.FafaConfig.DefaultConfig.Debug {
		config.SetLogLevel("DEBUG")
	}

	config.Log.Notice("Hi! FaFa CMS!")
	config.Log.Debugf("Hi! Config is %#v", config.FafaConfig)

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
		})
	}

	// Server Run
	engine := server.Server()

	// Storage static API
	engine.Static("/storage", config.FafaConfig.DefaultConfig.StoragePath)

	// Web welcome home!
	router.SetRouter(engine)

	// V1 API
	v1 := engine.Group("/v1")
	v1.Use(controllers.AuthManager)
	router.SetAPIRouter(v1)

	config.Log.Noticef("Run in %s", config.FafaConfig.DefaultConfig.WebPort)
	err = engine.Run(config.FafaConfig.DefaultConfig.WebPort)
	if err != nil {
		panic(err)
	}
}
