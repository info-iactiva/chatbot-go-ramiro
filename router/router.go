package router

import (
	"net/http"

	"github.com/gorilla/mux"

	"langtools/ws"

)

// NewRouter crea un nuevo router con soporte para WebSocket
func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// Ruta para WebSocket
	r.HandleFunc("/ws", ws.WebSocketHandler).Methods("GET")

	// Otras rutas (e.g., para pruebas o salud del servidor)
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server is healthy!"))
	}).Methods("GET")

	return r
}
