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
		userApi.GET("/ping", controllers.UserController{}.Ping)
		userApi.POST("/register", controllers.UserController{}.Register)
		userApi.POST("/login", controllers.UserController{}.Login)
		userApi.GET("/info", controllers.UserController{}.GetUserInfo)
		userApi.POST("/update", middleware.AuthMiddleware(), controllers.UserController{}.Update)
	}
}

func MessageApi(r *gin.Engine) {
	messageApi := r.Group(configs.AppConfig.Api.Prefix+"/message", middleware.AuthMiddleware())
	{
		//测试rabbitmq
		messageApi.POST("/send/string", controllers.MessageController{}.SendString)
		messageApi.POST("/send/json", controllers.MessageController{}.SendJson)
		messageApi.POST("/read", controllers.MessageController{}.Read)
		messageApi.POST("/query", controllers.MessageController{}.Query)
	}
}

func GroupApi(r *gin.Engine) {
	groupApi := r.Group(configs.AppConfig.Api.Prefix+"/group", middleware.AuthMiddleware())
	{
		groupApi.POST("/create", controllers.GroupController{}.Create)
		groupApi.POST("/update", controllers.GroupController{}.Update)
		groupApi.GET("/join", controllers.GroupController{}.Join)
		groupApi.GET("/quit", controllers.GroupController{}.Quit)
		groupApi.POST("/search", controllers.GroupController{}.Search)
		groupApi.GET("/member", controllers.GroupController{}.Member)
	}
}
