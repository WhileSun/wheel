package gresp

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type GrespConf struct {
	ctx *gin.Context
}

type RespData struct {
	Code uint        `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
	Time int64       `json:"time"`
}

func New(req *gin.Context) *GrespConf {
	return &GrespConf{
		ctx: req,
	}
}

func (grespConf GrespConf) JSON(code uint, msg string, data interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}
	grespConf.ctx.JSON(http.StatusOK, &RespData{
		Code: code,
		Msg:  msg,
		Data: data,
		Time: time.Now().Unix(),
	})
}
