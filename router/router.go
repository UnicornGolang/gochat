package router

import (
	"gochat/docs"
	"gochat/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()

	// 引入 swagger 的服务端点
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.GET("/index", service.GetIndex)
	r.GET("/user/list", service.GetUserList)
	r.POST("/user/login", service.Login)
	r.POST("/user/createUser", service.CreateUser)
	r.DELETE("/user/deleteUser", service.DeleteUser)
	r.PATCH("/user/updateUser", service.UpdateUser)
	r.GET("/message/send", service.SendMessage)

	return r
}
