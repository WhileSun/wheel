package gdb

import (
	"fmt"

	"gorm.io/driver/mysql"
)

func (gdbConf *GdbConf) mySqlInit() {
	dbURI = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		gdbConf.User,
		gdbConf.Password,
		gdbConf.Host,
		gdbConf.Port,
		gdbConf.Name,
		gdbConf.Charset)
	dialector = mysql.New(mysql.Config{
		DSN:                       dbURI, // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	})
}
