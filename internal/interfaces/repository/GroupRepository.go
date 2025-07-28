package interfaces

import (
	"github.com/lty120712/gorm-pagination/pagination"
	"go-chat/internal/model"
	request "go-chat/internal/model/request"
	"gorm.io/gorm"
)

type GroupRepositoryInterface interface {
	ExistsByCode(code string, tx ...*gorm.DB) bool
	Save(group *model.Group, tx ...*gorm.DB) error

	UpdateRole(member *model.GroupMember, tx ...*gorm.DB) error

	Page(req request.GroupSearchRequest, tx ...*gorm.DB) (*pagination.PageResult[model.Group], error)
}
