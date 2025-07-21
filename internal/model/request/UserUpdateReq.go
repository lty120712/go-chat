package model

// UserUpdateRequest 用户基础信息更新请求结构体
type UserUpdateRequest struct {
	ID       uint   `json:"id"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
}
