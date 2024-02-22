package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type config struct {
	Mysql Mysql
	Jwt   Jwt
}

var CONFIG config

func LoadConfig() {
	v := viper.New()
	v.SetConfigFile("./config/config.yaml")
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		fmt.Println("load config error, ", err)
		return
	}

	// 监听配置文件
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		v.Unmarshal(&CONFIG)
	})

	v.Unmarshal(&CONFIG)
}
