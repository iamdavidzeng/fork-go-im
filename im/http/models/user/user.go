package user

import "fork_go_im/pkg/model"

type Users struct {
	ID              uint64 `json:"id"`
	Email           string `valid:"email" json:"email"`
	Password        string `valid:"password"`
	Avatar          string `json:"avatar"`
	Name            string `json:"name"`
	OauthType       int
	OauthId         string
	CreatedAt       string `json:"created_at"`
	PasswordConfirm string `gorm:"-" valid:"password_confirm"`
	Bio             string `json:"bio"`
}

func (Users) TableName() string {
	return "users"
}

var AuthUser *Users

func (u Users) GetAvater() string {
	if u.Avatar == "" {
		return "https://learnku.com/users/27407"
	}
	return u.Avatar
}

func SetUserStatus(id uint64, status int) {
	model.DB.Model(&Users{}).Where("id=?", id).Update("status", status)
}
