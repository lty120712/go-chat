package service

import (
	"errors"
	"fmt"
	"go-chat/internal/db"
	models "go-chat/internal/model"
	model "go-chat/internal/model/request"
	response "go-chat/internal/model/response"
	"go-chat/internal/repository"
	"go-chat/internal/utils/jwtUtil"
	"go-chat/internal/utils/logUtil"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"sync"
)

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
func (u *UserService) Register(username, password, rePassword *string) (err error) {
	if *password != *rePassword {
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user = &models.User{Username: *username,
		Password:     string(hashedPassword),
		Nickname:     username,
		Status:       models.Enable,
		OnlineStatus: models.Offline,
	}
	err = userRepository.Save(user)
	if err != nil {
		logUtil.Errorf("保存用户失败: %v", err)
		return err
	}
	return nil
}

func (u *UserService) Login(username, password *string) (token string, err error) {
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
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(*password)); err != nil {
		return "", errors.New("用户名或密码错误")
	}
	//jwt 返回
	token, err = jwtUtil.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *UserService) UpdateUser(updateRequest *model.UserUpdateRequest) error {
	err := db.Mysql.Transaction(func(tx *gorm.DB) error {
		userRepository := repository.GetUserRepository()

		user, err := userRepository.GetById(updateRequest.ID, tx)
		if err != nil {
			return fmt.Errorf("查找用户失败: %v", err)
		}
		if user == nil {
			return errors.New("用户不存在")
		}
		updates := make(map[string]interface{})
		if updateRequest.Nickname != nil {
			updates["nickname"] = updateRequest.Nickname
		}
		if updateRequest.Desc != nil {
			updates["desc"] = updateRequest.Desc
		}
		if updateRequest.Avatar != nil {
			updates["avatar"] = updateRequest.Avatar
		}
		if updateRequest.Phone != nil {
			updates["phone"] = updateRequest.Phone
		}
		if updateRequest.Email != nil {
			updates["email"] = updateRequest.Email
		}
		if len(updates) == 0 {
			return errors.New("没有可更新的字段")
		}
		err = userRepository.UpdateFields(user.ID, updates, tx)
		if err != nil {
			return fmt.Errorf("更新用户失败: %v", err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) GetUserInfo(id uint) (response.UserVO, error) {
	userRepository := repository.GetUserRepository()
	userVo, err := userRepository.GetVoById(id)
	return userVo, err
}
