package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tmc/langchaingo/llms"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"langtools/config"
	"langtools/utils"
)

type SaveChatArgs struct {
	UserID string `json:"userId"`
	Name   string `json:"name"`
	Mail   string `json:"mail"`
	ChatID string `json:"chatId"`
}

type Message struct {
	Text string `json:"text"`
	User string `json:"user"`
}

func NewMongoDBToolLLM() llms.Tool {
	return llms.Tool{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "loadDataIntoMongoDB",
			Description: "Update chat history in MongoDB by uploading name and mail.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"name": map[string]any{
						"type":        "string",
						"description": "Name of the user",
					},
					"mail": map[string]any{
						"type":        "string",
						"description": "Email of the user",
					},
					"userId": map[string]any{
						"type":        "string",
						"description": "The user ID to associate with the chat history",
					},
					"chatId": map[string]any{
						"type":        "string",
						"description": "The ID of the chat to update",
					},
				},
				"required": []string{"name", "mail", "userId", "chatId"},
			},
		},
	}
}

func ExecuteMongoDBTool(ctx context.Context, argsJson string) (string, error) {
	var args SaveChatArgs
	if err := json.Unmarshal([]byte(argsJson), &args); err != nil {
		utils.Error("Error unmarshalling arguments for saveChatToMongoDB: %v", err)
		return "", fmt.Errorf("invalid arguments for saveChatToMongoDB: %v", err)
	}

	client, err := config.GetMongoDBStore()
	if err != nil {
		utils.Error("Error getting MongoDB client: %v", err)
		return "", err
	}

	collection := client.Database("dessa-chats").Collection("chats")

	// Convertir chatId a ObjectID
	chatID, err := primitive.ObjectIDFromHex(args.ChatID)
	if err != nil {
		utils.Error("Invalid ChatID format: %v", err)
		return "", fmt.Errorf("invalid ChatID format: %v", err)
	}

	// Filtro para encontrar el documento a actualizar
	filter := bson.M{"_id": chatID, "userID": args.UserID}

	// Actualización de los campos "name" y "mail"
	update := bson.M{
		"$set": bson.M{
			"name": args.Name,
			"mail": args.Mail,
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		utils.Error("Error updating chat in MongoDB: %v", err)
		return "", err
	}

	// Validar si algún documento fue actualizado
	if result.ModifiedCount == 0 {
		return "", fmt.Errorf("no chat updated; ensure ChatID and UserID are correct")
	}

	utils.Info("Chat successfully updated in MongoDB")
	response := fmt.Sprintf("Chat successfully updated in MongoDB with ID: %v", chatID.Hex())

	return response, nil
}

type UpdateChatArgs struct {
	UserID string `json:"userId"`
	ChatID string `json:"chatId"`
}

// func NewMongoDBToolUpdateLLM() llms.Tool {
// 	return llms.Tool{
// 		Type: "function",
// 		Function: &llms.FunctionDefinition{
// 			Name:        "updateChatInMongoDB",
// 			Description: "Update chat history in MongoDB by appending new messages.",
// 			Parameters: map[string]any{
// 				"type": "object",
// 				"properties": map[string]any{
// 					"userId": map[string]any{
// 						"type":        "string",
// 						"description": "The user ID associated with the chat history",
// 					},
// 					"chatId": map[string]any{
// 						"type":        "string",
// 						"description": "The ID of the chat to update",
// 					},
// 				},
// 				"required": []string{"userId", "chatId"},
// 			},
// 		},
// 	}
// }

// func ExecuteUpdateMongoDBTool(ctx context.Context, argsJson string, messageHistory []llms.MessageContent) (string, error) {
// 	var args UpdateChatArgs
// 	if err := json.Unmarshal([]byte(argsJson), &args); err != nil {
// 		utils.Error("Error unmarshalling arguments for updateChatInMongoDB: %v", err)
// 		return "", fmt.Errorf("invalid arguments for updateChatInMongoDB: %v", err)
// 	}

// 	client, err := config.GetMongoDBStore()
// 	if err != nil {
// 		utils.Error("Error getting MongoDB client: %v", err)
// 		return "", err
// 	}

// 	// Convert ChatID to ObjectID
// 	chatID, err := primitive.ObjectIDFromHex(args.ChatID)
// 	if err != nil {
// 		utils.Error("Invalid ChatID format: %v", err)
// 		return "", fmt.Errorf("invalid ChatID format: %v", err)
// 	}

// 	collection := client.Database("dessa-chats").Collection("chats")

// 	// Append new messages to the history
// 	update := bson.M{
// 		"$set": bson.M{
// 			"history": formatMessageHistory(messageHistory),
// 		},
// 	}

// 	filter := bson.M{"_id": chatID, "userID": args.UserID}

// 	result, err := collection.UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		utils.Error("Error updating chat in MongoDB: %v", err)
// 		return "", err
// 	}

// 	if result.ModifiedCount == 0 {
// 		return "", fmt.Errorf("no chat updated; ensure ChatID and UserID are correct")
// 	}

// 	utils.Info("Chat successfully updated in MongoDB")
// 	return "Chat successfully updated in MongoDB", nil
// }

func formatMessageHistory(messageHistory []llms.MessageContent) []Message {
	formattedHistory := []Message{}
	for _, msg := range messageHistory {
		jsonText, err := msg.MarshalJSON()
		if err != nil {
			utils.Error("Error marshalling message to JSON: %v", err)
			return nil
		}
		formattedHistory = append(formattedHistory, Message{
			Text: string(jsonText),
			User: string(msg.Role),
		})
	}
	return formattedHistory
}

func CreateEmptyChat(ctx context.Context, userID string) (string, error) {

	client, err := config.GetMongoDBStore()
	if err != nil {
		return "", fmt.Errorf("error getting MongoDB client: %v", err)
	}

	collection := client.Database("dessa-chats").Collection("chats")
	chat := bson.M{
		"userID":    userID,
		"history":   []Message{},
		"createdAt": primitive.NewDateTimeFromTime(time.Now()),
	}

	result, err := collection.InsertOne(ctx, chat)
	if err != nil {
		return "", fmt.Errorf("Error inserting empty chat: %v", err)
	}

	chatID := result.InsertedID.(primitive.ObjectID).Hex()
	return chatID, nil
}

func UpdateChatHistory(chatID, userID string, messageHistory []llms.MessageContent) error {
	ctx := context.Background()

	// Obtener el cliente de MongoDB
	client, err := config.GetMongoDBStore()
	if err != nil {
		return fmt.Errorf("error getting MongoDB client: %v", err)
	}

	// Seleccionar la colección de chats
	collection := client.Database("dessa-chats").Collection("chats")

	// Convertir chatID a ObjectID
	objectID, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return fmt.Errorf("invalid ChatID format: %v", err)
	}

	// Formatear el historial de mensajes
	formattedHistory := formatMessageHistory(messageHistory)

	// Construir el filtro para identificar la conversación
	filter := bson.M{
		"_id":    objectID,
		"userID": userID,
	}

	// Crear la actualización que añade el historial
	update := bson.M{
		"$set": bson.M{
			"history":   formattedHistory,
			"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	// Realizar la actualización en MongoDB
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating chat history in MongoDB: %v", err)
	}

	// Validar si se actualizó algún documento
	if result.ModifiedCount == 0 {
		return fmt.Errorf("no chat updated; ensure ChatID and UserID are correct")
	}

	utils.Info("Chat history successfully updated in MongoDB", zap.String("ID", objectID.Hex()))

	return nil
}
