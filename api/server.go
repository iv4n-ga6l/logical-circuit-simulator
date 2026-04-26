package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	handlers "iv4n-ga6l/logical-circuit-simulator/handlers"
)

func StartServer() {
	r := mux.NewRouter()

	r.HandleFunc("/simulate", handlers.SimulationHandler).Methods(http.MethodPost)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}