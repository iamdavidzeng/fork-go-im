package auth

import (
	messageModel "fork_go_im/im/http/models/msg"
	userModel "fork_go_im/im/http/models/user"
	"fork_go_im/pkg/helper"
	"fork_go_im/pkg/model"
	"fork_go_im/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type (
	UsersController struct{}

	UsersList struct {
		ID       uint64 `json:"id"`
		Email    string `json:"email"`
		Avatar   string `json:"avatar"`
		Name     string `json:"name"`
		Msg      string `json:"msg"`
		Status   int    `json:"status"`
		IsRead   int    `json:"is_read"`
		SendTime string `json:"send_time"`
		SendMsg  string `json:"send_msg"`
		MsgToTal int    `json:"msg_total"`
	}
)

func (*UsersController) GetUsersList(c *gin.Context) {
	name := c.Query("name")
	user := userModel.AuthUser
	var users []UsersList
	query := model.DB.Model(userModel.Users{}).Where("id <> ?", user.ID)
	if len(name) > 0 {
		query = query.Where("name like ?", "%"+name+"%")
	}
	query.Select("id", "name", "avatar", "status", "created_at").Find(&users)
	response.SuccessResponse(map[string]interface{}{
		"list": users,
	}, 200).ToJson(c)
}

func (*UsersController) ReadMessage(c *gin.Context) {
	user := userModel.AuthUser
	cA, cB := helper.ProduceChannelName(strconv.Itoa(int(user.ID)), c.Query("to_id"))
	messageModel.ReadMsg(cA, cB)
	response.SuccessResponse(gin.H{}, 200).ToJson(c)
}
