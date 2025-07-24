package model

import "gorm.io/gorm"

type GroupMember struct {
	gorm.Model
	GroupId   uint   `json:"group_id"`    //群ID
	MemberId  uint   `json:"member_id"`   //成员ID
	GNickName string `json:"g_nick_name"` //群昵称
	Role      Role   `json:"role"`        //成员角色（0=普通成员,1=群主，2=管理员）
	Status    Status `json:"status"`      //成员状态(0禁言1正常)
}

func (m *GroupMember) TableName() string {
	return "group_members"
}
