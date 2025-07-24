package model

import (
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	Code    string `json:"code"`        //群号
	Name    string `json:"name"`        //群组名称
	Desc    string `json:"description"` //群组简介
	OwnerId uint   `json:"owner_id"`    //群主ID
	MaxNum  int    `json:"max_num"`     //群组最大人数
	Status  Status `json:"status" `     //群组状态（1=正常，0=关闭）
}

func (m *Group) TableName() string {
	return "groups"
}
