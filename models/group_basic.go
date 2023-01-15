package models

import "gorm.io/gorm"

// 群聊
type GroupBasic struct {
	gorm.Model
	Name    string // 群名称
	OwnerId uint   // 群拥有人
	Type    int    // 类型
	Icon    string // 群头像
	Desc    string // 描述信息
}

func (table *GroupBasic) TableName() string {
	return "groupbasic"
}
