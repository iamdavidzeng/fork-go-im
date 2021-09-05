package router

import (
	"fork_go_im/im/http/controller/im"
	"fork_go_im/im/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterImRouters(router *gin.Engine) {
	IMService := new(im.IMService)
	ws := router.Group("/im").Use(middleware.Auth())
	{
		ws.GET("/connect", IMService.Connect)
	}
}
