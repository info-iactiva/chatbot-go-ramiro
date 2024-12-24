package ws

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/tmc/langchaingo/schema"
)

// Mapa para renombrar las claves de Metadata
var metadataKeyMap = map[string]string{
	"Aspecto Variante #1":               "variantAspect1",
	"Aspecto Variante #2":               "variantAspect2",
	"Aspecto Variante #3":               "variantAspect3",
	"Categoria":                         "category",
	"Descripcion":                       "description",
	"Enlace de imagen":                  "imageUrl",
	"Etiquetas/Tags":                    "tags",
	"Opciones para Aspecto Variante #1": "variantOptions1",
	"Opciones para Aspecto Variante #2": "variantOptions2",
	"Opciones para Aspecto Variante #3": "variantOptions3",
	"Peso (kg)":                         "weight",
	"Precio (MXN)":                      "price",
	"Requiere envio":                    "requiresShipping",
	"SKU":                               "sku",
	"Tipo":                              "type",
	"Titulo":                            "title",
	"Vendedor":                          "seller",
}

// Función para transformar Metadata
func transformMetadata(metadata map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})
	for originalKey, newKey := range metadataKeyMap {
		if value, exists := metadata[originalKey]; exists {
			transformed[newKey] = value
		}
	}
	return transformed
}

// Función para formatear los resultados
func FormatResults(results []schema.Document) []map[string]interface{} {
	formatted := make([]map[string]interface{}, 0, len(results))
	for _, result := range results {
		formatted = append(formatted, transformMetadata(result.Metadata))
	}
	return formatted
}

type WSResponse struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Records []map[string]interface{} `json:"records,omitempty"`
}

// Función para enviar la respuesta al cliente
func SendResponse(conn *websocket.Conn, status string, message string, results []map[string]interface{}) error {
	response := WSResponse{
		Status:  status,
		Message: message,
		Records: results,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return err
	}

	if err := conn.WriteMessage(websocket.TextMessage, jsonResponse); err != nil {
		log.Printf("Error writing message: %v", err)
		return err
	}

	return nil
}
