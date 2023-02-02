package main

import (
	"gochat/models"
	"gochat/router"
	"gochat/utils"
	"time"

	"github.com/spf13/viper"
)

func main() {
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	InitTimer()
	r := router.Router() // router.Router()
	// listen and server on 0.0.0.0:8080 （for windows "localhost:8080")
	r.Run(viper.GetString("server.port"))
}

// 初始化定时器
func InitTimer() {
	utils.Timer(
		time.Duration(viper.GetInt("timeout.taskDelay"))*time.Second,
		time.Duration(viper.GetInt("timeout.taskIdle"))*time.Second,
		models.CleanConnection,
		"",
	)
}
