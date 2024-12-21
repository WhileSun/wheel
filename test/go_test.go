package pkg

import (
	"fmt"
	"testing"

	"github.com/WhileSun/wheel/core/gconfig"
	"github.com/WhileSun/wheel/core/glog"
	"github.com/WhileSun/wheel/database/gdb"
)

func TestGconfig(t *testing.T) {
	type Database struct {
		Host string `default:"110"`
	}
	var dbconfig Database
	gconfig.NewFlagFile(&dbconfig, "./config.yaml", "config file")
	fmt.Println(dbconfig)
	config := gconfig.NewViper("./config.yaml", "yaml")
	fmt.Println(config.GetString("host"))
}

func TestGlog(t *testing.T) {
	var glogConf glog.GlogConf
	gconfig.NewLoadFile(&glogConf, "./config.yaml")
	logger := glog.New(glogConf)
	logger.Info("test")
	logger.Error("test")
}

func TestGdb(t *testing.T) {
	type ProductList struct {
		ID   uint `gorm:"primarykey"`
		Name string
	}
	// var glogConf glog.GlogConf
	// gconfig.NewLoadFile(&glogConf, "./config.yaml")
	// glogConf.Path = "runtime/dbs"
	// glog := glog.New(glogConf)
	var gdbConf gdb.GdbConf
	gconfig.NewLoadFile(&gdbConf, "./config.yaml")
	// gdbConf.LogWriter = glog
	db := gdb.New(gdbConf)
	var user ProductList
	db.Where("name =?", "John").First(&user)
	fmt.Println(user)
}
