package wsHandler

import (
	"go-chat/internal/model"
	"go-chat/internal/repository"
	wsClient "go-chat/internal/ws/client"
	wsMessage "go-chat/internal/ws/message"
	"net/http"
	"time"
)

// OnlineStatusNotice 在线状态通知,我的在线状态改变时通知与我相关的朋友或群组
func (ws *WebSocketHandler) OnlineStatusNotice(sendId int64, onlineStatusNotice model.OnlineStatusNotice) {

	//todo 3. 暂时先通知群组 之后还有好友
	memberList, _ := repository.GroupMemberRepositoryInstance.GetRelatedMemberByUserId(uint(sendId))

	var userIdList []int64
	seen := make(map[int64]bool)
	for _, member := range memberList {
		userId := int64(member.UserId)
		if !seen[userId] {
			seen[userId] = true
			userIdList = append(userIdList, userId)
		}
	}
	wsClient.WebSocketClient.SendMessageToMultiple(userIdList,
		&model.Response{
			Code:    http.StatusOK,
			Message: "success",
			Data: &wsMessage.Message{
				Type:   wsMessage.OnlineStatus,
				SendId: sendId,
				Data:   onlineStatusNotice,
				Time:   time.Now(),
			},
		})
}
