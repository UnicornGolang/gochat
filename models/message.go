package models

import (
	"context"
	"encoding/json"
	"fmt"
	"gochat/utils"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	UserId     uint   // 消息的发送者
	TargetId   uint   // 消息的接收者
	Type       int    // 消息的类型, 1.私发，2.群发，3.广播
	Media      int    // 消息内容的类型，文字，图片，音频
	Content    string // 消息的内容
	CreateTime uint64 // 创建时间
	ReadTime   uint64 // 读取时间
	Pic        string // 图片的地址
	Url        string // 图片的 Url
	Desc       string // 文件描述
	Amount     int    // 内容的大小
}

type Node struct {
	//
	Conn          *websocket.Conn
	Addr          string
	HeartbeatTime uint64
	LoginTime     uint64
	DataQueue     chan []byte
	GroupSets     set.Interface
}

// 存储的连接与用户的映射关系
var clientMap map[uint]*Node = make(map[uint]*Node, 0)

// 读写锁控制 map 的安全性问题
var rwLock sync.RWMutex

func init() {
	go udpSendProc()
	go udpRecvProc()
}

func Chat(writer http.ResponseWriter, request *http.Request) {
	// 校验 token
	query := request.URL.Query()
	id := query.Get("userId")
	userId, _ := strconv.Atoi(id)
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
	currentTime := uint64(time.Now().Unix())
	// 获取连接
	node := &Node{
		Conn:          conn,
		Addr:          conn.RemoteAddr().String(),
		HeartbeatTime: currentTime,
		LoginTime:     currentTime,
		DataQueue:     make(chan []byte, 1),
		GroupSets:     set.New(set.ThreadSafe),
	}
	// 用户关系
	// userid  跟 node 绑定，并加锁
	rwLock.Lock()
	clientMap[uint(userId)] = node
	rwLock.Unlock()

	// 完成发送逻辑
	go sendProc(node)

	// 完成接收逻辑
	go recvProc(node)

	// 用户上线后在 redis 中存储在线信息
	SetUserOnlineInfo("online_"+id, []byte(node.Addr), time.Duration(viper.GetInt("timeout.redisOnlineTime"))*time.Hour)
	fmt.Println("---------------- user login --------------------------")
	fmt.Printf("%d online : %s\n", userId, time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("------------------------------------------------------")
}

// 设置用户在线状态
func SetUserOnlineInfo(key string, val []byte, timeTTL time.Duration) {
	ctx := context.Background()
	utils.RDP.Set(ctx, key, val, timeTTL)
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
		msg := Message{}
		err = json.Unmarshal(data, &msg)
		if err != nil {
			fmt.Println("json unmarshal failure: ", err)
			return
		}
		if msg.Type == 3 {
			//fmt.Println("[heartbeat]:", string(data))
			currentTime := uint64(time.Now().Unix())
			node.Heartbeat(currentTime)
		} else {
			// dispatch(data) // 不广播直接发送
			broadMsg(data) // 通过 upd 服务外送
			fmt.Println("[ws]<<<<", string(data))
		}
	}
}

var udpsendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendChan <- data
}

// 完成 ddp 数据发送协程
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 3001,
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
		Port: 3001,
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
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1: // 私发
		fmt.Println("[私聊]===", msg.TargetId)
		sendMsg(msg.TargetId, data, &msg)
	case 2: // 群发
		fmt.Println("[群聊]===", msg.TargetId)
		sendGroupMsg(msg.TargetId, data, &msg)
		//case 3: // 广播
		//case 4:
	}
}

// 群发消息
func sendGroupMsg(targetId uint, data []byte, msg *Message) {
	userIds := SearchUserByGroupId(targetId)
	for i := 0; i < len(userIds); i++ {
		if targetId != userIds[i] {
			sendMsg(userIds[i], data, msg)
		}
	}
}

func getMsgStoreKey(userId, targetId uint) (key string) {
	keyPattern := "msg_%d_%d"
	if targetId > userId {
		key = fmt.Sprintf(keyPattern, userId, targetId)
	} else {
		key = fmt.Sprintf(keyPattern, targetId, userId)
	}
	return
}

