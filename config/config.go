package config

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"os"
)

var Conf config

type mysqlDB struct {
	DataSource   string `toml:"data_source"`
	MaxIdleConns int    `toml:"max_idle_conns"`
	MaxOpenConns int    `toml:"max_open_conns"`
}

type http struct {
	Listen string `toml:"listen"`
}

type config struct {
	Http    http    `toml:"http"`
	MysqlDB mysqlDB `toml:"mysql_db"`
}

func ParseConfig(configFile string) {
	if fileInfo, err := os.Stat(configFile); err != nil {
		if os.IsNotExist(err) {
			log.Panicf("configuration file %v does not exist.", configFile)
		} else {
			log.Panicf("configuration file %v can not be stated. %v", configFile, err)
		}
	} else {
		if fileInfo.IsDir() {
			log.Panicf("%v is a directory name", configFile)
		}
	}

	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Panicf("read configuration file error. %v", err)
	}
	content = bytes.TrimSpace(content)

	err = toml.Unmarshal(content, &Conf)
	if err != nil {
		log.Panicf("unmarshal toml object error. %v", err)
	}
}
