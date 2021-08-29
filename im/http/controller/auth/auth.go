package auth

import (
	userModel "fork_go_im/im/http/models/user"
	"fork_go_im/im/http/validates"
	"fork_go_im/pkg/config"
	"fork_go_im/pkg/helper"
	"fork_go_im/pkg/jwt"
	"fork_go_im/pkg/model"
	"fork_go_im/pkg/response"
	"strconv"
	"time"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type (
	AuthController  struct{}
	WeiBoController struct{}
	Me              struct {
		ID             uint64 `json:"id"`
		Name           string `json:"name"`
		Avatar         string `json:"avatar"`
		Email          string `json:"email"`
		Token          string `json:"token"`
		ExpirationTime int64  `json:"expiration_time"`
	}
)

func (*AuthController) Me(c *gin.Context) {
	user := userModel.AuthUser
	response.SuccessResponse(user, 200).ToJson(c)
}

func (a *AuthController) Login(c *gin.Context) {
	var (
		params validates.LoginParams
		users  userModel.Users
	)
	c.ShouldBind(&params)

	model.DB.Model(&userModel.Users{}).Where("name = ?", params.Name).Find(&users)
	if users.ID == 0 {
		response.FailResponse(403, "user not exists").ToJson(c)
		return
	}
	if !helper.ComparePasswords(users.Password, params.Password) {
		response.FailResponse(403, "account or password is wrong").ToJson(c)
		return
	}
	generateToken(c, &users)
}

func generateToken(c *gin.Context, user *userModel.Users) {
	signKey := config.GetString("app.jwt.sign_key")
	expirationTime := config.GetInt64("app.jwt.expiration_time")

	j := &jwt.JWT{[]byte(signKey)}
	claims := jwt.CustomClaims{
		strconv.FormatUint(user.ID, 10),
		user.Name,
		user.Avatar,
		user.Email,
		jwtGo.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,
			ExpiresAt: time.Now().Unix() + expirationTime,
			Issuer:    signKey,
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		response.FailResponse(403, "failed to issue token.").ToJson(c)
		return
	} else {
		data := Me{
			ID:             user.ID,
			Name:           user.Name,
			Avatar:         user.Avatar,
			Email:          user.Email,
			Token:          token,
			ExpirationTime: expirationTime,
		}
		response.SuccessResponse(data, 200).ToJson(c)
		return
	}

}
