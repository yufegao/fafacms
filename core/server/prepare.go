package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexedwards/scs"
	"github.com/alexedwards/scs/stores/memstore"
	"github.com/alexedwards/scs/stores/redisstore"
	"github.com/hunterhug/fafacms/core/config"
	"github.com/hunterhug/fafacms/core/model"
	"github.com/hunterhug/fafacms/core/router"
	"github.com/hunterhug/fafacms/core/util/rdb"
	"github.com/hunterhug/fafacms/core/util/session"
	"io/ioutil"
	"time"
)

func InitConfig(configFilePath string) error {
	c := new(config.Config)
	if configFilePath == "" {
		return errors.New("config file empty")
	}

	raw, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(raw, c)
	if err != nil {
		return err
	}

	config.FafaConfig = c
	return nil
}

func InitRdb(dbConfig rdb.MyDbConfig) error {
	db, err := rdb.NewDb(dbConfig)
	if err != nil {
		return err
	}

	config.FafaRdb = db
	return nil
}

func InitSession(redisConf session.MyRedisConf) error {
	pool, err := session.NewRedis(&redisConf)
	if err != nil {
		return err
	}
	redisStore := redisstore.New(pool)
	config.FafaSessionMgr = scs.NewManager(redisStore)
	return nil
}

func InitMemorySession() {
	config.FafaSessionMgr = scs.NewManager(memstore.New(time.Hour * 1))
}

func CreateTable(tables []interface{}) {
	for _, table := range tables {
		ok, err := config.FafaRdb.IsTableExist(table)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if !ok {
			err = config.FafaRdb.CreateTables(table)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
		}

		err = config.FafaRdb.Client.CreateIndexes(table)
		if err != nil {
			fmt.Println(err.Error())
			//continue
		}
		err = config.FafaRdb.Client.CreateUniques(table)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	}

	InitResource()
}

func InitResource() {
	for url, handler := range router.V1Router {
		r := new(model.Resource)
		r.Url = fmt.Sprintf("/v1%s", url)
		r.Name = handler.Name
		r.Describe = handler.Name
		r.Admin = handler.Admin
		err := r.InsertOne()
		if err != nil {
			fmt.Println(err.Error())
		}

	}
}
