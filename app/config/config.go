package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func InitViper() (err error) {
	viper.SetConfigType("yaml")     // 指定配置文件路径
	viper.SetConfigName("config")   // 配置文件名称(无扩展名)
	viper.AddConfigPath("./config") // 指定查找目录
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件被修改...")
	})
	return
}
