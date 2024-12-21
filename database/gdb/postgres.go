package gdb

import (
	"fmt"

	"gorm.io/driver/postgres"
)

func (gdbConf *GdbConf) postGreSqlInit() {
	dbURI = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		gdbConf.Host,
		gdbConf.Port,
		gdbConf.User,
		gdbConf.Name,
		gdbConf.Password)
	dialector = postgres.New(postgres.Config{
		DSN:                  dbURI,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	})
}
