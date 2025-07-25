package interfacesservice

import (
	request "go-chat/internal/model/request"
	response "go-chat/internal/model/response"
)

type GroupServiceInterface interface {
	// Create 创建群组
	Create(req *request.GroupCreateRequest) error
	Join(groupId uint, userId uint) error

	Quit(groupId uint, memberId uint) error

	Member(groupId uint) (memberList []response.MemberVo, err error)
}
