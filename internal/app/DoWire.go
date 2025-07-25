package app

import (
	"github.com/sirupsen/logrus"
	controllers "go-chat/internal/controller"
	"go-chat/internal/repository"
	"go-chat/internal/service"
	wsHandler "go-chat/internal/ws/handler"
)

func doWire() {
	//repo
	repository.InitUserRepository()
	repository.InitMessageRepository()
	repository.InitGroupRepository()
	repository.InitGroupMemberRepository()
	repository.InitFriendRepository()
	repository.InitFriendRequestRepository()
	repository.InitFriendGroupRepository()

	//ws
	wsHandler.InitWebSocketHandler(nil, nil, nil)
	//service
	service.InitUserService(wsHandler.WebSocketHandlerInstance, repository.UserRepositoryInstance)
	service.InitMessageService(repository.MessageRepositoryInstance, repository.UserRepositoryInstance,
		repository.GroupMemberRepositoryInstance)
	service.InitGroupService(repository.GroupRepositoryInstance, repository.MessageRepositoryInstance,
		repository.UserRepositoryInstance, repository.GroupMemberRepositoryInstance)
	service.InitFriendService(repository.FriendRepositoryInstance, repository.FriendRequestRepositoryInstance,
		repository.FriendGroupRepositoryInstance, repository.UserRepositoryInstance, wsHandler.WebSocketHandlerInstance)
	//controller
	controllers.InitUserController(service.UserServiceInstance)
	controllers.InitMessageController(service.MessageServiceInstance)
	controllers.InitGroupController(service.GroupServiceInstance)
	controllers.InitFriendController(service.FriendServiceInstance)

	//延迟注入
	wsHandler.InitWebSocketHandler(service.UserServiceInstance, service.MessageServiceInstance, service.GroupServiceInstance)

	logrus.Info("=======================依赖注入完成=====================")
}
