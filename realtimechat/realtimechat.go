package realtimechat

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type RealTimeChat struct {
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
}

func NewRealTimeChat(maxClient int) *RealTimeChat {
	return &RealTimeChat{
		maxClients: maxClient,
		clientGroup: clientGroup{
			clients: make(map[string]*Client),
		},
		messageChan: make(chan Message, 100),
		stopChan:    make(chan struct{}),
	}
}

func (chat *RealTimeChat) Run() {
	go chat.broadcastMessage()
}

func (chat *RealTimeChat) Stop() {
	chat.clientGroup.mu.Lock()
	defer chat.clientGroup.mu.Unlock()
	close(chat.stopChan)
	for clientId, client := range chat.clientGroup.clients {
		close(client.Messages)
		delete(chat.clientGroup.clients, clientId)
	}
}

func (chat *RealTimeChat) AddClient(name string) (*Client, error) {
	chat.clientGroup.mu.Lock()
	defer chat.clientGroup.mu.Unlock()

	if len(chat.clientGroup.clients) >= chat.maxClients {
		return nil, errors.New("chat room is full")
	}

	clientId := uuid.New().String()
	client := &Client{
		Id:       clientId,
		Name:     name,
		Messages: make(chan Message, 100),
	}
	go client.receiveMessage()

	chat.clientGroup.clients[clientId] = client
	return client, nil
}

func (chat *RealTimeChat) RemoveClient(id string) error {
	chat.clientGroup.mu.Lock()
	defer chat.clientGroup.mu.Unlock()

	client, exists := chat.clientGroup.clients[id]
	if !exists {
		return errors.New("client not found")
	}

	close(client.Messages)
	delete(chat.clientGroup.clients, id)
	return nil
}

func (client *Client) receiveMessage() {
	fmt.Printf("start for client %+v\n", client.Name)
	for msg := range client.Messages {
		fmt.Printf("%s: %s\n", msg.AuthorName, msg.Content)
	}
	fmt.Printf("end for client %+v\n", client.Name)
}

func (chat *RealTimeChat) SendMessage(content string, authorId string) (*Message, error) {
	chat.clientGroup.mu.RLock()
	client, exists := chat.clientGroup.clients[authorId]
	chat.clientGroup.mu.RUnlock()

	if !exists {
		return nil, errors.New("author not found")
	}

	msg := Message{
		Id:         uuid.New().String(),
		Content:    content,
		AuthorName: client.Name,
		CreatedAt:  time.Now(),
		AuthorId:   authorId,
	}

	select {
	case chat.messageChan <- msg:
		return &msg, nil
	default:
		return nil, errors.New("message channel is full")
	}
}

func (chat *RealTimeChat) broadcastMessage() {
	for {
		select {
		case <-chat.stopChan:
			fmt.Printf("the chat has been stopped\n")
			return
		case msg := <-chat.messageChan:
			chat.clientGroup.mu.RLock()
			for _, client := range chat.clientGroup.clients {
				select {
				case client.Messages <- msg:
				default:
					fmt.Printf("Failed to send message to client %s: buffer full\n", client.Id)
				}
			}
			chat.clientGroup.mu.RUnlock()
		}
	}
}

func (chat *RealTimeChat) GetConnectedClientsCount() int {
	chat.clientGroup.mu.RLock()
	defer chat.clientGroup.mu.RUnlock()
	return len(chat.clientGroup.clients)
}

func (chat *RealTimeChat) GetClient(id string) (*Client, error) {
	chat.clientGroup.mu.RLock()
	defer chat.clientGroup.mu.RUnlock()

	client, exists := chat.clientGroup.clients[id]
	if !exists {
		return nil, errors.New("client not found")
	}
	return client, nil
}
