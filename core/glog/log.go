package glog

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type GlogConf struct {
	Type         string `default:"file"`
	Path         string `default:"runtime/logs"`
	FileName     string `default:"sys"`
	MaxAge       int    `default:"168"`
	RotationTime int    `default:"24"`
	Stdout       bool   `default:"false"`
	LogLevel     string `default:"info"`
}

// New 初始化日志
func New(glogConf GlogConf) *logrus.Logger {
	return glogConf.run()
}

// new logrus
func (glogConf *GlogConf) run() *logrus.Logger {
	//logrus初始化
	logger := logrus.New()
	// 设置日志级别
	level, err := logrus.ParseLevel(glogConf.LogLevel)
	if err != nil {
		log.Fatalf("glog config level [%s] is not support, choose types [panic,fatal,error,warn,info,debug,trace]", glogConf.LogLevel)
		return nil
	}
	logger.SetLevel(level)
	logger.SetReportCaller(true)
	//logger标准化日志
	if glogConf.Stdout {
		logger.SetFormatter(new(LogFormatter))
		logger.SetOutput(os.Stdout)
		return logger
	}
	//判断日志类型
	logType := strings.ToLower(glogConf.Type)
	if logType == "file" {
		src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			log.Fatalf("glog Open Src File err %+v", err)
		}
		writer := bufio.NewWriter(src)
		logger.SetOutput(writer)
		glogConf.configLocalFileLogger(logger)
	} else {
		log.Fatalf("glog config type [%s] is not support, choose types [file]", logType)
	}
	return logger
}

/*
* ConfigLocalFileLogger 写入文件
logPath logs文件目录
logFileName 文件名
maxAge 文件最大保存时间
rotationTime 日志切割时间
*/
func (glogConf *GlogConf) configLocalFileLogger(log *logrus.Logger) {
	logPath := glogConf.Path
	logFileName := glogConf.FileName
	maxAge := glogConf.MaxAge
	rotationTime := glogConf.RotationTime
	// 文件目录
	baseLogPath := path.Join(logPath, logFileName)

	createLogWriter := func(logType string) *rotatelogs.RotateLogs {
		logFilePath := fmt.Sprintf("%s_%s_log", baseLogPath, logType)
		writer, err := rotatelogs.New(
			logFilePath+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(baseLogPath),                               // 生成软链，指向最新日志文件
			rotatelogs.WithMaxAge(time.Duration(maxAge)*time.Hour),             // 文件最大保存时间
			rotatelogs.WithRotationTime(time.Duration(rotationTime)*time.Hour), // 日志切割时间间隔
		)
		if err != nil {
			log.Fatalf("glog ailed to create %s log writer: %s", logType, err.Error())
		}
		return writer
	}
	writerTrace := createLogWriter("trace")
	writer := createLogWriter("access")
	writerError := createLogWriter("error")

	writeMap := lfshook.WriterMap{
		logrus.TraceLevel: writerTrace,
		logrus.InfoLevel:  writer,
		logrus.DebugLevel: writer,
		logrus.WarnLevel:  writer,
		logrus.FatalLevel: writerError,
		logrus.ErrorLevel: writerError,
		logrus.PanicLevel: writerError,
	}
	Hook := lfshook.NewHook(writeMap, new(LogFormatter))
	log.AddHook(Hook)
}

//func ConfigESLogger(esUrl string, esHOst string, index stringm) {
//	client, err := elastic.NewClient(elastic.SetURL(esUrl))
//	if err != nil {
//		log.Errorf("config es logger error. %+v", errors.WithStack(err))
//	}
//	esHook, err := elogrus.NewElasticHook(client, esHOst, log.DebugLevel, index)
//	if err != nil {
//		log.Errorf("config es logger error. %+v", errors.WithStack(err))
//	}
//	log.AddHook(esHook)
//}
