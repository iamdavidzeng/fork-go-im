package group_user

import (
	"fork_go_im/im/http/models/user"
	"fork_go_im/pkg/model"
	"strconv"
	"time"
)

type ImGroupUsers struct {
	ID        uint64 `json:"id"`
	UserId    uint64 `json:"user_id"`
	CreatedAt string `json:"created_at"`
	GroupId   uint64 `json:"group_id"`
	Remark    string `json:"remark"`
	Avatar    string `json:"avatar"`
}

func (ImGroupUsers) TableName() string {
	return "im_group_users"
}

func CreatedAll(userIds map[string]string, groupId uint64, uId uint64) (err error) {
	var groupUsers = make([]*ImGroupUsers, len(userIds)+1)
	var userId = make([]int, len(userIds)+1)
	userId = append(userId, int(uId))
	for _, value := range userIds {
		valueNum, _ := strconv.Atoi(value)
		userId = append(userId, valueNum)
	}
	var users []user.Users

	err = model.DB.Where("id in (?)", userId).Find(&users).Error
	if err != nil {
		return err
	}

	var i = 0
	for _, value := range users {
		groupUsers[i] = &ImGroupUsers{
			UserId:    value.ID,
			GroupId:   groupId,
			CreatedAt: time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05"),
			Remark:    value.Email,
			Avatar:    value.Avatar,
		}
		i++
	}

	err = model.DB.Model(&ImGroupUsers{}).Create(&groupUsers).Error
	if err != nil {
		return err
	}
	return nil
}
