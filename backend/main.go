package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"example.com/rogue/db"
	"example.com/rogue/pubsub"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleWebSocket(ps *pubsub.PubSub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return
	}
	defer conn.Close()

	user := pubsub.User{
		UserId:     uuid.New().String(),
		Connection: conn,
	}

	userIDMsg := pubsub.UserIdMessage{
		Type:   "userId",
		UserId: user.UserId,
	}
	if err := conn.WriteJSON(userIDMsg); err != nil {
		log.Println("Error sending user ID:", err)
		conn.Close()
		return
	}

	for {
		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket Read Error:", err)
			ps.PublishPlayerExit(user)
			break
		}

		handleMessage(ps, user, messageBytes)
	}
}

func handleMessage(ps *pubsub.PubSub, user pubsub.User, rawMsg []byte) {
	var baseMsg pubsub.BaseMessage
	err := json.Unmarshal(rawMsg, &baseMsg)
	if err != nil {
		log.Println("Error unmarshaling base message:", err)
		return
	}

	switch baseMsg.Type {
	case "subscribe":
		var msg pubsub.SubscribeMessage
		err := json.Unmarshal(baseMsg.Data, &msg)
		if err != nil {
			log.Println("Error unmarshaling subscribe message:", err)
			return
		}
		ps.PublishPlayerExit(user)
		ps.Subscribe(msg.Topic, user)
		ps.PublishPlayerEntrance(user)
	case "position":
		var msg pubsub.PositionMessage
		err := json.Unmarshal(baseMsg.Data, &msg)
		if err != nil {
			log.Println("Error unmarshaling position message:", err)
			return
		}
		ps.PublishPosition(user, msg)
	default:
		log.Println("Unknown message type:", baseMsg.Type)
	}
}

func GetChunkHandler(db db.Database, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	xStr := r.URL.Query().Get("x")
	yStr := r.URL.Query().Get("y")

	x, err := strconv.Atoi(xStr)
	if err != nil {
		http.Error(w, "Invalid x parameter", http.StatusBadRequest)
		return
	}

	y, err := strconv.Atoi(yStr)
	if err != nil {
		http.Error(w, "Invalid y parameter", http.StatusBadRequest)
		return
	}

	chunk, err := db.GetChunk(x, y)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chunk)
}

func main() {
	host := "localhost"
	dbPort := 5432
	user := "username"
	name := "database"
	password := "password"

	database := db.GetDatabase(host, user, name, password, dbPort)
	db.GenerateWorld(database, 10)

	http.HandleFunc("/api/chunk", func(w http.ResponseWriter, r *http.Request) {
		GetChunkHandler(database, w, r)
	})

	ps := pubsub.NewPubSub()
	http.HandleFunc("/api/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(ps, w, r)
	})

	port := ":8080"
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
