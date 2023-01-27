package service

import (
	"fmt"
	"gochat/models"
	"html/template"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1
// @Accept json
// @Produce json
// PingExample godoc
// @Summary ping example
// @Schemes

// @Description 冒烟测试
// @Tags 首页
// @Success 200 {string} welcome.
// @Router /index [get]
func GetIndex(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "welcome!!",
	})
}

// 首页
func GetHome(c *gin.Context) {
	template, err := template.ParseFiles("index.html", "views/chat/head.html")
	if err != nil {
		fmt.Println("===", err)
		panic(err)
	}
	err = template.Execute(c.Writer, "index")
	if err != nil {
		fmt.Println("---", err)
	}
}

// 请求注册页面
func ToRegister(c *gin.Context) {
	template, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		fmt.Println("===", err)
		panic(err)
	}
	err = template.Execute(c.Writer, "register")
	if err != nil {
		fmt.Println("---", err)
	}
}

// 加载聊天页面
func ToChat(c *gin.Context){
	template, err := template.ParseFiles(
    "views/chat/index.html",
    "views/chat/head.html",
    "views/chat/main.html",
    "views/chat/foot.html",
    "views/chat/tabmenu.html",
    "views/chat/profile.html",
    "views/chat/concat.html",
    "views/chat/userinfo.html",
    "views/chat/createcom.html",
    "views/chat/group.html",
  )
	if err != nil {
		fmt.Println("===", err)
		panic(err)
	}
  userId, _ := strconv.Atoi(c.Query("userId"))
  token := c.Query("token")

  user := models.UserBasic{}
  db_user := models.GetUserById(uint(userId))
  if (db_user.ID == 0 || db_user.Identity == token ) {
    err = template.Execute(c.Writer, &user)
    if err != nil {
      fmt.Println("---", err)
    }
    return
  }

  user.ID = uint(userId)
  user.Identity = token
  err = template.Execute(c.Writer, &user)
	if err != nil {
		fmt.Println("---", err)
	}
}
