package main

import (
	"log"
	"net/http"

	"github.com/SandroK0/sync-tube-go/backend/api"
	"github.com/rs/cors"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/rooms", api.GetRoomsHandler)

	mux.HandleFunc("/ws", api.HandleConnections)
	go api.HandleMessages()

	handler := cors.Default().Handler(mux)
	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
