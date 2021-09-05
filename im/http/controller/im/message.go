package im

import (
	"fmt"
	userModel "fork_go_im/im/http/models/user"
	"fork_go_im/pkg/helper"
	"fork_go_im/pkg/model"
	"fork_go_im/pkg/response"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"github.com/spf13/cast"
)

type MessageController struct{}

type ImMessage struct {
	ID          uint64          `json:"id"`
	Msg         string          `json:"msg"`
	CreatedAt   string          `json:"created_at"`
	FromId      uint64          `json:"user_id"`
	ToId        uint64          `json:"send_id"`
	Channel     string          `json:"channel"`
	Status      int             `json:"is_read"`
	MsgType     int             `json:"msg_type"`
	ChannelType int             `json:"channel_type"`
	Users       userModel.Users `json:"users" gorm:"foreignKey:FromId;references:ID"`
}

func (*MessageController) InfomationHistory(c *gin.Context) {
	toId := c.Query("to_id")
	channdelType := c.DefaultQuery("channel_type", "1")
	user := userModel.AuthUser
	fromId := cast.ToString(user.ID)
	if len(toId) < 0 {
		response.FailResponse(http.StatusBadRequest, "User id cannot be null.").ToJson(c)
	}

	var MsgList []ImMessage
	channelA, channelB := helper.ProduceChannelName(fromId, toId)
	fmt.Println(channelA, channelB)
	list := model.DB.Model(ImMessage{}).
		Where("(channel = ? or channel = ?) and channel_type = ? order by created_at desc", channelA, channelB, channdelType).
		Limit(40).
		Select("id, msg, created_at, from_id, to_id, channel, msg_type").
		Find(&MsgList)
	if list.Error != nil {
		return
	}

	fromIds, _ := cast.ToUint64E(user.ID)
	for key, value := range MsgList {
		MsgList[key].CreatedAt = carbon.Parse(value.CreatedAt).SetLocale("zh-CN").DiffForHumans()
		if value.FromId == fromIds {
			MsgList[key].Status = 0
		} else {
			MsgList[key].Status = 1
		}
	}
	SortByAge(MsgList)
	response.SuccessResponse(MsgList, 200).ToJson(c)
}

func SortByAge(list []ImMessage) {
	sort.Slice(list, func(i, j int) bool {
		return list[i].ID < list[j].ID
	})
}

func (*MessageController) GetGroupMessageList(c *gin.Context) {
	toId := c.Query("to_id")
	channelType := c.DefaultQuery("channel_type", "1")
	user := userModel.AuthUser

	if len(toId) < 0 {
		response.FailResponse(500, "User id cannot be null.").ToJson(c)
	}

	var MsgList []ImMessage
	channelA := helper.ProduceChannelGroupName(toId)
	list := model.DB.Preload("Users").
		Where("channel = ? and channel_type = ? order by created_at desc", channelA, channelType).
		Limit(40).
		Select("id, msg, created_at, from_id, to_id, channel, msg_type").
		Find(&MsgList)
	if list.Error != nil {
		return
	}
	fromIds, _ := cast.ToUint64E(user.ID)
	for key, value := range MsgList {
		MsgList[key].CreatedAt = carbon.Parse(value.CreatedAt).SetLocale("zh-CN").DiffForHumans()
		if value.FromId == fromIds {
			MsgList[key].Status = 0
		} else {
			MsgList[key].Status = 1
		}
	}
	SortByAge(MsgList)
	response.SuccessResponse(MsgList, 200).ToJson(c)
}
