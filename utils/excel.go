package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GenerateExcel() (string, string, error) {
	mongoURI := "mongodb+srv://admin:admin@cluster0.zgyky.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Error al conectar a MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Error al desconectar MongoDB: %v", err)
		}
	}()

	// Conexión a la base de datos y colección
	db := client.Database("dessa-chats")
	collection := db.Collection("chats")
	Info("Conexión exitosa a MongoDB.")

	// Filtrar documentos por fecha actual (día del servidor)
	// startOfDay := time.Now().Truncate(24 * time.Hour)
	// endOfDay := startOfDay.Add(24 * time.Hour)

	// filter := bson.M{
	// 	"createdAt": bson.M{
	// 		"$gte": startOfDay,
	// 		"$lt":  endOfDay,
	// 	},
	// }

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatalf("Error al obtener documentos: %v", err)
	}
	defer cursor.Close(context.TODO())

	var records []bson.M
	if err := cursor.All(context.TODO(), &records); err != nil {
		log.Fatalf("Error al leer documentos: %v", err)
	}

	if len(records) == 0 {
		fmt.Println("No hay datos para exportar del día actual.")
		return "", "", nil
	}

	// Crear archivo Excel
	f := excelize.NewFile()
	sheetName := "Sheet1"

	// Procesar datos
	headers := []string{"ID", "Fecha Creación", "Conversación", "Correo", "Nombre", "Actualizado"}
	for colIndex, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+colIndex)))
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			log.Fatalf("Error al escribir encabezado en Excel: %v", err)
		}
	}

	// Iterar sobre los registros y escribir filas en Excel
	for rowIndex, record := range records {
		colIndex := 0 // Reiniciar el índice de columna para cada fila

		// ID
		if id, ok := record["_id"]; ok {
			cell := fmt.Sprintf("%s%d", string(rune('A'+colIndex)), rowIndex+2)
			f.SetCellValue(sheetName, cell, id)
			colIndex++
		}

		// Fecha Creación
		if createdAt, ok := record["createdAt"].(primitive.DateTime); ok {
			cell := fmt.Sprintf("%s%d", string(rune('A'+colIndex)), rowIndex+2)
			f.SetCellValue(sheetName, cell, createdAt.Time().Format("2006-01-02 15:04:05"))
			colIndex++
		}

		// Conversación
		if history, ok := record["history"].(primitive.A); ok {
			log.Printf("Historial de conversación: ENTRO")
			conversation := []string{}
			for _, entry := range history {
				// log.Printf("Tipo de entrada en history: %T, Valor: %+v", entry, entry)

				// Revisar si es un map[string]interface{} o un tipo BSON específico
				switch v := entry.(type) {
				case primitive.M: // Manejar primitive.M
					log.Printf("Es un primitive.M: %+v", v)

					// Extraer el campo "text" como JSON
					if textField, ok := v["text"].(string); ok {
						var textData map[string]interface{}
						if err := json.Unmarshal([]byte(textField), &textData); err == nil {
							log.Printf("Texto deserializado: %+v", textData)

							// Extraer role y text del JSON deserializado
							if role, ok := textData["role"].(string); ok {
								if message, ok := textData["text"].(string); ok {
									if role == "human" {
										conversation = append(conversation, fmt.Sprintf("- Usuario: %s", message))
									} else if role == "ai" {
										conversation = append(conversation, fmt.Sprintf("- AI: %s", message))
									}
								}
							} else {
								log.Printf("No existe el campo 'role' en el texto deserializado: %+v", textData)
							}
						} else {
							log.Printf("Error al deserializar el campo 'text': %v", err)
						}
					} else {
						log.Printf("El campo 'text' no es una string o no existe.")
					}
				case map[string]interface{}: // En caso de que sea un map normal
					log.Printf("Es un map[string]interface{}: %+v", v)
					if role, ok := v["role"].(string); ok {
						if text, ok := v["text"].(string); ok {
							if role == "human" {
								conversation = append(conversation, fmt.Sprintf("- Usuario: %s", text))
							} else if role == "ai" {
								conversation = append(conversation, fmt.Sprintf("- AI: %s", text))
							}
						}
					}
				default:
					log.Printf("Tipo inesperado en history: %T, Valor: %+v", entry, entry)
				}
			}
			cell := fmt.Sprintf("%s%d", string(rune('A'+colIndex)), rowIndex+2)
			f.SetCellValue(sheetName, cell, strings.Join(conversation, "\n"))
			colIndex++
		} else {
			log.Printf("El campo history no es un primitive.A o está vacío para el documento ID: %v", record["_id"])
		}

		// Correo
		if mail, ok := record["mail"].(string); ok {
			cell := fmt.Sprintf("%s%d", string(rune('A'+colIndex)), rowIndex+2)
			f.SetCellValue(sheetName, cell, mail)
			colIndex++
		}

		// Nombre
		if name, ok := record["name"].(string); ok {
			cell := fmt.Sprintf("%s%d", string(rune('A'+colIndex)), rowIndex+2)
			f.SetCellValue(sheetName, cell, name)
			colIndex++
		}

		// Actualizado
		if updatedAt, ok := record["updatedAt"].(primitive.DateTime); ok {
			cell := fmt.Sprintf("%s%d", string(rune('A'+colIndex)), rowIndex+2)
			f.SetCellValue(sheetName, cell, updatedAt.Time().Format("2006-01-02 15:04:05"))
			colIndex++
		}
	}

	// Generar nombre de archivo con la fecha actual
	currentDate := time.Now().Format("02_01_2006")
	fileName := fmt.Sprintf("Reporte_Chats_Dessa_%s.xlsx", currentDate)
	outputFile := fmt.Sprintf("./%s", fileName)

	// Guardar archivo
	if err := f.SaveAs(outputFile); err != nil {
		log.Fatalf("Error al guardar archivo Excel: %v", err)
	}

	Info(fmt.Sprintf("Datos exportados a %s con éxito.\n", outputFile))
	return outputFile, fileName, nil
}
