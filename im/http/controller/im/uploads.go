package im

import (
	"fork_go_im/im/utils"
	"fork_go_im/pkg/config"
	"fork_go_im/pkg/response"

	"github.com/gin-gonic/gin"
)

type UploadCrotroller struct{}

var ym = config.GetString("app.ym")

func (*UploadCrotroller) UploadVoiceFile(c *gin.Context) {
	voice, _ := c.FormFile("voice")
	dir := utils.GetCurrentDirectory()
	path := dir + "/voice/" + voice.Filename
	c.SaveUploadedFile(voice, path)
	response.SuccessResponse(map[string]interface{}{
		"url": ym + "voice/" + voice.Filename,
	}).ToJson(c)
	return
}
