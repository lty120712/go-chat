package controller

import (
	"github.com/gin-gonic/gin"
	"go-chat/configs"
	interfacesservice "go-chat/internal/interfaces/service"
	"go-chat/internal/manager"
	request "go-chat/internal/model/request"
)

// MessageController 消息相关控制器
// @Tags Message
// @Description 消息相关控制器
type MessageController struct {
	BaseController
	messageService interfacesservice.MessageServiceInterface
}

var MessageControllerInstance *MessageController

func InitMessageController(messageService interfacesservice.MessageServiceInterface) {
	MessageControllerInstance = &MessageController{
		messageService: messageService,
	}
}

// SendString 发送消息接口
// @Summary 发送消息
// @Description 发送消息到指定队列，消息体为 JSON 格式
// @Tags Message
// @Accept json
// @Produce json
// @Param msg body map[string]interface{} true "消息内容"
// @Success 200 {object} model.Response "成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 500 {object} model.Response "发送消息失败"
// @Router /message/send/string [post]
func (con MessageController) SendString(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
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
// @Tags Message
// @Accept json
// @Produce json
// @Param msg body map[string]interface{} true "消息内容"
// @Success 200 {object} model.Response "成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 500 {object} model.Response "发送消息失败"
// @Router /message/send/json [post]
func (con MessageController) SendJson(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
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

// SendJson 已读消息接口
// @Summary 已读消息
// @Description 对该条消息已读
// @Tags Message
// @Accept json
// @Produce json
// @Param msg body model.ReadMessageReq true "消息内容"
// @Success 200 {object} model.Response "成功"
// @Failure 500 {object} model.Response "发送消息失败"
// @Router /message/read [post]
func (con MessageController) Read(c *gin.Context) {
	var req request.ReadMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, "参数错误")
		return
	}
	if err := con.messageService.ReadMessage(req.MessageId, req.UserId); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
}

// Query godoc
// @Summary 查询历史消息（分页，支持游标分页）
// @Description 根据目标ID和目标类型查询聊天消息历史，支持分页、时间范围等过滤
// @Tags Message
// @Accept application/json
// @Produce application/json
// @Param query body model.QueryMessagesRequest true "查询参数"
// @Success 200 {object} model.QueryMessagesResponse "查询成功，返回消息列表及分页信息"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器内部错误"
// @Router /message/query [post]
func (con MessageController) Query(c *gin.Context) {
	var req request.QueryMessagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}
	//userid来自中间件
	idStr, _ := c.Get("id")
	id := idStr.(uint)
	if data, err := con.messageService.QueryMessages(id, &req); err != nil {
		con.Error(c, err.Error())
		return
	} else {
		con.Success(c, data)
	}
}
