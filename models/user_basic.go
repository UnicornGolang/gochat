package models

import (
	"fmt"
	"gochat/utils"
	"time"

	"gorm.io/gorm"
)

// 可以使用 `` 符号定义数据模型的格式化信息
// 其中 gorm: 规定了与数据库中的对应字段
// json 限定了以 json 格式序列化数据模型的时候的输出的字段格式

// 时间类型对应的 time.Time 在数据库显示为时间类
type UserBasic struct {
	gorm.Model
	Name          string
	Password      string
	Phone         string
	Email         string
	Identity      string
	ClientIp      string
	ClientPort    string
	LoginTime     time.Time
	HeartbeatTime time.Time
	LogoutTime    time.Time `gorm:"column:logout_time" json:"logout_time"`
	IsLogout      bool
	DeviceInfo    string
	Salt          string
	Avatar        string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

func GetUser() *UserBasic {
  user := UserBasic{}
  utils.DB.Find(&user)

  return &user
}
