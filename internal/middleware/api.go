package middleware

import (
	"net/http"

	"github.com/abxuz/dns-manager/internal/service"
	"github.com/gin-gonic/gin"
)

var Api = &mApi{}

type mApi struct {
}

func (m *mApi) BasicAuth(c *gin.Context) {
	cfg := service.Config.GetCachedConfig()
	if cfg.App.Auth == nil {
		c.Next()
		return
	}

	auth := cfg.App.Auth
	username, password, ok := c.Request.BasicAuth()
	authOk := ok && (username == auth.Username) && (password == auth.Password)
	if authOk {
		c.Next()
		return
	}

	c.Header("WWW-Authenticate", "Basic realm=Authorization Required")
	c.AbortWithStatus(http.StatusUnauthorized)
}

func (m *mApi) ApiResponse(c *gin.Context) {
	c.Next()

	err := c.Errors.Last()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"errno":  1,
			"errmsg": err.Error(),
		})
		return
	}

	obj := gin.H{
		"errno": 0,
	}
	data, exists := c.Get("data")
	if exists {
		obj["data"] = data
	}
	c.JSON(http.StatusOK, obj)
}
