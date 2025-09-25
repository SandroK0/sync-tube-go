package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/SandroK0/sync-tube-go/backend/entities"
)

func GetRoomsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	WriteJson(Rooms, w)
}

func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	roomName, ok := data["roomName"].(string)
	if !ok || roomName == "" {
		log.Printf("Invalid or missing roomName")
		http.Error(w, "Bad request: roomName must be a string", http.StatusBadRequest)
		return
	}

	if _, exists := Rooms[roomName]; exists {
		WriteJson(`{"message": "Room already exists"}`, w)
	} else {
		Rooms[roomName] = &entities.Room{Name: roomName}
		WriteJson(`{"message": "Room created"}`, w)
	}
}
