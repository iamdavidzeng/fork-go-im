package service

import (
	"encoding/json"
	"fmt"
	messageModel "fork_go_im/im/http/models/msg"
	"fork_go_im/pkg/helper"
	"fork_go_im/pkg/model"
	"strconv"
	"time"
)

func EnMessage(message []byte) (msg *Message) {
	err := json.Unmarshal([]byte(string(message)), &msg)
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}
	return
}

func PutGroupData(msg *Msg, isRead int, channelType int) {
	channelA := helper.ProduceChannelGroupName(strconv.Itoa(msg.ToId))

	fid := uint64(msg.FromId)
	tid := uint64(msg.ToId)
	user := messageModel.ImMessage{
		FromId:      fid,
		ToId:        tid,
		Msg:         msg.Msg,
		CreatedAt:   time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05"),
		Channel:     channelA,
		IsRead:      isRead,
		MsgType:     msg.MsgType,
		ChannelType: channelType,
	}
	model.DB.Create(&user)
}

func PutData(msg *Msg, isRead int, channelType int) {
	channelA, _ := helper.ProduceChannelName(
		strconv.Itoa(msg.FromId),
		strconv.Itoa(msg.ToId),
	)
	fid := uint64(msg.FromId)
	tid := uint64(msg.ToId)
	user := messageModel.ImMessage{
		FromId:      fid,
		ToId:        tid,
		Msg:         msg.Msg,
		CreatedAt:   time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05"),
		Channel:     channelA,
		IsRead:      isRead,
		MsgType:     msg.MsgType,
		ChannelType: channelType,
	}
	model.DB.Create(&user)
}

func GetGroupUid(groupId int) ([]GroupId, error) {
	var groups []GroupId
	err := model.DB.Table("im_group_users").
		Where("group_id=?", groupId).
		Find(&groups).
		Error
	if err != nil {
		return groups, err
	}
	return groups, nil
}
