package ws

import (
	"go-chat/internal/model"
	"go-chat/internal/service"
	"go-chat/internal/utils"
	"go-chat/internal/utils/jsonUtil"
	"net/http"
)

func ChatMessageHandler(sendId int64, data interface{}) {
	bytes, err := jsonUtil.MarshalValue(data)
	if err != nil {
		WebSocketClient.SendMessageToOne(sendId, &model.Response{
			Code:    http.StatusBadRequest,
			Message: "数据格式错误",
			Data:    nil,
		})
		return
	}

	var message = &model.Message{}
	if err := jsonUtil.UnmarshalValue(bytes, message); err != nil {
		WebSocketClient.SendMessageToOne(sendId, &model.Response{
			Code:    http.StatusBadRequest,
			Message: "数据格式错误，无法反序列化成消息",
			Data:    nil,
		})
		return
	}
	message.InitFields()
	messageService := service.GetMessageService()
	vo, err := messageService.SendMessage(message)
	if err != nil {
		WebSocketClient.SendMessageToOne(sendId, &model.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	//消息发送后返回发送者:
	WebSocketClient.SendMessageToOne(sendId, &model.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data: &Message{
			SendId: sendId,
			Type:   ChatAck,
			Data:   vo,
		},
	})
	//消息发送后返回接收者,这里接收者私聊或群聊处理方式不同:
	if *message.TargetType == model.PrivateTarget {
		WebSocketClient.SendMessageToOne(*message.ReceiverId, &model.Response{
			Code:    http.StatusOK,
			Message: "success",
			Data: &Message{
				SendId: sendId,
				Type:   Chat,
				Data:   vo,
			},
		})
	} else if *message.TargetType == model.GroupTarget {
		groupService := service.GetGroupService()
		memberList, err := groupService.Member(uint(*message.GroupId))
		if err != nil {
			WebSocketClient.SendMessageToOne(sendId, &model.Response{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
		onlineUserIds := WebSocketClient.GetOnlineUserIds()
		groupOnlineUserIds := make([]int64, 0)
		for _, member := range memberList {
			//todo 需要加到下面条件中 && member.OnlineStatus == model.Online
			if utils.Contains(onlineUserIds, int64(member.UserId)) {
				groupOnlineUserIds = append(groupOnlineUserIds, int64(member.UserId))
			}
		}
		WebSocketClient.SendMessageToMultiple(groupOnlineUserIds, &model.Response{
			Code:    http.StatusOK,
			Message: "success",
			Data: &Message{
				SendId: sendId,
				Type:   Chat,
				Data:   vo,
			},
		})
	}
}
