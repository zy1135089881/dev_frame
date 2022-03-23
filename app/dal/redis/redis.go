package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var (
	rdb *redis.Client
)

// 初始化连接
func InitRedis() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port")),
		Password: "",  // no password set
		DB:       0,   // use default DB
		PoolSize: 100, // 连接池大小
	})
	_, err = rdb.Ping().Result()
	return err
}

func Close() {
	_ = rdb.Close()
}
