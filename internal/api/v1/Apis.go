package apiv1

import (
	"github.com/gin-gonic/gin"
	"go-chat/configs"
	controllers "go-chat/internal/controller"
	"go-chat/internal/middleware"
)

func InitRouter(r *gin.Engine) {
	//配置路由中间件
	RegisterMiddlewares(r)
	//配置控制器的路由
	UserApi(r)
	MessageApi(r)
	GroupApi(r)
}

func UserApi(r *gin.Engine) {
	userApi := r.Group(configs.AppConfig.Api.Prefix + "/user")
	{
		userApi.GET("/ping", controllers.UserControllerInstance.Ping)
		userApi.POST("/register", controllers.UserControllerInstance.Register)
		userApi.POST("/login", controllers.UserControllerInstance.Login)
		userApi.GET("/logout", middleware.AuthMiddleware(), controllers.UserControllerInstance.Logout)
		userApi.GET("/info", controllers.UserControllerInstance.GetUserInfo)
		userApi.POST("/update", middleware.AuthMiddleware(), controllers.UserControllerInstance.Update)
	}
}

func MessageApi(r *gin.Engine) {
	messageApi := r.Group(configs.AppConfig.Api.Prefix+"/message", middleware.AuthMiddleware())
	{
		//测试rabbitmq
		messageApi.POST("/send/string", controllers.MessageControllerInstance.SendString)
		messageApi.POST("/send/json", controllers.MessageControllerInstance.SendJson)
		messageApi.POST("/read", controllers.MessageControllerInstance.Read)
		messageApi.POST("/query", controllers.MessageControllerInstance.Query)
	}
}

func GroupApi(r *gin.Engine) {
	groupApi := r.Group(configs.AppConfig.Api.Prefix+"/group", middleware.AuthMiddleware())
	{
		groupApi.POST("/create", controllers.GroupControllerInstance.Create)
		groupApi.POST("/update", controllers.GroupControllerInstance.Update)
		groupApi.GET("/join", controllers.GroupControllerInstance.Join)
		groupApi.GET("/quit", controllers.GroupControllerInstance.Quit)
		groupApi.POST("/search", controllers.GroupControllerInstance.Search)
		groupApi.GET("/member", controllers.GroupControllerInstance.Member)
	}
}
