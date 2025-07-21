package model

import (
	"encoding/json"
	"gorm.io/gorm"
)

// 聊天消息结构体
type Message struct {
	gorm.Model
	SenderId   string          `json:"sender_id"`   // 发送者Id
	ReceiverId string          `json:"receiver_id"` // 接收者Id(如果是群聊则为空)
	GroupId    string          `json:"group_id"`    // 群组Id（如果是私聊则为空）
	ReplyId    string          `json:"reply_id"`    // 回复消息Id
	TargetType TargetType      `json:"target_type"` // 消息目标类型（私聊、群聊）
	Content    []MessagePart   `json:"content"`     // 消息的内容部分（多个片段）
	Type       MessageType     `json:"type"`        // 消息类型（聊天消息、红包消息等）
	ExtraData  json.RawMessage `json:"extra_data"`  // 额外json数据,例如红包消息需要额外字段
}

// 消息类型 MessageType
type MessageType int

const (
	TextMessage  MessageType = iota // 聊天消息
	ImageMessage                    // 红包消息
	EmojiMessage                    // 表情消息
	LinkMessage                     // 链接消息
)

// 消息接收类型
type TargetType int

const (
	Private TargetType = iota // 私聊
	Group                     // 群聊
)

// 消息内容类型
type ContentType string

const (
	TextContent  ContentType = "text"  // 文本
	EmojiContent ContentType = "emoji" // 表情
	ImageContent ContentType = "image" // 图片
	LinkContent  ContentType = "link"  // 链接
)

type MessagePart struct {
	Type    ContentType `json:"type"`    // 内容类型（text, emoji, image, link）
	Content string      `json:"content"` // 内容（如文本、图片 URL、链接等）
}
