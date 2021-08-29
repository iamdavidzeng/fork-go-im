package middleware

import (
	"errors"
	"fork_go_im/im/http/models/user"
	NewJwt "fork_go_im/pkg/jwt"
	"fork_go_im/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

var (
	token  string
	err    error
	claims *NewJwt.CustomClaims
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token = c.DefaultQuery("token", c.GetHeader("authorization"))
		token, err = ValidatedToken(token)
		if err != nil {
			response.ErrorResponse(401, err.Error()).WriteTo(c)
			c.Abort()
			return
		}
		jwt := NewJwt.NewJWT()
		claims, err = jwt.ParseToken(token)
		if err != nil {
			response.ErrorResponse(401, err.Error()).WriteTo(c)
			c.Abort()
			return
		}
		c.Set("claims", claims)
		setAuthUser(c)

		c.Next()
	}
}

func ValidatedToken(token string) (string, error) {
	if len(token) == 0 {
		return "", errors.New("Token cannot be null")
	}

	t := strings.Split(token, "Bearer ")
	if len(t) > 1 {
		return t[1], nil
	}
	return token, nil
}

func setAuthUser(c *gin.Context) {
	claims := c.MustGet("claims").(*NewJwt.CustomClaims)
	id, _ := cast.ToUint64E(claims.ID)
	user.AuthUser = &user.Users{
		ID:     id,
		Email:  claims.Email,
		Avatar: claims.Avatar,
		Name:   claims.Name,
	}
}
