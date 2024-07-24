package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BasicAuth(username string, password string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_username, _password, ok := ctx.Request.BasicAuth()
		if ok && username == _username && password == _password {
			ctx.Next()
			return
		}

		ctx.Header("WWW-Authenticate", "Basic realm=Authorization Required")
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
}
