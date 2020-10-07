package micro

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/avayayu/micro/net/http"
	ztime "github.com/avayayu/micro/time"
	"go.uber.org/zap"
)

type Micro interface {
	//定义微服务名
	ServiceName() string

	//自定义的服务配置
	MicroConfig() *http.ServerConfig
	//定义http路由
	HttpEngine(*http.ServerConfig) *http.Engine

	//获得服务的logger
	Logger() *zap.Logger
	// GRPCEngine() *Rpc.Engine
	// Dependency()

	//自定义路由
	SetupHandler()

	Close()
}

// func InJectCommonDep(micro Micro) error {
// 	micro.db
// }

func loadDefaultConfig(micro Micro) *http.ServerConfig {

	serverConfig := http.ServerConfig{
		Network:      "tcp",
		Addr:         "0.0.0.0:8080",
		Timeout:      ztime.Duration(time.Second * 60),
		ReadTimeout:  ztime.Duration(time.Second * 60),
		WriteTimeout: ztime.Duration(time.Second * 60),
	}
	return &serverConfig
}

func RunService(micro Micro) {

	var conf *http.ServerConfig

	conf = micro.MicroConfig()
	if conf == nil {
		conf = loadDefaultConfig(micro)
	}

	engine := micro.HttpEngine(conf)

	defer func() {
		micro.Close()
		engine.Shutdown(context.Background())
	}()

	//设置server的工作模式 开发环境或生产环境
	engine.SetServerMode(JudgeEnv())

	// engine.Use(http.CORS([]string{"http://127.0.0.1:8080", "http://localhost:8080", "http://192.168.100.104:8080", "http://192.168.100.180:8080"}))
	//server.engine.Use(http.JWTAuth())
	engine.Use(http.Recovery())
	engine.Use(http.TimeMiddleWare(micro.Logger()))
	engine.Use(http.RequestIDMiddleware(micro.Logger()))

	//加载自定义的中间件

	//挂载路由
	micro.SetupHandler()

	//开启服务
	engine.Run(conf.Addr)

	micro.Logger().Info("http server starts to receive http request")

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

}
