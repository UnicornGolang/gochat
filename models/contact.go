package models

import "gorm.io/gorm"

// 人员关系
type Contact struct {
  gorm.Model
	OwnerId  uint   // 关系的主体
	TargetId uint   // 关联的主体
	Type     int    // 对应的关系类型
	Desc     string // 关系的描述
}

func (table *Contact) TableName() string {
	return "contact"
}
