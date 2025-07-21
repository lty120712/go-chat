package manager

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

type WebSocketManager struct {
	server      *http.Server
	upgrader    websocket.Upgrader
	connections sync.Map // 存储所有连接，键为用户ID，值为WebSocket连接
}

var WebSocketClient *WebSocketManager

// 初始化 WebSocket
func InitWebSocket() {
	WebSocketClient = &WebSocketManager{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// 这里可以加检查来源的逻辑
				return true
			},
		},
	}
	// 监听客户端连接
	http.HandleFunc("/ws", handleWebSocket)

	// 启动 WebSocket 服务
	go func() {
		WebSocketClient.server = &http.Server{Addr: ":80"}
		if err := WebSocketClient.server.ListenAndServe(); err != nil {
			logrus.Errorf("WebSocket 服务启动失败: %s", err)
		}
	}()
	logrus.Info("WebSocket 服务已启动")
}

// WebSocket 连接处理
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := WebSocketClient.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Errorf("WebSocket 连接失败: %s", err)
		return
	}
	defer conn.Close()

	// 获取用户ID
	id := r.URL.Query().Get("id")
	if id == "" {
		logrus.Warn("连接未提供有效的id")
		conn.WriteMessage(websocket.TextMessage, []byte("id is required"))
		return
	}
	// 存储连接
	WebSocketClient.connections.Store(id, conn)
	// 连接成功后的回调
	onOpen(conn, id)
	// 处理 WebSocket 消息
	for {
		// 读取消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			onError(conn, id, err)
			break
		}
		// 消息处理函数
		WebSocketClient.SendMessageToOne(id, "不要学我"+string(msg))
	}
	// 连接断开时调用 onClose 并删除连接
	onClose(conn, id)
}

func (ws *WebSocketManager) SendMessageToOne(id string, message interface{}) {
	conn, ok := WebSocketClient.connections.Load(id)
	if !ok {
		logrus.Warnf("没有找到用户 %s 的连接", id)
		return
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		logrus.Errorf("消息序列化失败: %s", err)
		return
	}
	err = conn.(*websocket.Conn).WriteMessage(websocket.TextMessage, messageBytes)
	if err != nil {
		logrus.Errorf("向用户 %v 发送消息失败: %s", id, err)
	}
}
func (ws *WebSocketManager) SendMessageToMultiple(ids []string, message interface{}) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		logrus.Errorf("消息序列化失败: %s", err)
		return
	}
	for _, id := range ids {
		ws.SendMessageToOne(id, messageBytes)
	}
}

func (ws *WebSocketManager) SendMessageToAll(message interface{}) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		logrus.Errorf("消息序列化失败: %s", err)
		return
	}
	WebSocketClient.connections.Range(func(key, value interface{}) bool {
		conn := value.(*websocket.Conn)
		err := conn.WriteMessage(websocket.TextMessage, messageBytes)
		if err != nil {
			logrus.Errorf("向用户 %v 发送消息失败: %s", key, err)
		}
		return true
	})
}

func onOpen(conn *websocket.Conn, id string) {
	logrus.Infof("WebSocket 客户端(%v)已连接: %s", conn.RemoteAddr(), id)
	//todo 更新心跳 时间
	WebSocketClient.SendMessageToOne(id, "连接成功")
}

func onClose(conn *websocket.Conn, id string) {
	logrus.Infof("WebSocket 客户端(%v)已断开: %s", id, conn.RemoteAddr())
	WebSocketClient.connections.Delete(id)
}

func onError(conn *websocket.Conn, id string, err error) {
	logrus.Infof("WebSocket 客户端(%v)%s发生错误:%v", id, conn.RemoteAddr(), err)
}

func (ws *WebSocketManager) Close() {
	if ws.server != nil {
		if err := ws.server.Close(); err != nil {
			logrus.Errorf("WebSocket 服务关闭失败: %s", err)
		}
	}
}
