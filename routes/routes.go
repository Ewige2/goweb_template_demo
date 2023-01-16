package routes

import (
	"github.com/gin-gonic/gin"
	"web_app_template/logger"
)

func RegisterRoute() *gin.Engine {
	r := gin.New()
	//注册 日志  中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello ")

	})

	return r
}
