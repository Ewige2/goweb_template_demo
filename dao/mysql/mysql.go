package mysql

import (
	"fmt"
	"web_app_template/settings"

	"go.uber.org/zap"

	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
	//匿名导入
	_ "github.com/go-sql-driver/mysql"
)

// 首字母 小写是私有的
var db *sqlx.DB

func InitDB(cfg *settings.MySqlConfig) (err error) {
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
	//	viper.GetString("mysql.user"),
	//	viper.GetString("mysql.password"),
	//	viper.GetString("mysql.host"),
	//	viper.GetInt("mysql.port"),
	//	viper.GetString("mysql.dbname"),
	//)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect  mysql   filed ", zap.Error(err))
		return
	}

	db.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns"))
	return
}

func Close() {
	_ = db.Close()

}
