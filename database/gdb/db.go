package gdb

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// GdbConf 数据库连接设置
type GdbConf struct {
	Type        string         `default:"postgres"`
	Host        string         `default:"127.0.0.1"`
	Port        string         `default:"5432"`
	User        string         `default:"postgres"`
	Password    string         `default:"123456"`
	Name        string         `default:"test"`
	TablePrefix string         `default:""`
	Charset     string         `default:"utf8"`
	MaxIdleConn int            `default:"20"`
	MaxOpenConn int            `default:"200"`
	Log         bool           `default:"false"`
	LogWriter   *logrus.Logger //日志插件
}

var (
	dbURI     string
	dialector gorm.Dialector
)

// New 初始化数据库连接
func New(gdbConf GdbConf) *gorm.DB {
	return gdbConf.run()
}

// 设置日志
func (gdbConf *GdbConf) SetLogger(logger *logrus.Logger) {
	gdbConf.LogWriter = logger
}

// run 初始化运行
func (gdbConf *GdbConf) run() *gorm.DB {
	dbType := strings.ToLower(gdbConf.Type)
	if dbType == "mysql" {
		gdbConf.mySqlInit()
	} else if dbType == "postgres" {
		gdbConf.postGreSqlInit()
	} else {
		fmt.Printf("gdb config type [%s] is not setting, use [mysql,postgres] \n", dbType)
		return nil
	}
	var newLogger logger.Interface
	// 开启数据库日式记录
	if gdbConf.Log {
		var writer logger.Writer

		if gdbConf.LogWriter != nil {
			writer = &Writer{log: gdbConf.LogWriter}
		} else {
			writer = log.New(os.Stdout, "\r\n", log.LstdFlags) // gorm io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		}
		newLogger = logger.New(
			writer,
			logger.Config{
				SlowThreshold:             time.Second, // 慢 SQL 阈值
				LogLevel:                  logger.Info, // 日志级别
				IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  false,       // 禁用彩色打印
			},
		)
	}
	conn, err := gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{TablePrefix: gdbConf.TablePrefix, SingularTable: true},
		Logger:         newLogger,
	})
	if err != nil {
		fmt.Println("gdb gorm conn failed")
		return nil
	}
	sqlDB, err := conn.DB()
	if err != nil {
		fmt.Printf("gdb connect server failed,Error: %s \n", err.Error())
		return nil
	}
	sqlDB.SetMaxIdleConns(gdbConf.MaxIdleConn)  // 空闲进程数
	sqlDB.SetMaxOpenConns(gdbConf.MaxOpenConn)  // 最大进程数
	sqlDB.SetConnMaxLifetime(time.Second * 600) // 设置了连接可复用的最大时间
	if err := sqlDB.Ping(); err != nil {
		fmt.Printf("gdb ping is failed,Error:%s", err.Error())
	}
	return conn
}
