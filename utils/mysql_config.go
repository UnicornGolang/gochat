package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var (
	DB *gorm.DB
)

func InitMySQL() {
	// 自定义 SQL 日志
	customLog := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // 日志级别
			Colorful:      true,        // 日志彩色显示
		},
	)
	var err error
	DB, err = gorm.Open(mysql.Open(viper.GetString("mysql.url")), &gorm.Config{

		// 开启默认的日志
		// Logger: logger.Default.LogMode(logger.Info),

		// 设置自定义的日志
		Logger: customLog,

		// 日志的命名策略
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		fmt.Println("mysql connect error: ", err.Error())
		return
	}
	fmt.Println("##################################################")
	fmt.Println("### ---------     mysql config    ------------###")
	fmt.Println("##################################################")
	fmt.Println("config mysql: ", viper.Get("mysql"))
	fmt.Println("mysql init successful!!")
	fmt.Println("")
}
