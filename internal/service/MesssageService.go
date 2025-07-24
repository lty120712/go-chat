package service

import (
	"errors"
	"go-chat/internal/model"
	request "go-chat/internal/model/request"
	response "go-chat/internal/model/response"
	"go-chat/internal/repository"
	"go-chat/internal/utils"
	"sync"
)

type MessageService struct{}

var (
	messageServiceInstance *MessageService
	messageOnce            sync.Once
)

func GetMessageService() *MessageService {
	messageOnce.Do(func() {
		messageServiceInstance = &MessageService{}
	})
	return messageServiceInstance
}

// SendMessage 发送消息（支持私聊和群聊）
// msg 是已经构造好的 message 对象（建议外部构建 content 等）
func (s *MessageService) SendMessage(msg *model.Message) (*response.MessageVo, error) {
	if msg == nil {
		return nil, errors.New("消息不能为空")
	}
	if msg.SenderId == 0 {
		return nil, errors.New("发送者 Id 不能为空")
	}
	if *msg.TargetType == model.PrivateTarget && (msg.ReceiverId == nil) {
		return nil, errors.New("私聊消息必须有接收者 Id")
	}
	if *msg.TargetType == model.GroupTarget && (msg.GroupId == nil) {
		return nil, errors.New("群聊消息必须有群组 Id")
	}
	if msg.Type == nil {
		return nil, errors.New("消息类型不能为空")
	}
	if len(*msg.Content) == 0 {
		return nil, errors.New("消息内容不能为空")
	}
	messageRepo := repository.GetMessageRepository()
	if err := messageRepo.Save(msg); err != nil {
		return nil, err
	}
	vo, err := s.GetMessageById(msg.ID)
	if err != nil {
		return nil, err
	}
	return vo, nil
}

// GetMessageById  获取消息
func (s *MessageService) GetMessageById(id uint) (*response.MessageVo, error) {
	repo := repository.GetMessageRepository()
	message, err := repo.GetById(id)
	if err != nil {
		return nil, err
	}
	var messageVo = &response.MessageVo{}
	messageVo.GetFieldsFromMessage(message)
	//获取发送者信息
	userRepo := repository.GetUserRepository()
	sender, _ := userRepo.GetById(uint(message.SenderId))
	messageVo.SenderNickName = new(string)
	messageVo.SenderNickName = sender.Nickname
	messageVo.SenderAvatar = new(string)
	messageVo.SenderAvatar = sender.Avatar
	messageVo.SenderOnlineStatus = new(model.OnlineStatus)
	*messageVo.SenderOnlineStatus = sender.OnlineStatus
	return messageVo, nil
}

func (s *MessageService) ReadMessage(messageId uint, userId uint) error {
	//1.消息是否存在
	messageRepo := repository.GetMessageRepository()
	message, err := messageRepo.GetById(messageId)
	if err != nil {
		return err
	}
	if message == nil {
		return errors.New("消息不存在")
	}
	//2.将自己插入
	if message.ReaderIdList == nil {
		// 初始化为空切片
		message.ReaderIdList = &model.ReaderIdList{}
	}
	if !utils.Contains(*message.ReaderIdList, userId) {
		*message.ReaderIdList = append(*message.ReaderIdList, userId)
	}
	//3.更新数据库
	updateFields := map[string]interface{}{
		"reader_id_list": message.ReaderIdList,
	}
	err = messageRepo.UpdateFields(messageId, updateFields)
	if err != nil {
		return err
	}
	return nil
}

func (s *MessageService) QueryMessages(userId uint, req *request.QueryMessagesRequest) (*response.QueryMessagesResponse, error) {
	repo := repository.GetMessageRepository()
	userRepository := repository.GetUserRepository()
	memberRepository := repository.GetGroupMemberRepository()
	messages, err := repo.QueryHistoryMessages(userId, req)
	if err != nil {
		return nil, err
	}

	hasMore := false
	if len(messages) > req.Limit {
		hasMore = true
		messages = messages[:req.Limit]
	}
	senderIds := make([]uint, len(messages))
	for i, msg := range messages {
		senderIds[i] = uint(msg.SenderId)
	}
	var idToUserMap = make(map[uint]*model.User)
	userList, _ := userRepository.GetByIdList(senderIds)
	for _, user := range userList {
		user1 := user
		idToUserMap[user.ID] = &user1
	}
	if *req.TargetType == model.GroupTarget {
		memberList, _ := memberRepository.GetMemberListByGroupId(req.TargetId)
		if memberList != nil {
			for _, member := range memberList {
				if user, exists := idToUserMap[member.UserId]; exists {
					*user.Nickname = member.Nickname
				}
			}
		}
	}
	var list []*response.MessageVo
	for _, msg := range messages {
		sender := idToUserMap[uint(msg.SenderId)]
		list = append(list, &response.MessageVo{
			ID:                 msg.ID,
			CreatedAt:          msg.CreatedAt,
			UpdatedAt:          msg.UpdatedAt,
			SenderId:           msg.SenderId,
			ReceiverId:         msg.ReceiverId,
			GroupId:            msg.GroupId,
			ReplyId:            msg.ReplyId,
			ReaderIdList:       msg.ReaderIdList,
			TargetType:         msg.TargetType,
			Content:            msg.Content,
			Type:               msg.Type,
			Status:             msg.Status,
			ExtraData:          msg.ExtraData,
			IsRead:             msg.ReaderIdList != nil && utils.Contains(*msg.ReaderIdList, userId),
			SenderNickName:     sender.Nickname,
			SenderAvatar:       sender.Avatar,
			SenderOnlineStatus: &sender.OnlineStatus,
		})
	}

	var cursor int64 = 0
	if len(list) > 0 {
		cursor = int64(list[len(list)-1].ID)
	}

	return &response.QueryMessagesResponse{
		List:    list,
		Cursor:  cursor,
		HasMore: hasMore,
	}, nil
}
