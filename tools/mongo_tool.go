package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"

	"langtools/config"
	"langtools/globals"
	"langtools/utils"
)

type SaveChatArgs struct {
	UserID     string `json:"userId"`
	ClientName string `json:"clientName"`
	ClientMail string `json:"clientMail"`
}

type Message struct {
	Text string `json:"text"`
	User string `json:"user"`
}

func NewMongoDBToolLLM() llms.Tool {
	return llms.Tool{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "saveChatToMongoDB",
			Description: "Save chat history to MongoDB, including user information and message history.",
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
				},
				"required": []string{"name", "mail", "userId"},
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

	// Create formatted history object
	chatHistory := globals.GlobalMemory.GetHistory(args.UserID)
	formattedHistory := []Message{}
	for _, msg := range chatHistory {
		jsonText, err := msg.MarshalJSON()
		if err != nil {
			utils.Error("Error marshalling message to JSON: %v", err)
			return "", err
		}
		formattedHistory = append(formattedHistory, Message{
			Text: string(jsonText),
			User: string(msg.Role),
		})
	}

	// Save chat to MongoDB
	collection := client.Database("dessa-chats").Collection("chats")
	chat := bson.M{
		"userID":     args.UserID,
		"clientName": args.ClientName,
		"clientMail": args.ClientMail,
		"history":    formattedHistory,
	}

	result, err := collection.InsertOne(ctx, chat)
	if err != nil {
		utils.Error("Error inserting chat to MongoDB: %v", err)
		return "", err
	}

	chatID := result.InsertedID
	utils.Info("Chat successfully saved to MongoDB", zap.String("ID", fmt.Sprintf("%v", chatID)))
	return fmt.Sprintf("Chat successfully saved to MongoDB with ID: %v", chatID), nil
}
