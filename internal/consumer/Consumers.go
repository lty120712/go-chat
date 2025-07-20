package consumer

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

// 处理队列 "string" 的消息逻辑
func HandleStringConsumer(msg []byte) {
	message := string(msg)
	logrus.Printf("处理 message 队列的消息: %s", message)
}

// 处理队列 "json" 的消息逻辑
func HandleJsonConsumer(msg []byte) {
	var jsonMessage map[string]interface{}
	err := json.Unmarshal(msg, &jsonMessage)
	if err != nil {
		logrus.Printf("消息反序列化失败: %s", err)
		return
	}
	logrus.Printf("处理 notifications 队列的通知: %+v", jsonMessage)
}
