package service

import (
	"errors"
	"fmt"
	models "go-chat/internal/model"
	"go-chat/internal/repository"
	"go-chat/internal/utils"
	"go-chat/internal/utils/logUtil"
	"golang.org/x/crypto/bcrypt"
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

// Register 注册
func (u *UserService) Register(username, password, rePassword string) (err error) {
	if password != rePassword {
		return errors.New("密码不一致")
	}
	userRepository := repository.GetUserRepository()
	user, err := userRepository.GetByName(username)
	if err != nil {
		logUtil.Errorf("GetUserByName error: %v", err)
		return err
	}
	if user != nil {
		logUtil.Warnf("用户(%v)已存在", username)
		return errors.New(fmt.Sprintf("用户(%v)已存在", username))
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user = &models.User{Username: username, Password: string(hashedPassword), Nickname: username}
	err = userRepository.Save(user)
	if err != nil {
		logUtil.Errorf("保存用户失败: %v", err)
		return err
	}
	return nil
}

func (u *UserService) Login(username, password string) (token string, err error) {
	//根据username查询
	userRepository := repository.GetUserRepository()
	user, err := userRepository.GetByName(username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("用户名或密码错误")
	}
	//验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("用户名或密码错误")
	}
	//jwt 返回
	token, err = utils.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *UserService) UpdateUser(id uint, nickname string, avatar string, phone string, email string) error {
	return errors.New("未实现")
}
