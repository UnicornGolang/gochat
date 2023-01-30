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

func GetCommunityByNameOrId(comId string) *Community {
  community := Community{}
  utils.DB.Where("id = ? or name = ?", comId, comId).Find(&community)
  return &community
}

func AddCommunity(community *Community) {
	utils.DB.Create(community)
}

func LoadCommunity(userId uint) []*Community {
	data := make([]*Contact, 0)
  comminityIds := make([]uint,0)
	utils.DB.Where("owner_id = ? and type = 2", userId).Find(&data)
  for _, v := range data {
    comminityIds = append(comminityIds, v.TargetId)
  }
  comminities := make([]*Community, 10)
  utils.DB.Where("id in ?", comminityIds).Find(&comminities)
	return comminities
}
