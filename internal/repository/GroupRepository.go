package repository

import (
	"go-chat/internal/db"
	"go-chat/internal/model"
	"gorm.io/gorm"
	"sync"
)

type GroupRepository struct {
}

var (
	GroupRepositoryInstance *GroupRepository
	groupOnce               sync.Once
)

func GetGroupRepository() *GroupRepository {
	groupOnce.Do(func() {
		GroupRepositoryInstance = &GroupRepository{}
	})
	return GroupRepositoryInstance
}

func (g *GroupRepository) ExistsByCode(code string, tx ...*gorm.DB) bool {
	gormDB := db.GetGormDB(tx...)
	var exists bool
	rawSql := "SELECT EXISTS(SELECT 1 FROM `groups` WHERE code = ?)"
	err := gormDB.Raw(rawSql, code).Scan(&exists).Error

	if err != nil {
		return false
	}
	return exists
}

func (g *GroupRepository) Save(group *model.Group, tx ...*gorm.DB) error {
	gormDB := db.GetGormDB(tx...)
	return gormDB.Create(group).Error
}
