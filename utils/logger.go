package utils

import (
	"go.uber.org/zap"

)

// Logger global
var Logger *zap.Logger

// InitLogger inicializa el logger
func InitLogger() {
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer Logger.Sync() // Asegura que los logs se escriban
}

// Info registra información general
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Error registra errores
func Error(msg string, err error, fields ...zap.Field) {
	Logger.Error(msg, append(fields, zap.Error(err))...)
}

