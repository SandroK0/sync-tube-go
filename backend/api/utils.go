package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJsonHTTP(data any, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func HandleEventError(err error, context string) {
	log.Printf("Error in %s: %v", context, err)
}
