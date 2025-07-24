package ws

import (
	"go-chat/internal/model"
	"go-chat/internal/utils/jsonUtil"
	"go-chat/internal/utils/logUtil"
	"net/http"
)

func MessageHandler(messageBytes []byte) {
	// 处理消息
	wsMessage := &Message{}
	err := jsonUtil.UnmarshalValue(messageBytes, wsMessage)
	if err != nil {
		logUtil.Errorf("消息反序列化失败: %s", err)
		return
	}
	switch wsMessage.Type {
	case Chat:
		ChatMessageHandler(wsMessage.SendId, wsMessage.Data)
	default:
		WebSocketClient.SendMessageToOne(wsMessage.SendId, &model.Response{
			Code:    http.StatusBadRequest,
			Message: "数据格式错误",
			Data:    nil,
		})
	}
}
