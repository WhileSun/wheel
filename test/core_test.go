package pkg

import (
	"fmt"
	"testing"

	"github.com/WhileSun/wheel/core/gconfig"
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
