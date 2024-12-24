package uresp

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UrespConf struct {
	ctx *gin.Context
}

type RespData struct {
	Code uint        `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
	Time int64       `json:"time"`
}

func New(req *gin.Context) *UrespConf {
	return &UrespConf{
		ctx: req,
	}
}

func (urespConf UrespConf) JSON(code uint, msg string, data interface{}) {
	urespConf.ctx.JSON(http.StatusOK, &RespData{
		Code: code,
		Msg:  msg,
		Data: data,
		Time: time.Now().Unix(),
	})
}
