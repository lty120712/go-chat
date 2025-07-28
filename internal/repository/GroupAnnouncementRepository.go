package repository

import (
	"go-chat/internal/db"
	"go-chat/internal/model"
	"gorm.io/gorm"
	"sync"
)

type GroupAnnouncementRepository struct {
}

var (
	GroupAnnouncementRepositoryInstance *GroupAnnouncementRepository
	groupAnnouncementOnce               sync.Once
)

// InitGroupAnnouncementRepository 初始化群公告仓库实例
func InitGroupAnnouncementRepository() {
	groupAnnouncementOnce.Do(func() {
		GroupAnnouncementRepositoryInstance = &GroupAnnouncementRepository{}
	})
}

// Create 创建群组公告
// @param announcement 群组公告模型
// @param tx 事务对象，可选
// @return error 创建失败时返回错误
func (r *GroupAnnouncementRepository) Create(announcement *model.GroupAnnouncement, tx ...*gorm.DB) error {
	gormDB := db.GetGormDB(tx...)
	return gormDB.Create(announcement).Error
}

// Update 更新群组公告
// @param announcement 群组公告模型
// @param tx 事务对象，可选
// @return error 更新失败时返回错误
func (r *GroupAnnouncementRepository) Update(announcement *model.GroupAnnouncement, tx ...*gorm.DB) error {
	gormDB := db.GetGormDB(tx...)
	return gormDB.Save(announcement).Error
}

// Delete 删除群组公告
// @param announcementId 公告ID
// @param tx 事务对象，可选
// @return error 删除失败时返回错误
func (r *GroupAnnouncementRepository) Delete(announcementId uint, tx ...*gorm.DB) error {
	gormDB := db.GetGormDB(tx...)
	return gormDB.Delete(&model.GroupAnnouncement{}, announcementId).Error
}

// GetByID 根据公告ID查询群组公告
// @param announcementId 公告ID
// @param tx 事务对象，可选
// @return *model.GroupAnnouncement 公告对象，错误信息
func (r *GroupAnnouncementRepository) GetByID(announcementId uint, tx ...*gorm.DB) (*model.GroupAnnouncement, error) {
	gormDB := db.GetGormDB(tx...)
	var announcement model.GroupAnnouncement
	err := gormDB.Where("id = ?", announcementId).First(&announcement).Error
	if err != nil {
		return nil, err
	}
	return &announcement, nil
}

// GetLatestByGroupID 获取群组的最新公告
// @param groupId 群组ID
// @param tx 事务对象，可选
// @return *model.GroupAnnouncement 最新公告，错误信息
func (r *GroupAnnouncementRepository) GetLatestByGroupID(groupId uint, tx ...*gorm.DB) (*model.GroupAnnouncement, error) {
	gormDB := db.GetGormDB(tx...)
	var announcement model.GroupAnnouncement
	err := gormDB.Where("group_id = ?", groupId).
		Order("created_at DESC"). // 按创建时间倒序排列，获取最新公告
		First(&announcement).Error
	if err != nil {
		return nil, err
	}
	return &announcement, nil
}

// GetListByGroupID 获取群组的所有公告
// @param groupId 群组ID
// @param tx 事务对象，可选
// @return []model.GroupAnnouncement 群组所有公告，错误信息
func (r *GroupAnnouncementRepository) GetListByGroupID(groupId uint, tx ...*gorm.DB) ([]model.GroupAnnouncement, error) {
	gormDB := db.GetGormDB(tx...)
	var announcements []model.GroupAnnouncement
	err := gormDB.Where("group_id = ?", groupId).
		Order("created_at DESC"). // 按创建时间倒序排列
		Find(&announcements).Error
	if err != nil {
		return nil, err
	}
	return announcements, nil
}
