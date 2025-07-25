package interfacesservice

import (
	request "go-chat/internal/model/request"
	response "go-chat/internal/model/response"
)

// UserServiceInterface 接口
type UserServiceInterface interface {
	Register(username, password, rePassword *string) (err error)
	Login(username, password *string) (token string, err error)
	UpdateUser(updateRequest *request.UserUpdateRequest) error
	GetUserInfo(id uint) (response.UserVO, error)
}
