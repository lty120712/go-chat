package apiv1

import (
	"github.com/gin-gonic/gin"
	"go-chat/configs"
	controllers "go-chat/internal/controller"
	"go-chat/internal/middleware"
)

func InitRouter(r *gin.Engine) {
	//配置路由中间件
	DoMiddlewares(r)
	//配置控制器的路由
	UserApi(r)
	MessageApi(r)
}

func UserApi(r *gin.Engine) {
	userApi := r.Group(configs.AppConfig.Api.Prefix + "/user")
	{
		userApi.GET("/ping", controllers.UserController{}.Ping)
		userApi.POST("/register", controllers.UserController{}.Register)
		userApi.POST("/login", controllers.UserController{}.Login)
		userApi.POST("/update_info", middleware.AuthMiddleware(), controllers.UserController{}.UpdateInfo)
	}
}

func MessageApi(r *gin.Engine) {
	messageApi := r.Group(configs.AppConfig.Api.Prefix + "/message")
	{
		//测试rabbitmq
		messageApi.POST("/send/string", controllers.MessageController{}.SendString)
		messageApi.POST("/send/json", controllers.MessageController{}.SendJson)
	}
}
