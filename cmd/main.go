package main

import (
	"log"

	"langtools/app"
	"langtools/utils"

)

func main() {
	// Inicializar logger
	utils.InitLogger()

	// Ejecutar aplicación
	if err := app.Run(); err != nil {
		utils.Error("Application failed", err)
		log.Fatalf("Application failed: %v", err)
	}
}
