package service

import (
	"fmt"
	"gochat/models"
	"gochat/utils"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// @Summary 加入某个群聊
// @Description 加入群聊
// @Tags 用户服务
// @Produce json
// @Success 200 {json} Result
// @Router /contact/joinGroup [post]
func JoinGroup(c *gin.Context) {
	userId, _ := strconv.Atoi(c.PostForm("userId"))
	dbuser := models.GetUserById(uint(userId))
	if dbuser.ID == 0 {
		c.JSON(http.StatusOK, models.Failure("用户未登录"))
		return
	}
	comId := c.PostForm("comId")
	if len(comId) < 1 {
		c.JSON(http.StatusOK, models.Failure("未获取到需要加入的群组的信息"))
		return
	}
	community := models.GetCommunityByNameOrId(comId)
	if community.Name == "" {
		c.JSON(http.StatusOK, models.Failure("未找到对应的群组信息"))
		return
	}
	contact := models.Contact{}
	contact.OwnerId = uint(userId)
	contact.TargetId = community.ID
	models.JoinCommunity(&contact)
	c.JSON(http.StatusOK, models.Success(nil))
}

// @Summary 加载用户群组
// @Description 加载群组群聊
// @Tags 用户服务
// @Produce json
// @Success 200 {json} []Community
// @Router /contact/loadCommunity [post]
func LoadCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.PostForm("ownerId"))
	if ownerId == 0 {
		c.JSON(http.StatusOK, models.Failure("用户需先登录"))
		return
	}
	dbuser := models.GetUserById(uint(ownerId))
	if dbuser.ID == 0 {
		c.JSON(http.StatusOK, models.Failure("用户不存在"))
		return
	}
	communities := models.LoadCommunity(dbuser.ID)
	utils.RespOKList(c.Writer, communities, len(communities))
}

// @Summary 创建群聊
// @Description 创建群聊
// @Tags 用户服务
// @Produce json
// @Success 200 {json} Result
// @Router /contact/addCommunity [post]
func AddCommunity(c *gin.Context) {

	ownerId, _ := strconv.Atoi(c.PostForm("ownerId"))
	name := c.PostForm("name")
	icon := c.PostForm("icon")
	desc := c.PostForm("desc")
	if ownerId == 0 {
		c.JSON(http.StatusOK, models.Failure("用户需先登录"))
		return
	}
	if len(name) < 1 {
		c.JSON(http.StatusOK, models.Failure("群组名称不能为空"))
		return
	}
	dbuser := models.GetUserById(uint(ownerId))
	if dbuser.ID == 0 {
		c.JSON(http.StatusOK, models.Failure("用户不存在"))
		return
	}
  tx := utils.DB.Begin()
  defer func() {
    if r := recover(); r != nil {
      fmt.Println("add community err : ", r)
      tx.Rollback()
    }
  }()
	community := &models.Community{
		OwnerId: uint(ownerId),
		Name:    name,
		Img:     icon,
		Desc:    desc,
	}
  models.AddCommunity(community)

  // 添加团队成员
  concact := models.Contact{}
  concact.TargetId = community.ID
  concact.OwnerId = uint(ownerId)
  concact.Type = 2
  if err := utils.DB.Create(&concact).Error; err != nil {
    fmt.Println("create Community err ", err)
    tx.Rollback()
  }
	c.JSON(http.StatusOK, models.Success(nil))
  tx.Commit()
}

// @Summary 添加好友
// @Description 添加好友
// @Tags 用户服务
// @Produce json
// @Success 200 {json} Result
// @Router /user/addFriend [post]
func AddFriend(c *gin.Context) {
	userId := c.PostForm("userId")
	friendName := c.PostForm("targetName")
	if userId == "" || friendName == "" {
		c.JSON(http.StatusOK, models.Failure("未获取到请求参数"))
		return
	}
	OwnId, _ := strconv.Atoi(userId)
	// 查找对应的名称的好友
	friend := models.GetUserByName(friendName)
	if friend.ID == 0 {
		c.JSON(http.StatusOK, models.Failure("未找到指定的好友"))
		return
	}
	if uint(OwnId) == friend.ID {
		c.JSON(http.StatusOK, models.Failure("不能添加自己为好友"))
		return
	}
	models.AddFriendRelation(uint(OwnId), friend.ID)
	c.JSON(http.StatusOK, models.Success(nil))
}

