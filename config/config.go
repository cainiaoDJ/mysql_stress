package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

type AppConfig struct {
	Host     string `yaml:"host"`
	Port     uint   `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	InitRows uint   `yaml:"init_rows"`
}

var Cfg *AppConfig

func LoadAppConfig(cfgPath string) {
	data, err := ioutil.ReadFile(filepath.Join(cfgPath, "/", "cfg.yml"))
	if err != nil {
		fmt.Println("read config err, load default value")
		Cfg = &AppConfig{
			Host:     "127.0.0.1",
			Port:     3306,
			Username: "root",
			Password: "123456",
			InitRows: 10000,
		}
		return
	}
	err = yaml.Unmarshal(data, &Cfg)
	if err != nil {
		panic("parse config err:" + err.Error())
	}
}