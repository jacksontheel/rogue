package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"example.com/rogue/db"
)

func GetChunkHandler(db db.Database) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

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
}

func main() {
	host := "localhost"
	dbPort := 5432
	user := "username"
	name := "database"
	password := "password"

	database := db.GetDatabase(host, user, name, password, dbPort)
	db.GenerateWorld(database, 10)

	http.HandleFunc("/api/chunk", GetChunkHandler(database))

	port := ":8080"
	log.Printf("Starting server on %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
