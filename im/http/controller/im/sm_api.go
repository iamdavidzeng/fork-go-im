package im

import (
	"encoding/json"
	"fmt"
	"fork_go_im/im/utils"
	"fork_go_im/pkg/config"
	log2 "fork_go_im/pkg/log"
	"fork_go_im/pkg/redis"
	"fork_go_im/pkg/response"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	username = config.GetString("app.sm_name")
	password = config.GetString("app.sm_password")
	smToken  = config.GetString("app.sm_token")
)

type (
	SmAPICrontrooler struct{}

	Data struct {
		Token string `json:"token"`
	}
	ResponseData struct {
		Success   bool   `json:"success"`
		Code      string `json:"code"`
		Message   string `json:"message"`
		Data      Data   `json:"data"`
		RequestId string `json:"request_id"`
	}
)

func (*SmAPICrontrooler) GetApiToken(c *gin.Context) {
	stringCmd := redis.RedisDB.Get("sm_token")
	if len(stringCmd.Val()) != 0 {
		resp := new(ResponseData)
		resp.Code = "success"
		resp.Data.Token = stringCmd.Val()
		resp.Success = true
		fmt.Println(resp)
		c.JSON(200, resp)
		return
	}
	data := url.Values{"username": {username}, "password": {password}}
	j, err := http.PostForm("https://sm/ms/api/v2/token", data)
	log2.Warning(err.Error())
	defer j.Body.Close()

	bodyC, _ := ioutil.ReadAll(j.Body)
	resp := new(ResponseData)
	json.Unmarshal(bodyC, resp)
	if resp.Success {
		response.FailResponse(500, resp.Message)
		return
	}
	redis.RedisDB.Set("sm_token", resp.Data.Token, 1*time.Hour)
	c.JSON(200, resp)
}

type (
	DataSuccess struct {
		FileId    int    `json:"file_id"`
		Width     int    `json:"width"`
		Height    int    `json:"height"`
		FileName  string `json:"filename"`
		Storename string `json:"storename"`
		Size      int    `json:"size"`
		Path      string `json:"path"`
		Hash      string `json:"hash"`
		Url       string `json:"url"`
		Delete    string `json:"delete"`
		Page      string `json:"page"`
	}
	ResponseUploadData struct {
		Success   bool        `json:"success"`
		Code      string      `json:"code"`
		Message   string      `json:"message"`
		Data      DataSuccess `json:"data"`
		RequestId string      `json:"request_id"`
	}
)

func (*SmAPICrontrooler) UploadImg(c *gin.Context) {
	file, _ := c.FormFile("Smfile")
	dir := utils.GetCurrentDirectory()
	path := dir + "/docs/" + file.Filename
	err := c.SaveUploadedFile(file, path)
	log2.LogError(err)
	header := new(utils.Header)
	header.Authorizatoin = "Authorization"
	header.Token = smToken
	resp, err := utils.PostFile(path, "https://sm.ms/api/v2/upload", header)
	log2.LogError(err)
	bodyC, _ := ioutil.ReadAll(resp.Body)
	data := new(ResponseUploadData)
	json.Unmarshal(bodyC, data)
	c.JSON(200, data)
}
