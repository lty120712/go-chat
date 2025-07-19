package controllers

import (
	"github.com/gin-gonic/gin"
	"go-chat/configs"
	"go-chat/internal/manager"
)

// MessageController 消息相关控制器
// @Tags Message
// @Description 消息相关控制器
type MessageController struct {
	BaseController
}

// SendString 发送消息接口
// @Summary 发送消息
// @Description 发送消息到指定队列，消息体为 JSON 格式
// @Tags message
// @Accept json
// @Produce json
// @Param msg body map[string]interface{} true "消息内容"
// @Success 200 {object} model.Response "成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 500 {object} model.Response "发送消息失败"
// @Router /message/send/string [post]
func (con MessageController) SendString(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		con.Error(c, "参数错误")
	}
	msg := req["msg"]
	// 发送消息
	chatMq := configs.AppConfig.Mq[0]
	err := manager.RabbitClient.SendMessage(chatMq.Exchange, chatMq.RoutingKey, msg)
	if err != nil {
		con.Error(c, "发送消息失败")
		return
	}
	con.Success(c, req)
}

// SendJson 发送消息接口
// @Summary 发送消息
// @Description 发送消息到指定队列，消息体为 JSON 格式
// @Tags message
// @Accept json
// @Produce json
// @Param msg body map[string]interface{} true "消息内容"
// @Success 200 {object} model.Response "成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 500 {object} model.Response "发送消息失败"
// @Router /message/send/json [post]
func (con MessageController) SendJson(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		con.Error(c, "参数错误")
	}
	// 发送消息
	chatMq := configs.AppConfig.Mq[0]
	err := manager.RabbitClient.SendMessage(chatMq.Exchange, chatMq.RoutingKey, req)
	if err != nil {
		con.Error(c, "发送消息失败")
		return
	}
	con.Success(c, req)
}
