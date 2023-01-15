package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

const (
	PublishKey string = "websocket"
)

var (
	RDP *redis.Client
)

func InitRedis() {
	RDP = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConns"),
	})
	pong, err := RDP.Ping(RDP.Context()).Result()
	if err != nil {
		fmt.Println("cofing redis init err...")
	} else {
		fmt.Println("##################################################")
		fmt.Println("### ---------      mysql config    ------------###")
		fmt.Println("##################################################")
		fmt.Println("config redis: ", viper.Get("redis"))
		fmt.Println("redis init successful! redis heartbeat " + pong)
		fmt.Println("")
	}
}

// Publish 消息到 redis
func Publish(ctx context.Context, channel string, msg string) error {
	err := RDP.Publish(ctx, channel, msg).Err()
	return err
}

// Subscribe 订阅 redis 的消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := RDP.Subscribe(ctx, channel)
	fmt.Println("Subscribe ....", sub)
	msg, err := sub.ReceiveMessage(ctx)
	fmt.Println(">>::", msg.Payload)
	return msg.Payload, err
}
