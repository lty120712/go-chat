package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username      string
	Password      string
	Phone         string
	Email         string
	Identity      string
	ClientIp      string
	ClientPort    string
	LoginTime     string
	HeartbeatTime uint64
	LogoutTIme    uint64
	IsLogout      int32
	DeviceInfo    string
}

func (u *User) TableName() string {
	return "users"
}
