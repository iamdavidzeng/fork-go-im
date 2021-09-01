package im

import (
	"encoding/json"
	"fmt"
	"fork_go_im/im/http/models/group"
	"fork_go_im/im/http/models/group_user"
	userModel "fork_go_im/im/http/models/user"
	"fork_go_im/im/http/validates"
	"fork_go_im/pkg/helper"
	log2 "fork_go_im/pkg/log"
	"fork_go_im/pkg/model"
	"fork_go_im/pkg/response"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type (
	GroupController struct{}
	Groups          struct {
		GroupId string `json:"group_id"`
	}
)

func (*GroupController) List(c *gin.Context) {
	user := userModel.AuthUser
	var groupIds []Groups
	err := model.DB.Table("im_group_users").
		Where("user_id=?", user.ID).
		Group("group_id").
		Find(&groupIds).Error
	if err != nil {
		fmt.Println(err)
	}
	v := reflect.ValueOf(groupId)
	groupSlice := make([]string, v.Len())
	for key, value := range groupIds {
		groupSlice[key] = value.GroupId
	}
	list, err := group.GetGroupUserList(groupSlice)
	if err != nil {
		log2.Warning(err.Error())
		response.FailResponse(http.StatusInternalServerError, "Internal server error.")
	}
	response.SuccessResponse(list).ToJson(c)
}

func (*GroupController) Create(c *gin.Context) {
	user := userModel.AuthUser

	groups := validates.CreateGroupParams{
		GroupName: c.PostForm("group_name"),
		UserId:    c.PostFormMap("user_id"),
	}
	fmt.Println(groups)
	rules := govalidator.MapData{
		"group_name": []string{"required", "between:2,20"},
		//"user_id": []string{"required"},
	}
	opts := govalidator.Options{
		Data:          &groups,
		Rules:         rules,
		TagIdentifier: "valid",
	}
	errs := govalidator.New(opts).ValidateStruct()

	if len(errs) > 0 {
		data, _ := json.MarshalIndent(errs, "", " ")
		result := helper.JsonToMap(data)
		response.ErrorResponse(http.StatusBadRequest, "Invalid user input.", result).ToJson((c))
		return
	}
	if len(groups.UserId) > 50 {
		response.ErrorResponse(http.StatusBadRequest, "Group's capacity must to be less than 50.").ToJson(c)
	}
	id, err := group.Created(user.ID, groups.GroupName)
	if err != nil {
		fmt.Println("Exception.")
		response.ErrorResponse(http.StatusInternalServerError, "Exception when creating.").ToJson(c)
		return
	}
	err = group_user.CreatedAll(groups.UserId, id, user.ID)
	if err != nil {
		response.ErrorResponse(http.StatusInternalServerError, "Exception wehn creating.").ToJson(c)
		return
	}
	response.SuccessResponse().ToJson(c)
	return
}

func (*GroupController) RemoveGroup(c *gin.Context) {
	groupId := c.PostForm("group_id")
	if len(groupId) == 0 {
		response.ErrorResponse(http.StatusBadRequest, "Invalid user input.").ToJson(c)
		return
	}
	model.DB.Where("id=?", groupId).Delete(&group.ImGroups{})
	model.DB.Where("group=?", groupId).Delete(&group_user.ImGroupUsers{})
	response.SuccessResponse().ToJson(c)
	return
}

func (*GroupController) DeleteUser(c *gin.Context) {
	// TODO: Implements it
}
