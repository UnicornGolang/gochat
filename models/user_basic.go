package models

import (
	"gochat/utils"
	"time"

	"gorm.io/gorm"
)

// 可以使用 `` 符号定义数据模型的格式化信息, 其中
// gorm: 规定了与数据库中的对应字段
// json:  限定了以 json 格式序列化数据模型的时候的输出的字段格式
// valid: 限定了数据格式校验的规则

// 时间类型对应的 time.Time 在数据库显示为时间类
type UserBasic struct {
	gorm.Model
	Name          string
	Password      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9})"`
	Email         string `valid:"email"`
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
	return data
}

// 根据 id 获取用户
func GetUserById(id uint) UserBasic {
	user := UserBasic{}
	utils.DB.Find(&user, id)

	return user
}

// 名称查找
func GetUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name=?", name).First(&user)
	return user
}

// 电话查找
func GetUserByPhone(phone string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("phone=?", phone).First(&user)
	return user
}

// email 查找
func GetUserByEmail(email string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("email=?", email).First(&user)
	return user
}

func CreateUser(user UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}

func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{
		Name:     user.Name,
		Password: user.Password,
	})
}

func StoreIdentity(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{
		Identity: user.Identity,
	})
}
