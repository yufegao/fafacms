package main

import (
	"flag"
	"github.com/hunterhug/fafacms/core/config"
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

	config.Log.Notice("Hi! Fafa blog!")
	config.Log.Debugf("Hi! %#v", config.FafaConfig)

	// Init Db
	err = server.InitRdb(config.FafaConfig.DbConfig)
	if err != nil {
		panic(err)
	}

	err = server.InitSession(config.FafaConfig.SessionConfig)
	if err != nil {
		panic(err)
	}

	// Create Table
	if createTable {
		server.CreateTable([]interface{}{
			model.User{},
		})
	}

	// Server Run
	engine := server.Server()
	router.SetRouter(engine)

	config.Log.Noticef("Run in %s", config.FafaConfig.DefaultConfig.WebPort)
	err = engine.Run(config.FafaConfig.DefaultConfig.WebPort)
	if err != nil {
		panic(err)
	}
}
