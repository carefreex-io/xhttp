package xhttp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func SuccessResponse(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "OK",
		Data: data,
	})
}

func ErrorResponse(ctx *gin.Context, code int, err error) {
	ctx.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  err.Error(),
	})
}
