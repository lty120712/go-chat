package wsHandler

import (
	"go-chat/internal/model"
	"go-chat/internal/service"
	"go-chat/internal/utils"
	"go-chat/internal/utils/jsonUtil"
	wsClient "go-chat/internal/ws/client"
	wsMessage "go-chat/internal/ws/message"
	"net/http"
	"time"
)

func ChatHandler(sendId int64, data interface{}) {
	bytes, err := jsonUtil.MarshalValue(data)
	if err != nil {
		wsClient.WebSocketClient.SendMessageToOne(sendId, &model.Response{
			Code:    http.StatusBadRequest,
			Message: "数据格式错误",
			Data:    nil,
		})
		return
	}

	var message = &model.Message{}
	if err := jsonUtil.UnmarshalValue(bytes, message); err != nil {
		wsClient.WebSocketClient.SendMessageToOne(sendId, &model.Response{
			Code:    http.StatusBadRequest,
			Message: "数据格式错误，无法反序列化成消息",
			Data:    nil,
		})
		return
	}
	message.InitFields()
	vo, err := service.MessageServiceInstance.SendMessage(message)
	if err != nil {
		wsClient.WebSocketClient.SendMessageToOne(sendId, &model.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	//消息发送后返回发送者:
	wsClient.WebSocketClient.SendMessageToOne(sendId, &model.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data: &wsMessage.Message{
			SendId: sendId,
			Type:   wsMessage.ChatAck,
			Data:   vo,
			Time:   time.Now(),
		},
	})
	//消息发送后返回接收者,这里接收者私聊或群聊处理方式不同:
	if *message.TargetType == model.PrivateTarget {
		wsClient.WebSocketClient.SendMessageToOne(*message.ReceiverId, &model.Response{
			Code:    http.StatusOK,
			Message: "success",
			Data: &wsMessage.Message{
				SendId: sendId,
				Type:   wsMessage.Chat,
				Data:   vo,
				Time:   time.Now(),
			},
		})
	} else if *message.TargetType == model.GroupTarget {
		memberList, err := service.GroupServiceInstance.Member(uint(*message.GroupId))
		if err != nil {
			wsClient.WebSocketClient.SendMessageToOne(sendId, &model.Response{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
		onlineUserIds := wsClient.WebSocketClient.GetOnlineUserIds()
		groupOnlineUserIds := make([]int64, 0)
		for _, member := range memberList {
			if utils.Contains(onlineUserIds, int64(member.UserId)) && member.OnlineStatus == model.Online {
				groupOnlineUserIds = append(groupOnlineUserIds, int64(member.UserId))
			}
		}
		wsClient.WebSocketClient.SendMessageToMultiple(groupOnlineUserIds, &model.Response{
			Code:    http.StatusOK,
			Message: "success",
			Data: &wsMessage.Message{
				SendId: sendId,
				Type:   wsMessage.Chat,
				Data:   vo,
				Time:   time.Now(),
			},
		})
	}
}
