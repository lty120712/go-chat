package repository

import (
	"github.com/lty120712/gorm-pagination/pagination"
	"go-chat/internal/db"
	"go-chat/internal/model"
	request "go-chat/internal/model/request"
	"gorm.io/gorm"
	"sync"
)

type GroupRepository struct {
}

var (
	GroupRepositoryInstance *GroupRepository
	groupOnce               sync.Once
)

func InitGroupRepository() {
	groupOnce.Do(func() {
		GroupRepositoryInstance = &GroupRepository{}
	})

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

func (g *GroupRepository) UpdateRole(member *model.GroupMember, tx ...*gorm.DB) error {
	gormDB := db.GetGormDB(tx...)

	return gormDB.Model(&model.GroupMember{}).
		Where("group_id = ? AND member_id = ?", member.GroupId, member.MemberId).
		Update("role", member.Role).Error
}

func (g *GroupRepository) Page(req request.GroupSearchRequest, tx ...*gorm.DB) (*pagination.PageResult[model.Group], error) {
	gormDB := db.GetGormDB(tx...)
	query := gormDB.Model(&model.Group{})

	if req.Code != "" {
		query = query.Where("code = ?", req.Code)
	}
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	result := &pagination.PageResult[model.Group]{Records: []model.Group{}}

	// 调用分页函数并获取结果
	_, err := pagination.Paginate(query, req.Page, req.PageSize, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
