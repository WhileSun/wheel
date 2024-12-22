package pkg

import (
	"fmt"
	"testing"
	"time"

	"github.com/WhileSun/wheel/core/gconfig"
	"github.com/WhileSun/wheel/core/glog"
	"github.com/WhileSun/wheel/core/gserver"
	"github.com/WhileSun/wheel/database/gdb"
	"github.com/WhileSun/wheel/web/gjwt"
	"github.com/gin-gonic/gin"
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
	var glogConf glog.GlogConf
	gconfig.NewLoadFile(&glogConf, "./config.yaml")
	logger := glog.New(glogConf)

	var gserverConf gserver.GserverConf
	gconfig.NewLoadFile(&gserverConf, "./config.yaml")
	s := gserver.New(gserverConf)
	// s.SetHttpServer(&http.Server{
	// 	Addr: ":8080",
	// })
	s.SetHttpHandler(func() *gin.Engine {
		r := gin.New()
		r.Use(gserver.MiddlewareLogger(logger))
		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "test",
			})
		})
		return r
	}())
	s.Run()
}
