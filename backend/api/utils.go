package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJson(data any, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	switch v := data.(type) {
	case string:
		w.Write([]byte(v))
	default:
		json.NewEncoder(w).Encode(v)
	}
}

func HandleEventError(err error, context string) {
	log.Printf("Error in %s: %v", context, err)
}
