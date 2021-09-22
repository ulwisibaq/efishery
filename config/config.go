package config

import (
	"log"

	gcfg "gopkg.in/gcfg.v1"
)

type MainConfig struct {
	Database struct {
		MysqlDSN string
	}

	CommoditiesAPI struct {
		Url string
	}
}

func ReadConfig(cfg interface{}) interface{} {
	ok := ReadModuleConfig(cfg, "files")
	if !ok {
		log.Fatalln("failed to read config")
	}
	return cfg
}

func ReadModuleConfig(cfg interface{}, path string) bool {
	filename := path + "/main.development.ini"
	err := gcfg.ReadFileInto(cfg, filename)
	if err == nil {
		return true
	}

	return false
}