// @Summary 获取用户列表
// @Description 获取用户列表
// @Tags 用户服务
// @Produce json
// @Success 200 {list} []UserBasic
// @Router /user/searchFriends [post]
func SearchFriends(c *gin.Context) {
	id := c.PostForm("userId")
	if id == "" {
		c.JSON(http.StatusOK, models.Failure("为获取到查询参数"))
		return
	}
	userId, _ := strconv.Atoi(id)
	frinds := models.SearchFriends(uint(userId))
	utils.RespOKList(c.Writer, frinds, len(frinds))

	//c.JSON(http.StatusOK, models.Success(frinds))
}

// @Summary 获取用户列表
// @Description 获取用户列表
// @Tags 用户服务
// @Produce json
// @Success 200 {list} []UserBasic
// @Router /user/list [get]
func GetUserList(c *gin.Context) {
	data := models.GetUserList()
	fmt.Println(data)
	// c.JSON(200, gin.H{
	// 	 "message": data,
	// })
	c.JSON(http.StatusOK, models.Success(&data))
}

// @Summary 用户登录
// @Description 用户登录服务
// @Tags 用户服务
// @Produce json
// @Param name formData string false "用户名"
// @Param password formData string false "登录密码"
// @Success 200 {string} json{"code", "message"}
// @Router /user/login [post]
func Login(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")
	user := models.GetUserByName(name)

	// 校验用户名是否存在
	if user.Name == "" {
		c.JSON(http.StatusOK, models.Failure("用户不存在"))
		return
	}
	// 校验密码是否正确
	if !utils.ValidatePassword(password, user.Salt, user.Password) {
		c.JSON(http.StatusOK, models.Failure("密码错误"))
		return
	}
	identity := fmt.Sprintf("%d:%s:%s", time.Now().Unix(), user.Name, user.Password)
	user.Identity = utils.Md5Encode(identity)
	models.StoreIdentity(user)

	// c.JSON(200, gin.H{
	// 	"message": "登录成功",
	// })
	// 敏感信息去除
	user.Password = ""
	user.Salt = ""
	c.JSON(http.StatusOK, models.Success(&user))
}

// @Summary  创建用户
// @Description 创建用户
// @Tags 用户服务
// @Produce json
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code", "message"}
// @Router /user/createUser [post]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.PostForm("name")
	password := c.PostForm("password")
	identity := c.PostForm("identity")

	if user.Name == "" || password == "" {
		c.JSON(http.StatusOK, models.Failure("用户名和密码不能为空"))
		return
	}

	if password != identity {
		c.JSON(http.StatusOK, models.Failure("两次输入的明码不同"))
		return
	}
	data := models.GetUserByName(user.Name)
	if data.Name != "" {
		c.JSON(http.StatusOK, models.Failure("用户已存在"))
		return
	}
	// 获取随机数作为盐值
	salt := fmt.Sprintf("%06d", rand.Int31())
	user.Password = utils.MakePassword(password, salt)
	user.Salt = salt
	models.CreateUser(user)
	c.JSON(http.StatusOK, models.Success(nil))
}

// @Summary  修改用户
// @Description 修改用户
// @Tags 用户服务
// @Produce json
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code", "message"}
// @Router /user/updateUser [patch]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	oldUser := models.GetUserById(uint(id))
	if oldUser.ID == 0 {
		c.JSON(http.StatusOK, models.Failure("用户不存在"))
		return
	}
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.Password = utils.MakePassword(c.PostForm("password"), oldUser.Salt)
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	// 数据格式校验，电话，邮箱
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, models.Failure("数据格式错误"))
		return
	}

	models.UpdateUser(user)
	c.JSON(http.StatusOK, models.Success(nil))
}

// @Summary  删除用户
// @Description 删除用户
// @Tags 用户服务
// @Produce json
// @param id query string false "id"
// @Success 200 {string} json{"code", "message"}
// @Router /user/deleteUser [delete]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	// c.JSON(200, gin.H{
	// 	"message": "删除成功",
	// })
	c.JSON(http.StatusOK, models.Success(nil))
}
