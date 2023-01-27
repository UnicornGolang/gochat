package models

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	FromId   int64  // 消息的发送者
	TargetId int64  // 消息的接收者
	Type     int    // 消息的类型, 私发，群发，广播
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

func Chat(writer http.ResponseWriter, request *http.Request) {
	// 校验 token
	query := request.URL.Query()
	id := query.Get("userId")
	userId, _ := strconv.ParseInt(id, 10, 64)
	fmt.Println("USERID", userId)
	//msgType := query.Get("type")
	//targetId := query.Get("targetId")
	//context := query.Get("context")
	isvalid := true // 校验 token
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalid
		},
	}).Upgrade(writer, request, nil)

	// 获取连接出错
	if err != nil {
		fmt.Println("[Conn Err]===:", err)
		return
	}
	// 获取连接
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 1),
		GroupSets: set.New(set.ThreadSafe),
	}
	// 用户关系
	// userid  跟 node 绑定，并加锁
	rwLock.Lock()
	clientMap[userId] = node
	rwLock.Unlock()

	// 完成发送逻辑
	go sendProc(node)

	// 完成接收逻辑
	go recvProc(node)

}

// 接收发送给当前用户的数据, 接收到消息后通过 websocket 写出
func sendProc(node *Node) {
	for data := range node.DataQueue {
		fmt.Println("[Rece Data]===", string(data))
		err := node.Conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println("[Cliend Err]===", err)
			return
		}
		broadMsg(data)
		fmt.Println("[ws]<<<<", string(data))
	}
}

var udpsendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendChan <- data
}

func init() {
	go udpSendProc()
	go udpRecvProc()
}

// 完成 ddp 数据发送协程
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	for data := range udpsendChan {
		_, err := con.Write(data)
		if err != nil {
			fmt.Println("[DataSend Err]===", err)
			return
		}
	}
}

func udpRecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println("[Read Err]===", err)
			return
		}
		dispatch(buf[0:n])
	}
}

func dispatch(data []byte) {
	msg := Message{}
	fmt.Println(string(data))
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	msg.Content = string(data)
	switch msg.Type {
	case 1: // 私发
		fmt.Println("[私聊]===", msg.TargetId)
		sendMsg(&msg)
		//case 2: // 群发
		//case 3: // 广播
		//case 4:
	}
}

// 发送消息
func sendMsg(msg *Message) {
	rwLock.RLock()
	node, ok := clientMap[msg.TargetId]
	rwLock.RUnlock()
	if ok {
		node.DataQueue <- []byte(msg.Content)
	}
}

// 返回表名
func (table *Message) TableName() string {
	return "message"
}
