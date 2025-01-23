package websockethandler

import (
	"encoding/json"
	"log"
	"net/http"

	realtimechat "github.com/brianwu291/go-playground/realtimechat"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	chat *realtimechat.RealTimeChat
}

func NewWebSocketHandler(chat *realtimechat.RealTimeChat) *WebSocketHandler {
	return &WebSocketHandler{
		chat: chat,
	}
}

func (ws *WebSocketHandler) HandleRealTimeChat(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // For development only
		},
		EnableCompression: true,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Handle WebSocket messages
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var message struct {
			Type     string `json:"type"`
			Content  string `json:"content"`
			AuthorId string `json:"authorId"`
			Username string `json:"username"`
		}

		if err := json.Unmarshal(p, &message); err != nil {
			continue
		}

		switch message.Type {
		case "join":
			client, err := ws.chat.AddClient(message.Username)
			if err == nil {
				conn.WriteJSON(map[string]interface{}{
					"type":     "join_response",
					"clientId": client.Id,
				})
				go func() {
					for msg := range client.Messages {
						err := conn.WriteJSON(map[string]interface{}{
							"type":       "message",
							"id":         msg.Id,
							"content":    msg.Content,
							"authorId":   msg.AuthorId,
							"authorName": msg.AuthorName,
							"createdAt":  msg.CreatedAt,
						})
						if err != nil {
							log.Printf("WebSocket write error: %v", err)
							return
						}
					}
				}()
			}
		case "message":
			ws.chat.SendMessage(message.Content, message.AuthorId)
		}
	}
}
