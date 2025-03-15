package handlers

import (
	"log"
	"net/http"

	"go-chat/internal/templates"

	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/ws", wsHandler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.IndexTemplate.Execute(w, nil); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
