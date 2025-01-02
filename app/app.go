package app

import (
	"net/http"

	"langtools/router"
	"langtools/utils"
)

func Run() error {
	// Cargar variables de entorno
	utils.LoadEnv()

	// Inicializar router
	r := router.NewRouter()

	// Print PID for debugging
	utils.Info("PID: " + utils.GetPID())

	// Configurar el servidor HTTP
	serverAddress := ":5000"
	utils.Info("Starting server on " + serverAddress)
	if err := http.ListenAndServe(serverAddress, r); err != nil {
		return err
	}

	return nil
}
