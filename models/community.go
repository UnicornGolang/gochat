package models

import (
	"gochat/utils"

	"gorm.io/gorm"
)

// Community : 群聊的模型类
type Community struct {
	gorm.Model
	OwnerId uint
	Name    string
	Img     string
	Desc    string
}

// 返回当前模型类
func (table *Community) TableName() string {
	return "community"
}

func AddCommunity(community *Community) {
	utils.DB.Create(community)
}

func LoadCommunity(userId uint) []*Community {
  data := make([]*Community, 10)
  utils.DB.Where("owner_id = ?", userId).Find(&data)
  return data
}
