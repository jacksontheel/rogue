package pubsub

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type    string `json:"type"`
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

type PubSub struct {
	mu             sync.RWMutex
	topicByClient  map[*websocket.Conn]string
	clientsByTopic map[string]map[*websocket.Conn]bool
}

func NewPubSub() *PubSub {
	return &PubSub{
		topicByClient:  make(map[*websocket.Conn]string),
		clientsByTopic: make(map[string]map[*websocket.Conn]bool),
	}
}

func (ps *PubSub) Subscribe(topic string, conn *websocket.Conn) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if previousTopic, exists := ps.topicByClient[conn]; exists {
		delete(ps.clientsByTopic[previousTopic], conn)
		if len(ps.clientsByTopic[previousTopic]) == 0 {
			delete(ps.clientsByTopic, previousTopic) // Clean up empty topics
		}
	}

	if ps.clientsByTopic[topic] == nil {
		ps.clientsByTopic[topic] = make(map[*websocket.Conn]bool)
	}
	ps.clientsByTopic[topic][conn] = true
	ps.topicByClient[conn] = topic
}

func (ps *PubSub) Unsubscribe(topic string, conn *websocket.Conn) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, ok := ps.clientsByTopic[topic]; ok {
		delete(ps.clientsByTopic[topic], conn)
		conn.Close()
		if len(ps.clientsByTopic[topic]) == 0 {
			delete(ps.clientsByTopic, topic)
		}
	}
}

func (ps *PubSub) Publish(topic, message string) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if clients, ok := ps.clientsByTopic[topic]; ok {
		for conn := range clients {
			err := conn.WriteJSON(Message{
				Type:    "message",
				Topic:   topic,
				Message: message,
			})
			if err != nil {
				conn.Close()
				delete(clients, conn)
			}
		}
	}
}
