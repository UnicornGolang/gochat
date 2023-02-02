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
func SearchFriends(userId uint) []*UserBasic {
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

// 该方法主要在于演示事务的控制
func AddFriendRelation(userId uint, friendId uint) {
	// 开启事务，保证好友数据的一致性
	tx := utils.DB.Begin()
	// 保险起见，在事务开启之后，不管最终发生了什么异常，都会进行回滚操作
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			return
		}
	}()
	mcontact := &Contact{}
	utils.DB.Where("owner_id=? and type=1 and target_id=?", userId, friendId).Find(mcontact)
	if mcontact.OwnerId == 0 || mcontact.TargetId == 0 {
		contact := Contact{
			OwnerId:  userId,
			TargetId: friendId,
			Type:     1,
			Desc:     "",
		}
		utils.DB.Create(&contact)
	}
	ncontact := &Contact{}
	utils.DB.Where("owner_id=? and type=1 and target_id=?", friendId, userId).Find(ncontact)
	if ncontact.OwnerId == 0 || ncontact.TargetId == 0 {
		friendRelation := Contact{
			OwnerId:  friendId,
			TargetId: userId,
			Type:     1,
			Desc:     "",
		}
		if err := utils.DB.Create(&friendRelation); err != nil {
			tx.Rollback()
			return
		}
	}
	tx.Commit()
}

func JoinCommunity(contact *Contact) {
	_contact := Contact{}

	utils.DB.Where(
		"owner_id = ? and target_id = ? and type = 2",
		contact.OwnerId,
		contact.TargetId,
	).Find(&_contact)

	// 没有加入群聊则加入, 否则不做任何操作
	if _contact.ID == 0 {
		contact.Type = 2
		utils.DB.Create(&contact)
	}
}

// 获取群组内的所有用户
func SearchUserByGroupId(communityId uint) []uint {
	contacts := make([]Contact, 0)
	userIds := make([]uint, 0)
	utils.DB.Where("target_id = ? and type = 2", communityId).Find(&contacts)
	for _, v := range contacts {
		userIds = append(userIds, v.OwnerId)
	}
	return userIds
}
