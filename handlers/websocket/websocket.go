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

type wsMessage struct {
	Type     string `json:"type"`
	Content  string `json:"content"`
	AuthorId string `json:"authorId"`
	Username string `json:"username"`
	RoomName string `json:"roomName"`
}

func NewWebSocketHandler(chat *realtimechat.RealTimeChat) *WebSocketHandler {
	return &WebSocketHandler{
		chat: chat,
	}
}

func (ws *WebSocketHandler) HandleRealTimeChat(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: true,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	var currentRoom *realtimechat.ChatRoom
	var currentClient *realtimechat.Client

	// handle WebSocket messages
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			if currentClient != nil && currentRoom != nil {
				currentRoom.RemoveClient(currentClient.Id)
			}
			return
		}

		var message wsMessage
		if err := json.Unmarshal(p, &message); err != nil {
			continue
		}

		switch message.Type {
		case "join":
			// get or create a room
			room, err := ws.chat.GetOrCreateRoom(message.RoomName, 10)
			if err != nil {
				log.Printf("Failed to get/create room: %v", err)
				continue
			}

			// add client to room
			client, err := room.AddClient(message.Username)
			if err != nil {
				log.Printf("Failed to add client: %v", err)
				continue
			}

			currentRoom = room
			currentClient = client

			// send join confirmation
			conn.WriteJSON(map[string]interface{}{
				"type":     "join_response",
				"clientId": client.Id,
				"roomName": message.RoomName,
			})

			// start message listener
			go func() {
				for msg := range client.Messages {
					err := conn.WriteJSON(map[string]interface{}{
						"type":       "message",
						"id":         msg.Id,
						"content":    msg.Content,
						"authorId":   msg.AuthorId,
						"authorName": msg.AuthorName,
						"createdAt":  msg.CreatedAt,
						"roomName":   msg.RoomName,
					})
					if err != nil {
						log.Printf("WebSocket write error: %v", err)
						return
					}
				}
			}()

		case "message":
			if currentRoom != nil && currentClient != nil {
				currentRoom.SendMessage(message.Content, message.AuthorId)
			}

		case "list_rooms":
			rooms := ws.chat.ListRooms()
			conn.WriteJSON(map[string]interface{}{
				"type":  "rooms_list",
				"rooms": rooms,
			})
		}
	}
}
