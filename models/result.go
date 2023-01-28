package models

// 统一的数据返回格式
type Result struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

// 请求成功失败的 code 常量
const (
	SUCCESS int = 0
	FAILURE int = -1
)

func Success(data interface{}) *Result {
	return &Result{
		Code:    SUCCESS,
		Msg:    "操作成功",
		Data:    data,
	}
}

func Failure(message string) *Result {
	return &Result{
		Code:    FAILURE,
		Msg:     message,
		Data:    nil,
	}
}
