package gconfig

import (
	"flag"
	"fmt"

	"github.com/jinzhu/configor"
	"github.com/spf13/viper"
)

// load flag file
func NewFlagFile(config interface{}, filepath string, fileUsage string) {
	configFile := flag.String("f", filepath, fileUsage)
	flag.Parse()
	configor.Load(config, *configFile)
}

// load files
func NewLoadFile(config interface{}, files ...string) {
	configor.Load(config, files...)
}

func NewViper(filepath string, filetype string) *viper.Viper {
	config := viper.New()
	config.SetConfigType(filetype)
	config.SetConfigFile(filepath)
	if err := config.ReadInConfig(); err != nil {
		fmt.Printf("读取配置文件失败: %v\n", err)
	}
	return config
}
