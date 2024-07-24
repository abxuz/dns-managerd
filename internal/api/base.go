package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type baseApi struct {
}

func (base *baseApi) Success(ctx *gin.Context, data any) {
	obj := gin.H{"errno": 0}
	if data != nil {
		obj["data"] = data
	}
	ctx.JSON(http.StatusOK, obj)
}

func (base *baseApi) Error(ctx *gin.Context, msg string) {
	obj := gin.H{"errno": 1}
	if msg != "" {
		obj["errmsg"] = msg
	}
	ctx.JSON(http.StatusOK, obj)
}
