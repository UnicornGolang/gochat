package service

import (
	"fmt"
	"gochat/models"
	"gochat/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}



// 从 redis 读取历史消息
func RedisMsg(c *gin.Context) {
  userIdA, _ := strconv.Atoi(c.PostForm("userIdA"))
  userIdB, _ := strconv.Atoi(c.PostForm("userIdB"))
  start, _ := strconv.Atoi(c.PostForm("start"))
  end, _ := strconv.Atoi(c.PostForm("end"))
  isRev, _ := strconv.ParseBool(c.PostForm("isRev"))
  res := models.RedisMsg(uint(userIdA), uint(userIdB), int64(start), int64(end), isRev)
  utils.RespOKList(c.Writer, res, len(res))
}

// 获取 ws 通信 连接
func SendChatMessage(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

// 发送消息
func SendMessage(c *gin.Context) {
	ws, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)

	MsgHandler(ws, c)

}

// 消息处理
func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	msg, err := utils.Subscribe(c, utils.PublishKey)
	if err != nil {
		fmt.Println(err)
	}
	tm := time.Now().Format("2006-01-02 15:04:05")
	m := fmt.Sprintf("[ws]\t[%s]:\t %s", tm, msg)
	fmt.Println("发送消息:>> ")
	err = ws.WriteMessage(1, []byte(m))
	if err != nil {
		fmt.Println(err)
	}
}
