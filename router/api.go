package router

import (
	Auth "fork_go_im/im/http/controller/auth"
	"fork_go_im/im/http/controller/im"
	"fork_go_im/im/http/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterApiRoutes(router *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{
		"tus-resumable",
		"upload-length",
		"upload-metadata",
		"chache-control",
		"x-requested-with",
		"*",
	}
	router.Use(cors.New(config))
	weibo := new(Auth.WeiBoController)
	auth := new(Auth.AuthController)
	users := new(Auth.UsersController)
	sm := new(im.SmAPICrontrooler)
	uploads := new(im.UploadCrotroller)
	group := new(im.GroupController)
	im := new(im.MessageController)

	apiRouter := router.Group("/api")
	apiRouter.Group("")
	{
		apiRouter.POST("/login", auth.Login)                 // account login
		apiRouter.GET("/WeiBoCallback", weibo.WeiBoCallback) // TODO: weibo auth
		apiRouter.GET("/getApiToken", sm.GetApiToken)        // get sm token
		apiRouter.Use(middleware.Auth())
		{
			apiRouter.POST("/me", auth.Me)                  // get user info
			apiRouter.GET("/UsersList", users.GetUsersList) // get user list

			apiRouter.GET("/InformationHistory", im.InfomationHistory)    // get message list
			apiRouter.GET("/GetGroupMessageList", im.GetGroupMessageList) // get message list
			apiRouter.POST("/UploadImg", sm.UploadImg)                    // upload image
			apiRouter.POST("/UploadBVoiceFile", uploads.UploadVoiceFile)  // upload voice file
			apiRouter.GET("/ReadMessage", users.ReadMessage)              // read message

			apiRouter.GET("/GetGroupList", group.List)   // get group list
			apiRouter.POST("/CreateGroup", group.Create) // add group
		}

	}
}
