package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
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
	fmt.Println("Conexión exitosa a MongoDB.")

	// Obtener registros
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
		fmt.Println("La colección está vacía. No hay datos para exportar.")
		return "", "", nil
	}

	// Crear archivo Excel
	f := excelize.NewFile()
	sheetName := "Sheet1"

	// Convertir registros BSON a JSON
	data, err := json.Marshal(records)
	if err != nil {
		log.Fatalf("Error al convertir registros a JSON: %v", err)
	}

	// Convertir JSON a un mapa para procesar columnas
	var rows []map[string]interface{}
	if err := json.Unmarshal(data, &rows); err != nil {
		log.Fatalf("Error al convertir JSON a mapa: %v", err)
	}

	// Procesar columnas y filas
	headers := []string{}
	if len(rows) > 0 {
		for header := range rows[0] {
			headers = append(headers, header)
		}
	}

	// Escribir encabezados
	for colIndex, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+colIndex)))
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			log.Fatalf("Error al escribir encabezado en Excel: %v", err)
		}
	}

	// Escribir datos
	for rowIndex, row := range rows {
		for colIndex, header := range headers {
			cell := fmt.Sprintf("%s%d", string(rune('A'+colIndex)), rowIndex+2)
			value := row[header]

			// Formatear fechas
			if header == "updatedAt" || header == "createdAt" {
				if t, ok := value.(string); ok {
					parsedTime, err := time.Parse(time.RFC3339, t)
					if err == nil {
						value = parsedTime.Format("2006-01-02 15:04:05")
					}
				}
			}

			if err := f.SetCellValue(sheetName, cell, value); err != nil {
				log.Fatalf("Error al escribir valor en Excel: %v", err)
			}
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

	fmt.Printf("Datos exportados a %s con éxito.\n", outputFile)
	return outputFile, fileName, nil
}
