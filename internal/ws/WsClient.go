package ws

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

type WebSocketManager struct {
	Server      *http.Server
	Upgrader    websocket.Upgrader
	Connections sync.Map // 存储所有连接，键为用户ID，值为WebSocket连接
}

var WebSocketClient *WebSocketManager

func (ws *WebSocketManager) Close() {
	if ws.Server != nil {
		if err := ws.Server.Close(); err != nil {
			logrus.Errorf("WebSocket 服务关闭失败: %s", err)
		}
	}
}
