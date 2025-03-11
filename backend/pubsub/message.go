package pubsub

import "encoding/json"

type BaseMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type SubscribeMessage struct {
	Topic string `json:"topic"`
}

type UserIdMessage struct {
	Type   string `json:"type"`
	UserId string `json:"userId"`
}

type PublishPositionMessage struct {
	Type            string          `json:"type"`
	UserId          string          `json:"userId"`
	PositionMessage PositionMessage `json:"data"`
}

type PublishPlayerExitMessage struct {
	Type   string `json:"type"`
	UserId string `json:"userId"`
}

type PositionMessage struct {
	X int `json:"x"`
	Y int `json:"y"`
}
