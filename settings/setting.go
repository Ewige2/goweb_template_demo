package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 将配置集成到结构体上
// 全局变量保存程序的所有配置信息

// 项目配置
var Conf = new(multipleConfig)

type multipleConfig struct {
	*AppConfig   `mapstructure:"app"`
	*LogConfig   `mapstructure:"log"`
	*MySqlConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Mode    string `mapstructure:"mode"`
	version string `mapstructure:"version"`
	Port    string `mapstructure:"port"`
}

// 日志配置映射
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySqlConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"db_name"`
	Port     int    `mapstructure:"port"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// 使用viper 管理配置

func Init() (err error) {
	viper.SetConfigName("config")
	//viper.SetConfigType("yaml")  	//专用于从远处获取配置信息， 配合远程配置中心使用， 告诉   viper 配置使用什么  格式解析
	//指定配置文件所在目录， 可以配置多个目录
	viper.AddConfigPath("./conf/")
	//法二： 直接指定相关文件
	//viper.SetConfigFile("./conf/config.yaml")

	//读取配置信息
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig  failed  err:", err)
		return
	}

	//把读取到的配置信息反序列化到  结构体上
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed: %v", err)
	}

	//配置文件修改，  将会触发回调函数
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("setting 配置文件被修改了.....")
		//重新  序列化
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed: %v", err)
		}

	})

	return

}
