package realtimechat

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type RealTimeChat struct {
	rooms      map[string]*ChatRoom
	roomsMutex sync.RWMutex
}

type ChatRoom struct {
	name        string
	maxClients  int
	clientGroup clientGroup
	messageChan chan Message
	stopChan    chan struct{}
}

type Client struct {
	Id       string
	Name     string
	Messages chan Message
}

type clientGroup struct {
	clients map[string]*Client
	mu      sync.RWMutex
}

type Message struct {
	Id         string
	AuthorName string
	Content    string
	CreatedAt  time.Time
	AuthorId   string
	RoomName   string
}

func NewRealTimeChat() *RealTimeChat {
	return &RealTimeChat{
		rooms: make(map[string]*ChatRoom),
	}
}

func (rt *RealTimeChat) GetOrCreateRoom(roomName string, maxClients int) (*ChatRoom, error) {
	rt.roomsMutex.Lock()
	defer rt.roomsMutex.Unlock()

	if room, exists := rt.rooms[roomName]; exists {
		return room, nil
	}

	room := &ChatRoom{
		name:       roomName,
		maxClients: maxClients,
		clientGroup: clientGroup{
			clients: make(map[string]*Client),
		},
		messageChan: make(chan Message, 100),
		stopChan:    make(chan struct{}),
	}

	go room.broadcastMessage()
	go func() {
		time.Sleep(90 * time.Minute)
		room.Stop()
		rt.removeRoom(roomName)
	}()

	rt.rooms[roomName] = room
	return room, nil
}

func (rt *RealTimeChat) removeRoom(roomName string) {
	rt.roomsMutex.Lock()
	defer rt.roomsMutex.Unlock()
	delete(rt.rooms, roomName)
}

func (room *ChatRoom) Stop() {
	room.clientGroup.mu.Lock()
	defer room.clientGroup.mu.Unlock()

	for _, client := range room.clientGroup.clients {
		select {
		case client.Messages <- Message{
			Id:         uuid.New().String(),
			Content:    "Chat room is closing",
			AuthorName: "System",
			CreatedAt:  time.Now(),
			RoomName:   room.name,
		}:
		default:
		}
	}

	time.Sleep(time.Second * 2)
	close(room.stopChan)
	for clientId, client := range room.clientGroup.clients {
		close(client.Messages)
		delete(room.clientGroup.clients, clientId)
	}
}

func (room *ChatRoom) AddClient(name string) (*Client, error) {
	room.clientGroup.mu.Lock()
	defer room.clientGroup.mu.Unlock()

	if len(room.clientGroup.clients) >= room.maxClients {
		return nil, errors.New("chat room is full")
	}

	clientId := uuid.New().String()
	client := &Client{
		Id:       clientId,
		Name:     name,
		Messages: make(chan Message, 100),
	}

	room.clientGroup.clients[clientId] = client
	return client, nil
}

func (room *ChatRoom) RemoveClient(id string) error {
	room.clientGroup.mu.Lock()
	defer room.clientGroup.mu.Unlock()

	client, exists := room.clientGroup.clients[id]
	if !exists {
		return errors.New("client not found")
	}

	close(client.Messages)
	delete(room.clientGroup.clients, id)
	return nil
}

func (room *ChatRoom) SendMessage(content string, authorId string) (*Message, error) {
	room.clientGroup.mu.RLock()
	client, exists := room.clientGroup.clients[authorId]
	room.clientGroup.mu.RUnlock()

	if !exists {
		return nil, errors.New("author not found")
	}

	msg := Message{
		Id:         uuid.New().String(),
		Content:    content,
		AuthorName: client.Name,
		CreatedAt:  time.Now(),
		AuthorId:   authorId,
		RoomName:   room.name,
	}

	select {
	case room.messageChan <- msg:
		return &msg, nil
	default:
		return nil, errors.New("message channel is full")
	}
}

func (room *ChatRoom) broadcastMessage() {
	for {
		select {
		case <-room.stopChan:
			fmt.Printf("the chat room %s has been stopped\n", room.name)
			return
		case msg := <-room.messageChan:
			room.clientGroup.mu.RLock()
			for _, client := range room.clientGroup.clients {
				select {
				case client.Messages <- msg:
				default:
					fmt.Printf("Failed to send message to client %s in room %s: buffer full\n",
						client.Id, room.name)
				}
			}
			room.clientGroup.mu.RUnlock()
		}
	}
}

func (room *ChatRoom) GetConnectedClientsCount() int {
	room.clientGroup.mu.RLock()
	defer room.clientGroup.mu.RUnlock()
	return len(room.clientGroup.clients)
}

func (rt *RealTimeChat) GetRoom(roomName string) (*ChatRoom, error) {
	rt.roomsMutex.RLock()
	defer rt.roomsMutex.RUnlock()

	room, exists := rt.rooms[roomName]
	if !exists {
		return nil, errors.New("room not found")
	}
	return room, nil
}

func (rt *RealTimeChat) ListRooms() []string {
	rt.roomsMutex.RLock()
	defer rt.roomsMutex.RUnlock()

	rooms := make([]string, 0, len(rt.rooms))
	for name := range rt.rooms {
		rooms = append(rooms, name)
	}
	return rooms
}
