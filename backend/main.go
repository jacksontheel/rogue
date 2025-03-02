package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"example.com/rogue/world"
)

func GetChunkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	xStr := r.URL.Query().Get("x")
	yStr := r.URL.Query().Get("y")

	_, err := strconv.Atoi(xStr)
	if err != nil {
		http.Error(w, "Invalid x parameter", http.StatusBadRequest)
		return
	}

	_, err = strconv.Atoi(yStr)
	if err != nil {
		http.Error(w, "Invalid y parameter", http.StatusBadRequest)
		return
	}

	chunk := world.GenerateChunk([]int{11, 23, 34}, []int{1, 23, 34})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chunk)
}

func main() {
	http.HandleFunc("/api/chunk", GetChunkHandler)

	port := ":8080"
	log.Printf("Starting server on %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
