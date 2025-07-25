package interfaces

import (
	"go-chat/internal/model"
	response "go-chat/internal/model/response"
	"gorm.io/gorm"
)

type GroupMemberRepositoryInterface interface {
	SaveBatch(list []*model.GroupMember, tx *gorm.DB) error

	Save(member *model.GroupMember, tx ...*gorm.DB) error

	ExistsByGroupIdAndUserId(groupId uint, memberId uint, tx ...*gorm.DB) bool

	RejoinGroupIfDeleted(groupId uint, memberId uint, tx ...*gorm.DB) bool

	DeleteByGroupIdAndUserId(groupId uint, memberId uint, tx ...*gorm.DB) error

	GetMemberListByGroupId(groupId uint, tx ...*gorm.DB) ([]response.MemberVo, error)

	IsOwner(groupId uint, memberId uint, tx ...*gorm.DB) bool

	GetRelatedMemberByUserId(id uint) (memberList []response.MemberVo, err error)
}
