package handlers

import (
	"context"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/pinecone"

	"langtools/tools"
	"langtools/utils"
)

func ExecuteToolCalls(ctx context.Context, messageHistory []llms.MessageContent, resp *llms.ContentResponse, store *pinecone.Store) ([]llms.MessageContent, []schema.Document) {
	var records []schema.Document
	var response string
	for _, toolCall := range resp.Choices[0].ToolCalls {
		utils.Info("Executing tool call: " + toolCall.FunctionCall.Name)
		response, records = HandleToolCall(ctx, tools.GetToolByName(toolCall.FunctionCall.Name), toolCall.FunctionCall.Arguments, store)

		toolResponse := llms.MessageContent{
			Role: llms.ChatMessageTypeTool,
			Parts: []llms.ContentPart{
				llms.ToolCallResponse{
					ToolCallID: toolCall.ID,
					Name:       toolCall.FunctionCall.Name,
					Content:    response,
				},
			},
		}
		log.Printf("Tool call response: %v", toolResponse)
		messageHistory = append(messageHistory, toolResponse)
	}
	return messageHistory, records
}

func HandleToolCall(ctx context.Context, tool llms.Tool, args string, store *pinecone.Store) (string, []schema.Document) {
	pineconeTool := tools.NewPineconeTool(store)

	results := []schema.Document{}

	switch tool.Function.Name {
	case "pineconeSearch":
		log.Printf("Executing Pinecone search with args: %v", args)
		response, results, err := pineconeTool.Execute(ctx, args)
		log.Printf("Pinecone search results: %v", results)
		if err != nil {
			log.Printf("Error executing Pinecone search: %v", err)
			return "Error executing Pinecone search.", results
		}

		return response, results
	case "loadDataIntoMongoDB":
		response, err := tools.ExecuteMongoDBTool(ctx, args)
		if err != nil {
			log.Printf("Error executing MongoDB tool: %v", err)
			return "Error saving chat to MongoDB.", results
		}
		return response, results

	// case "updateChatInMongoDB":
	// 	response, err := tools.ExecuteUpdateMongoDBTool(ctx, args, messageHistory)
	// 	if err != nil {
	// 		log.Printf("Error executing MongoDB tool: %v", err)
	// 		return "Error updating chat in MongoDB.", results
	// 	}
	// 	return response, results
	default:
		return "Tool not supported.", results
	}
}
