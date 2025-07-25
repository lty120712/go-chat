package interfaces

import (
	"go-chat/internal/model"
	"gorm.io/gorm"
)

type GroupRepositoryInterface interface {
	ExistsByCode(code string, tx ...*gorm.DB) bool
	Save(group *model.Group, tx ...*gorm.DB) error
}
