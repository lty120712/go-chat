package model

type KickMemberRequest struct {
	GroupID  uint `json:"group_id" binding:"required"`  // 群组ID
	MemberID uint `json:"member_id" binding:"required"` // 被踢的成员ID
}
