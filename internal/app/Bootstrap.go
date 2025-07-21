package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-chat/configs"
	_ "go-chat/docs"
	apiv1 "go-chat/internal/api/v1"
	"go-chat/internal/db"
	"go-chat/internal/manager"
	"go-chat/internal/timer"
)

func Start() {
	//加载配置文件
	err := configs.LoadConfig()
	if err != nil {
		logrus.Error("Error loading config: %v", err)
		return
	}
	//配置数据库
	db.InitMysql()
	db.InitRedis()
	//配置logrus
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	//配置路由基本信息
	router := gin.Default()
	// 配置swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//配置路由
	apiv1.InitRouter(router)
	//配置rabbitmq
	manager.InitRabbitMQ()
	//配置WebSocket
	manager.InitWebSocket()
	//配置定时任务
	timer.InitTimer()
	//启动服务
	router.Run(fmt.Sprintf(":%d", configs.AppConfig.Server.Port))
}
