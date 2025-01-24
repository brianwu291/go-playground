package websockethandler

import (
	"encoding/json"
	"fmt"
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

func (ws *WebSocketHandler) sendClientsList(room *realtimechat.ChatRoom, conn *websocket.Conn) error {
	clients := room.GetClientInfoList()

	err := conn.WriteJSON(map[string]interface{}{
		"type":    "clients_list",
		"clients": clients,
	})
	if err != nil {
		return err
	}
	return nil
}

func (ws *WebSocketHandler) leaveRealTimeChatAndNotify(clientId string, roomName string) {
	if room, err := ws.chat.GetRoom(roomName); err == nil {
		if client, _ := room.GetClientById(clientId); client != nil {
			clientName := client.Name
			room.RemoveClient(clientId)
			leftMsg := fmt.Sprintf(realtimechat.ClientLeftTemplate, clientName, roomName)
			room.SendSystemMessage(leftMsg)
		}
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

	var currentClientId string
	var currentRoomName string
	defer func() {
		fmt.Printf("closed!\n")
		ws.leaveRealTimeChatAndNotify(currentClientId, currentRoomName)
		conn.Close()
	}()

	// handle WebSocket messages
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
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
				return
			}

			// add client to room
			client, err := room.AddClient(message.Username)
			if err != nil {
				log.Printf("Failed to add client: %v", err)
				return
			}
			currentClientId = client.Id
			currentRoomName = room.Name
			joinedMsg := fmt.Sprintf(realtimechat.ClientJoinedTemplate, client.Name, room.Name)
			room.SendSystemMessage(joinedMsg)

			// send join confirmation
			conn.WriteJSON(map[string]interface{}{
				"type":       "join_response",
				"clientId":   client.Id,
				"clientName": client.Name,
				"roomName":   message.RoomName,
			})

			ws.sendClientsList(room, conn)

			// start message listener
			go func() {
				for msg := range client.Messages {
					err := conn.WriteJSON(map[string]interface{}{
						"type":       "message",
						"id":         msg.Id,
						"content":    msg.Content,
						"authorId":   msg.AuthorId,
						"authorName": msg.AuthorName,
						"createdAt":  int(msg.CreatedAt.Unix()),
						"roomName":   msg.RoomName,
					})
					if err != nil {
						log.Printf("WebSocket write error: %v", err)
						return
					}
					ws.sendClientsList(room, conn)
				}
			}()

		case "message":
			currentRoom, err := ws.chat.GetRoom(message.RoomName)
			if err != nil {
				log.Printf("cannot find room: %+v", err)
				return
			}
			client, err := currentRoom.GetClientById(message.AuthorId)
			if err != nil {
				log.Printf("cannot find client: %+v", err)
				return
			}
			currentRoom.SendMessage(message.Content, client.Id)

		case "list_rooms":
			rooms := ws.chat.ListRooms()
			conn.WriteJSON(map[string]interface{}{
				"type":  "rooms_list",
				"rooms": rooms,
			})
		}
	}
}
