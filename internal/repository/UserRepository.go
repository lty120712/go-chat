package repository

import (
	"errors"
	"go-chat/internal/db"
	models "go-chat/internal/model"
	"gorm.io/gorm"
	"sync"
)

type UserRepository struct {
}

var (
	UserRepositoryInstance *UserRepository
	once                   sync.Once
)

func GetUserRepository() *UserRepository {
	once.Do(func() {
		UserRepositoryInstance = &UserRepository{}
	})
	return UserRepositoryInstance
}

// GetById 只执行sql相关操作
func (r *UserRepository) GetById(id int) (user *models.User, err error) {
	user = &models.User{}
	err = db.Mysql.Where("id = ?", id).First(user).Error
	return
}

// GetByName 根据用户名查询用户
func (r *UserRepository) GetByName(username string) (*models.User, error) {
	user := &models.User{}
	err := db.Mysql.Where("username = ?", username).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return user, err
}

// Save 保存用户
func (r *UserRepository) Save(user *models.User) (err error) {
	err = db.Mysql.Create(user).Error
	return
}
