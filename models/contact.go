package models

import (
	"gochat/utils"

	"gorm.io/gorm"
)

// 人员关系
type Contact struct {
  gorm.Model
	OwnerId  uint   // 关系的主体
	TargetId uint   // 关联的主体
	Type     int    // 对应的关系类型 1 是好友，2 是群组， 3 备用
	Desc     string // 关系的描述
}

func (table *Contact) TableName() string {
	return "contact"
}

// 查找某个用户的所有朋友
func SearchFriends(userId uint)([]*UserBasic){
  contacts := make([]*Contact, 0)
  utils.DB.Where("owner_id = ? and type = 1", userId).Find(&contacts)
  frindIds := make([]uint, 0)
  for _, v := range contacts {
    frindIds = append(frindIds, uint(v.TargetId))
  }
  frinds := make([]*UserBasic, 0)
  utils.DB.Where("id in ?", frindIds).Find(&frinds)
  return frinds
}
