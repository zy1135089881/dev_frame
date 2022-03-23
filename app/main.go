package main

import (
	"app/config"
	"app/dal/mysql"
	"app/dal/redis"
	"app/logger"
	"app/routes"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	//1.加载配置
	if err := config.InitViper(); err != nil {
		fmt.Printf("InitViper failed,err:%v", err)
		return
	}

	//2.加载日志配置
	if err := logger.InitZap(); err != nil {
		fmt.Printf("InitZap failed,err:%v", err)
		return
	}
	defer zap.L().Sync()

	//3.数据库初始化
	if err := mysql.InitDB(); err != nil {
		fmt.Printf("InitDB failed,err:%v", err)
		return
	}
	defer mysql.Close()

	//4.redis初始化
	if err := redis.InitRedis(); err != nil {
		fmt.Printf("InitRedis failed,err:%v", err)
		return
	}
	defer redis.Close()

	//5.注册路由
	r := routes.Setup()

	//6.启动服务
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("Listing: %v\n", zap.Error(err))
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
		zap.L().Fatal("Server Shutdown: %v", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
