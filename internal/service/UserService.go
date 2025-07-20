package service

import (
	"errors"
	"go-chat/internal/repository"
	"go-chat/internal/utils"
	"gorm.io/gorm"
	"sync"
)

type UserServiceInterface interface {
	// Register 注册
	Register(username, password, rePassword string) (ok bool, err error)
}

type UserService struct{}

var (
	userServiceInstance *UserService
	once                sync.Once
)

func GetUserService() *UserService {
	once.Do(func() {
		userServiceInstance = &UserService{}
	})
	return userServiceInstance
}

func (u *UserService) Register(username, password, rePassword string) (ok bool, err error) {
	//进行一系列操作后得到user然后 调用repository
	return true, nil
}

func (u *UserService) Login(username, password string) (token string, err error) {
	//根据username查询
	userRepository := repository.GetUserRepository()
	user, err := userRepository.GetByName(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("用户不存在")
		}
		return "", err
	}
	if utils.IsZero(user) {
		return "", errors.New("用户不存在")
	}
	//jwt 返回
	token, err = utils.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}
