package redis

import (
	"fmt"
	"web_app_template/settings"

	"github.com/go-redis/redis"
)

//申明全局变量

var rdb *redis.Client

// 初始化连接
func InitRedis(cfg *settings.RedisConfig) (err error) {
	//rdb := redis.NewClient(&redis.Options{
	//	Addr:     fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port")),
	//	Password: viper.GetString("password"),    // 密码
	//	DB:       viper.GetInt("redis.dbname"),   // 数据库
	//	PoolSize: viper.GetInt("redis.poolsize"), // 连接池大小
	//})
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		return err
	}

	return nil

}

func Close() {
	_ = rdb.Close()

}
