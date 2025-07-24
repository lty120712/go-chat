package model

import "go-chat/internal/model"

type MemberVo struct {
	UserId       uint               `json:"user_id"`
	Nickname     string             `json:"nickname"`
	Role         model.Role         `json:"role"`
	Avatar       *string            `json:"avatar"`
	OnlineStatus model.OnlineStatus `json:"online_status"`
}
