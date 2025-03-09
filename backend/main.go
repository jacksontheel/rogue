package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"example.com/rogue/db"
	"example.com/rogue/pubsub"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleWebSocket(ps *pubsub.PubSub, w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		http.Error(w, "Missing userId query parameter", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return
	}
	defer conn.Close()

	var topic string

	for {
		var msg pubsub.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("WebSocket Read Error:", err)
			if topic != "" {
				ps.Unsubscribe(topic, conn)
			}
			break
		}

		switch msg.Type {
		case "subscribe":
			if msg.Topic == "" {
				log.Println("Subscription request missing topic")
				continue
			}
			topic = msg.Topic
			ps.Subscribe(topic, conn)
			log.Printf("Client subscribed to topic: %s", topic)
		case "publish":
			if msg.Topic == "" || msg.Message == "" {
				log.Println("Publish request missing topic or message")
				continue
			}
			ps.Publish(msg.Topic, msg.Message)
			log.Printf("Published message to topic: %s, Message: %s", msg.Topic, msg.Message)
		}
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
