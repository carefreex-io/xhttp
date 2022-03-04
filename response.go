package xhttp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func SuccessResponse(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Msg:  "OK",
		Data: data,
	})
}

func BadRequestResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, Response{
		Msg: err.Error(),
	})
}

func UnauthorizedResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusUnauthorized, Response{
		Msg: err.Error(),
	})
}

func ForbiddenResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusForbidden, Response{
		Msg: err.Error(),
	})
}

func NotFoundResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusNotFound, Response{
		Msg: err.Error(),
	})
}

func InternalServerErrorResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, Response{
		Msg: err.Error(),
	})
}

func NotImplementedResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusNotImplemented, Response{
		Msg: err.Error(),
	})
}

func ServiceUnavailableResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusServiceUnavailable, Response{
		Msg: err.Error(),
	})
}

func GatewayTimeoutResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusGatewayTimeout, Response{
		Msg: err.Error(),
	})
}
