package handler

import (
	"encoding/json"
	"fmt"
	"gin-web-framework/internal/service"
	"gin-web-framework/pkg/jwt"
	"gin-web-framework/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源，生产环境应该更严格
	},
}

// WebSocketHandler WebSocket处理器
type WebSocketHandler struct {
	realtimeService *service.RealtimeNotificationService
	logger          logger.LoggerInterface
}

// NewWebSocketHandler 创建WebSocket处理器
func NewWebSocketHandler(logger logger.LoggerInterface) *WebSocketHandler {
	return &WebSocketHandler{
		realtimeService: service.NewRealtimeNotificationService(logger),
		logger:          logger,
	}
}

// WebSocket连接
func (h *WebSocketHandler) WebSocket(c *gin.Context) {
	// 从查询参数获取token
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
		return
	}

	// 验证token并获取用户ID
	claims, err := jwt.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userID := claims.UserID

	// 升级HTTP连接为WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("Failed to upgrade connection: %v\n", err)
		return
	}
	defer conn.Close()

	// 添加客户端到实时通知服务
	h.realtimeService.AddClient(userID, conn)

	// 发送连接成功消息
	connectionMsg := map[string]interface{}{
		"type": "connection",
		"data": map[string]interface{}{
			"message": "WebSocket连接成功",
			"user_id": userID,
			"time":    time.Now().Format(time.RFC3339),
		},
	}

	if err := conn.WriteJSON(connectionMsg); err != nil {
		fmt.Printf("Failed to send connection message: %v\n", err)
		return
	}

	// 启动任务监控（如果还没有启动）
	go h.realtimeService.StartTaskMonitoring()

	// 保持连接活跃
	for {
		// 读取消息（保持连接）
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("WebSocket read error: %v\n", err)
			break
		}

		// 处理接收到的消息
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			fmt.Printf("Failed to unmarshal message: %v\n", err)
			continue
		}

		// 处理不同类型的消息
		switch msg["type"] {
		case "ping":
			// 响应ping消息
			pongMsg := map[string]interface{}{
				"type": "pong",
				"data": map[string]interface{}{
					"time": time.Now().Format(time.RFC3339),
				},
			}
			if err := conn.WriteJSON(pongMsg); err != nil {
				fmt.Printf("Failed to send pong: %v\n", err)
				return
			}
		case "get_notifications":
			// 获取最新通知
			notificationService := service.NewNotificationService(h.logger)
			notifications, err := notificationService.GetUserNotifications(userID, 10)
			if err != nil {
				fmt.Printf("Failed to get notifications: %v\n", err)
				continue
			}

			response := map[string]interface{}{
				"type": "notifications",
				"data": notifications,
			}

			if err := conn.WriteJSON(response); err != nil {
				fmt.Printf("Failed to send notifications: %v\n", err)
				return
			}
		}
	}

	// 移除客户端
	h.realtimeService.RemoveClient(userID)
}

// 获取WebSocket处理器实例
var websocketHandler *WebSocketHandler

// InitWebSocketHandler 初始化WebSocket处理器
func InitWebSocketHandler(logger logger.LoggerInterface) {
	websocketHandler = NewWebSocketHandler(logger)
}

// WebSocketHandlerFunc 返回WebSocket处理函数
func WebSocketHandlerFunc(c *gin.Context) {
	websocketHandler.WebSocket(c)
}
