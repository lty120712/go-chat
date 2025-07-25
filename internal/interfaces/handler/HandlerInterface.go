// Package interfaces
package interfaces

import "go-chat/internal/model"

// WsHandlerInterface  接口
type WsHandlerInterface interface {
	OnlineStatusNotice(sendId int64, data model.OnlineStatusNotice)
}
