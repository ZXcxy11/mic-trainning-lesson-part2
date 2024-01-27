package internal

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// 声明redis客户端对象
var RedisClient *redis.Client

func InitRedis() {
	//	从viper获取的数据中读取redis信息
	host := AppConf.RedisConfig.Host
	port := AppConf.RedisConfig.Port
	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Println(addr)
	//	获取redis客户端
	RedisClient = redis.NewClient(&redis.Options{
		Addr: addr,
	})
	//	Ping检测连接
	ping := RedisClient.Ping(context.Background())
	fmt.Println(ping.String())
	fmt.Println("Redis初始化完成...")
}
