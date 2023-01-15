package models

// 统一的数据返回格式
type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 请求成功失败的 code 常量
const (
	SUCCESS int = 0
	FAILURE int = -1
)

func Success(data ...interface{}) *Result {
	return &Result{
		Code:    SUCCESS,
		Message: "操作成功",
		Data:    data[0],
	}
}


func Failure(message string) *Result {
	return &Result{
		Code:    FAILURE,
		Message: message,
		Data:    nil,
	}
}
