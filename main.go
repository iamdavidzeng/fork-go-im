package main

import (
	"fork_go_im/config"
	"fork_go_im/im"
	"fork_go_im/im/service"
	conf "fork_go_im/pkg/config"
	log2 "fork_go_im/pkg/log"
	"fork_go_im/pkg/pool"
	"fork_go_im/pkg/wordsfilter"
	"fork_go_im/router"

	"github.com/gin-gonic/gin"
)

func init() {
	config.Initialize()
	wordsfilter.SetTexts()
}

func main() {
	app := gin.Default()

	im.SetupPool()

	pool.AntsPool.Submit(func() {
		service.ImManager.ImStart()
	})

	router.RegisterApiRoutes(app)
	router.RegisterImRouters(app)

	app.Use(log2.Recover)
	_ = app.Run(":" + conf.GetString("app.port"))
}
