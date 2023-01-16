package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web_app_template/dao/mysql"
	"web_app_template/dao/redis"
	"web_app_template/logger"
	"web_app_template/routes"
	"web_app_template/settings"

	"go.uber.org/zap"
)

// go web  通用开发脚手架

func main() {

	//1.加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init  settings failed: %v", err)

		return
	}
	//2.初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("logger.Init() failed: %v", err)
		return
	}
	//同步日志，将缓存区中的日志  追加到 文件中
	defer zap.L().Sync()
	zap.L().Debug("log  init success....")

	//3. 初始化Mysql连接
	if err := mysql.InitDB(settings.Conf.MySqlConfig); err != nil {
		fmt.Printf("init mysql.InitDB() failed :%v", err)
		return
	}

	defer mysql.Close()
	//4. 初始化redis 连接
	if err := redis.InitRedis(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init  redis failed %v", err)
		return
	}
	defer redis.Close()
	//5. 注册路由
	router := routes.RegisterRoute()
	//6. 启动服务（优雅关机）
	// Go 1.8版本之后， http.Server 内置的 Shutdown() 方法就支持优雅地关机，具体示例如下：

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: %s\n", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
