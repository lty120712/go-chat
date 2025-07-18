// Package model internal/models/response.go
package model

// Response 通用的 API 响应体结构体
type Response struct {
	Code    int         `json:"code"`    // 状态码
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 数据部分
}
