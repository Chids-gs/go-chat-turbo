package main

import (
	"log"
	"net/http"

	"go-chat/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Register routes
	handlers.SetupRoutes(r)

	// Serve static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}
