package model

import "time"

type GroupMuteRequest struct {
	GroupId uint      `json:"group_id" binding:"required"`
	MuteEnd time.Time `json:"mute_end" binding:"required"`
}
