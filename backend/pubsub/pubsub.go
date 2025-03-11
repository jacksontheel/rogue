package pubsub

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type    string `json:"type"`
	Topic   string `json:"topic"`
	UserId  string `json:"userId"`
	Message string `json:"message"`
}

type User struct {
	UserId     string
	Connection *websocket.Conn
}

type PubSub struct {
	mu           sync.RWMutex
	topicByuser  map[User]string
	usersByTopic map[string]map[User]bool
}

func NewPubSub() *PubSub {
	return &PubSub{
		topicByuser:  make(map[User]string),
		usersByTopic: make(map[string]map[User]bool),
	}
}

func (ps *PubSub) Subscribe(topic string, user User) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if previousTopic, exists := ps.topicByuser[user]; exists {
		delete(ps.usersByTopic[previousTopic], user)
		if len(ps.usersByTopic[previousTopic]) == 0 {
			delete(ps.usersByTopic, previousTopic) // Clean up empty topics
		}
	}

	if ps.usersByTopic[topic] == nil {
		ps.usersByTopic[topic] = make(map[User]bool)
	}
	ps.usersByTopic[topic][user] = true
	ps.topicByuser[user] = topic
}

func (ps *PubSub) PublishPlayerExit(user User) {
	ps.Publish(user, PublishPlayerExitMessage{
		Type:   "playerExit",
		UserId: user.UserId,
	})
}

func (ps *PubSub) PublishPosition(user User, message PositionMessage) {
	ps.Publish(user, PublishPositionMessage{
		Type:            "position",
		UserId:          user.UserId,
		PositionMessage: message,
	})
}

func (ps *PubSub) Publish(user User, message interface{}) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	topic := ps.topicByuser[user]

	if subscribers, ok := ps.usersByTopic[topic]; ok {
		for subscriber := range subscribers {
			conn := subscriber.Connection
			err := conn.WriteJSON(message)
			if err != nil {
				conn.Close()
				delete(subscribers, subscriber)
			}
		}
	}
}
