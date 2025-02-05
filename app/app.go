package app

import (
	"net/http"
	"time"

	"langtools/router"
	"langtools/tools"
	"langtools/utils"

)

func Run() error {
	// Cargar variables de entorno
	utils.LoadEnv()

	// Inicializar router
	r := router.NewRouter()

	go func() {
		excelTask := func() {
			// Ruta del archivo

			// Generar Excel
			filePath, fileName, err := utils.GenerateExcel()
			if err != nil {
				utils.Error("Error generating Excel", err)
				return
			}

			// Subir a SharePoint
			err = tools.UploadFileToSharePoint(filePath, fileName)
			if err != nil {
				utils.Error("Error uploading to SharePoint", err)
			}
		}
		time.Sleep(10 * time.Second)
		// Configura el scheduler para las 12:30 am de todos los días
		excelTask()
		// handlers.StartScheduler("30 0 * * *", excelTask)
	}()

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
