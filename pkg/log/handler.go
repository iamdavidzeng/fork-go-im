package log

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			errors := errorToString(r)

			Warning(errors)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "server error",
				"data": nil,
			})

			c.Abort()
		}
	}()
	c.Next()
}

func errorToString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	default:
		return r.(string)
	}
}
