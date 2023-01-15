package models

import (fatih/set.v0"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	FromId   string // 消息的发送者
	TargetId string // 消息的接收者
	Type     string // 消息的类型, 私发，群发，广播
	Media    int    // 消息内容的类型，文字，图片，音频
	Content  string // 消息的内容
	Pic      string // 图片的地址
	Url      string // 图片的 Url
	Desc     string // 文件描述
	Amount   int    // 内容的大小
}

type Node struct {
	//
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 存储的连接与用户的映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁控制 map 的安全性问题
var rwLock sync.RWMutex

func Chat(writer *http.ResponseWriter, request *http.Request) {
	// 校验 token
	query := request.URL.Query()
  userId := query.Get("userId")
  msgType := query.Get("type")
  targetId := query.Get("targetId")
  context := query.Get("context")
  isvalid := true  // 校验 token
  conn, err := (&websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
      return invalid
    },
  }).Upgrade(writer, request, nil)

  // 获取连接出错
  if err != nil {
    fmt.println("err: ", err)
    return
  }
  // 获取连接
  node := &Node {
    Conn: conn,
    DataQueue: make(char []byte, 50),
    GroupSets: set.New(set.ThreadSafe),
  }
  // 用户关系
  // userid 
  rwLock.
}

// 返回表名
func (table *Message) TableName() string {
	return "message"

}
