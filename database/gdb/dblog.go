package gdb

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

// Writer 定义自己的Writer
type Writer struct {
	log *logrus.Logger
}

// Printf 实现gorm/logger.Writer接口
func (w *Writer) Printf(format string, v ...interface{}) {
	logStr := fmt.Sprintf(format, v...)
	w.log.Info(logStr + "\r\n")
}
