package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"langtools/utils"

)

// GetMongoDBStore configura y retorna una instancia de cliente MongoDB.
func GetMongoDBStore() (*mongo.Client, error) {
	// Cargar la URI de conexión desde las variables de entorno
	connectionString := utils.GetEnv("MONGO_CONNECTION_STRING", "")
	if connectionString == "" {
		utils.Error("MONGO_CONNECTION_STRING no configurada en las variables de entorno", nil)
		return nil, nil
	}

	// Configurar opciones del cliente MongoDB
	clientOptions := options.Client().ApplyURI(connectionString)

	// Establecer una conexión con la base de datos
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		utils.Error("Error conectándose a MongoDB: %w", err)
		return nil, err
	}

	// Verificar conexión
	err = client.Ping(ctx, nil)
	if err != nil {
		utils.Error("Error verificando la conexión a MongoDB: %w", err)
		return nil, err
	}

	utils.Info("MongoDB store initialized")
	return client, nil
}
