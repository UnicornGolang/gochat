package main

import (
	"gochat/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 获取数据库连接
	db, err := gorm.Open(
		mysql.Open("root:root@tcp(172.21.96.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"),
		&gorm.Config{},
	)
	if err != nil {
		panic("Fail to connect database")
	}
	// 迁移 schema
	db.AutoMigrate(&models.UserBasic{})

	// create
	user := &models.UserBasic{}
	user.Name = "alex"
	user.Phone = "12200908708"
	user.IsLogout = true
	user.Identity = "v353535"
	db.Create(user)

	// read
	println("user1:", db.First(user, 30))    // 根据整形主键查找
	db.First(user, "Phone=?", "12200908708") // 查找 Phone 字段值为 12200908708 的记录

	// 更新多个字段, 更新传入了值的字段
	db.Model(user).Updates(&models.UserBasic{Identity: "v22342", IsLogout: true})
	db.Model(user).Updates(map[string]interface{}{"Identity": "v22345", "IsLogout": true})

	// Update - 将 product 的 IsLogOut 更新为 false
	db.Model(user).Update("IsLogout", false)

	// 删除
	db.Delete(user, 1)

}
