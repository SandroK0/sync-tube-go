package api

import (
	"net/http"
)

func GetRoomsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	WriteJsonHTTP(Rooms, w)
}
