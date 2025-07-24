package ws

type Message struct {
	Type   string      `json:"type"`    // 事件类型
	SendId int64       `json:"send_id"` // 发送者ID
	Data   interface{} `json:"data"`    // 具体数据
}

// 事件类型
const (
	Chat      = "chat"       //聊天
	ChatAck   = "chat_ack"   // 聊天确认
	Online    = "online"     // 上线
	Offline   = "offline"    // 下线
	Recall    = "recall"     //  撤回
	IdRequest = "id_request" // 请求获取真实ID,引入mq之后采用
)
