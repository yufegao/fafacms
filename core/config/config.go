package config

import (
	"encoding/json"
	"github.com/hunterhug/fafa/core/util/rdb"
	"github.com/hunterhug/fafa/core/util/session"
	"github.com/alexedwards/scs"
)

var (
	FafaConfig     *Config
	FafaRdb        *rdb.MyDb
	FafaSessionMgr *scs.Manager
)

type Config struct {
	DefaultConfig MyConfig
	DbConfig      rdb.MyDbConfig
	SessionConfig session.MyRedisConf
}

type MyConfig struct {
	WebPort     string
	LogPath     string
	StoragePath string
	Debug       bool
}

func JsonOutConfig(config Config) (string, error) {
	raw, err := json.Marshal(config)
	if err != nil {
		return "", err
	}

	back := string(raw)
	return back, nil
}