// 发送消息
func sendMsg(targetId uint, data []byte, msg *Message) {
	msg.CreateTime = uint64(time.Now().Unix())
	ctx := context.Background()
	// 判断消息的接收者是否在线
	online, err := utils.RDP.Get(ctx, fmt.Sprintf("online_%d", targetId)).Result()
	if err != nil {
		fmt.Println(err)
	}
	if online != "" {
		// 用户如果在线，则直接将消息先发给用户
		rwLock.RLock()
		node, ok := clientMap[targetId]
		rwLock.RUnlock()
		if ok {
			node.DataQueue <- data
		}
	} else {
		fmt.Println("用户不在线")
	}
  // 非私聊的消息不缓存
	if msg.Type > 1 {
		return
	}
	// 将消息持久化到 redis 中，方便后期上线后查看
	key := getMsgStoreKey(msg.UserId, targetId)
	res, err := utils.RDP.ZRevRange(ctx, key, 0, -1).Result() // MsgList
	if err != nil {
		fmt.Println("[Redis read err] === ", err)
	}
	score := float64(cap(res)) + 1
	result, e := utils.RDP.ZAdd(ctx, key, &redis.Z{Score: score, Member: data}).Result() // msg
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(result)
}

// 读取 Redis 中缓存的消息
func RedisMsg(userIdA, userIdB uint, start, end int64, isRev bool) []string {
	/**
	 * 方案 1. 获取消息然后逐条推送到客户端
	 */
	ctx := context.Background()
	key := getMsgStoreKey(userIdA, userIdB)

	/*
		rwLock.RLock()
		node, _ := clientMap[userIdA]
		rwLock.RUnlock()
		key := fmt.Sprintf(keyPattern, userIdB, userIdA)
		rels, err := utils.RDP.ZRevRange(ctx, key, 0, 10).Result() // 根据 score 来进行倒叙排列
		if err != nil {
			fmt.Println("查找数据异常")
		}
		for _, val := range rels {
			fmt.Println("sendMsg >> userID: ", userIdA, " msg: ", val)
			node.DataQueue <- []byte(val)
		}
	*/

	/**
	 * 方案2 : 一次查出所有数据后全部通过接口返回，交给前端自己显示
	 */

	var rels []string
	var err error
	// 是正序取出还是逆序取出
	if isRev {
		rels, err = utils.RDP.ZRange(ctx, key, start, end).Result()
	} else {
		rels, err = utils.RDP.ZRevRange(ctx, key, start, end).Result()
	}
	if err != nil {
		fmt.Println("查找数据异常")
	}
	return rels
}

// 定时任务清理超时连接
func CleanConnection(param interface{}) (result bool) {
	info := fmt.Sprintf("当前清理时间: %s, 参数 : %v", time.Now().Format("2006-01-02 15:04:05"), param)
	fmt.Println("[定时任务清理超时连接] === ", info)
  ctx := context.Background()
	result = true
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("cleanConnection err", r)
		}
	}()
	// 使用定时任务来清理已经过期的连接
	for i := range clientMap {
		node := clientMap[i]
		if node.IsHeartbeatTimeout() {
			fmt.Println("--------------------------------------------------")
			fmt.Println("超过最大时限未收到客户端心跳，自动下线, userId: ", i)
			fmt.Println("--------------------------------------------------")
			node.Conn.Close()
      delete(clientMap, i)
      utils.RDP.Del(ctx, fmt.Sprintf("online_%d", i))
		}
	}
	return
}

func (node *Node) Heartbeat(currentTime uint64) {
	node.HeartbeatTime = currentTime
}

func (node *Node) IsHeartbeatTimeout() (timeout bool) {
	currentTime := uint64(time.Now().Unix())
	if node.HeartbeatTime+viper.GetUint64("timeout.heartbeatMaxTime") <= currentTime {
		timeout = true
	}
	return
}

// 返回表名
func (table *Message) TableName() string {
	return "message"
}
