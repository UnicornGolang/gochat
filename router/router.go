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

	// 静态资源
	r.Static("/asset", "asset/")
	r.LoadHTMLGlob("views/**/*")

	// 首页
	r.GET("/", service.GetHome)
	r.GET("/index", service.GetIndex)
	r.GET("/toRegister", service.ToRegister)
	r.GET("/toChat", service.ToChat)

	// 用户服务
	r.GET("/user/list", service.GetUserList)
	r.POST("/user/login", service.Login)
	r.POST("/user/createUser", service.CreateUser)
	r.DELETE("/user/deleteUser", service.DeleteUser)
	r.PATCH("/user/updateUser", service.UpdateUser)
	r.POST("/user/searchFriends", service.SearchFriends)

	// 添加好友
	r.POST("/contact/addFriend", service.AddFriend)

	// 添加群组
	r.POST("/contact/createCommunity", service.AddCommunity)
	r.POST("/contact/loadCommunity", service.LoadCommunity)

	// 消息处理
	r.GET("/message/send", service.SendMessage)
	r.GET("/message/chat", service.SendChatMessage)

	// 附件上传
	r.POST("/attach/upload", service.Upload)

	return r
}
