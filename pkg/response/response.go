package response

import "github.com/gin-gonic/gin"

const (
	MaskNeedAuthor   = 8
	MaskParamMissing = 7
	StatusSuccess    = 200
	StatusError      = 200
)

type JsonResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (resp *JsonResponse) ToJson(ctx *gin.Context) {
	code := 200
	if resp.Code != StatusSuccess {
		code = resp.Code
	}
	ctx.JSON(code, resp)
}

func FailResponse(code int, message string) *JsonResponse {
	return &JsonResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}

func SuccessResponse(data ...interface{}) *JsonResponse {
	var r interface{}
	if len(data) > 0 {
		r = data[0]
	}
	return &JsonResponse{
		Code:    StatusSuccess,
		Message: "Success",
		Data:    r,
	}
}

func ErrorResponse(status int, message string, data ...interface{}) *JsonResponse {
	return &JsonResponse{
		Code:    status,
		Message: message,
		Data:    data,
	}
}

func (j *JsonResponse) WriteTo(ctx *gin.Context) {
	code := 200
	if j.Code != StatusSuccess {
		code = j.responseCode()
	}
	ctx.JSON(code, j)
}

func (j *JsonResponse) responseCode() int {
	if j.Code != StatusSuccess {
		return 200
	}
	return 200
}
