package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"langtools/utils"

)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Permitir todas las conexiones, ajustar para producción
	},
}

// WebSocketHandler maneja las solicitudes de conexión WebSocket
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Actualizar la conexión a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	utils.Info("Client connected to WebSocket from " + r.RemoteAddr)


	// Manejar la conexión WebSocket
	HandleConnection(conn)
}
