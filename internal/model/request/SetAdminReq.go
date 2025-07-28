package model

type SetAdminRequest struct {
	GroupID  uint `json:"group_id" binding:"required"`
	MemberID uint `json:"member_id" binding:"required"`
}
