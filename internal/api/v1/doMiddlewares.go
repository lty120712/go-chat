package apiv1

import (
	"github.com/gin-gonic/gin"
	"go-chat/internal/middleware"
)

func DoMiddlewares(r *gin.Engine) {
	// 添加中间件
	// r.Use(middleware.Cors())
	// 用户限流
	r.Use(middleware.UserRateLimiterMiddleware())
}
