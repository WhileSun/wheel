package pkg

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/WhileSun/wheel/core/gconfig"
	"github.com/WhileSun/wheel/core/glog"
	"github.com/WhileSun/wheel/core/gserver"
	"github.com/WhileSun/wheel/database/gdb"
	"github.com/WhileSun/wheel/web/gjwt"
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

func TestGjwt(t *testing.T) {
	var gjwtConf gjwt.GjwtConf
	gconfig.NewLoadFile(&gjwtConf, "./config.yaml")
	jwt := gjwt.New(gjwtConf)
	// create token
	token, _ := jwt.CreateToken(map[string]interface{}{
		"name": "John",
	})
	// parse token
	res, _ := jwt.ParseToken(token)
	fmt.Printf("jwt object %+v \n", res)

	time.Sleep(time.Second * 3)
	// get new token
	newToken, err := jwt.RefreshToken(token)
	fmt.Println("newToken", newToken, err)
}

func TestGserver(t *testing.T) {
	var gserverConf gserver.GserverConf
	gconfig.NewLoadFile(&gserverConf, "./config.yaml")
	gserver.New(gserverConf).Run(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}))
}
