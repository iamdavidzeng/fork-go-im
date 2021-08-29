package msg

import (
	"fmt"
	"fork_go_im/im/http/models/user"
	"fork_go_im/pkg/model"
)

type ImMessage struct {
	ID        uint64 `json:"id"`
	Msg       string `json:"msg"`
	CreatedAt string `json:"created_at"`
	FromId    uint64 `json:"user_id"`
	ToId      uint64 `json:"send_id"`
	Channel   string `json:"channel"`

	IsRead      int          `json:"is_read"`
	MsgType     int          `json:"msg_type"`
	ChannelType int          `json:"channel_type"`
	Users       []user.Users `json:"users" gorm:"foreignKey:ID;references:FromId"`
}

func (ImMessage) TableName() string {
	return "im_messages"
}

func GetOfflineMessage(id uint64) (msg *[]ImMessage) {
	list := model.DB.Where("id=?", id).Find(&msg)
	if list.Error != nil {
		fmt.Println(list.Error)
	}
	return msg
}

func ReadMsg(cA, cB string) {
	model.DB.Model(&ImMessage{}).Where("channel=? or channel=?", cA, cB).Update("is_read", 1)
}
