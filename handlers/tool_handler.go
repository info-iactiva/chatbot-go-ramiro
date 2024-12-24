package handlers

import (
	"context"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/vectorstores/pinecone"

	"langtools/tools"
	"langtools/utils"
)

func ExecuteToolCalls(ctx context.Context, messageHistory []llms.MessageContent, resp *llms.ContentResponse, store *pinecone.Store) []llms.MessageContent {
	for _, toolCall := range resp.Choices[0].ToolCalls {
		utils.Info("Executing tool call: " + toolCall.FunctionCall.Name)
		response := HandleToolCall(ctx, tools.GetToolByName(toolCall.FunctionCall.Name), toolCall.FunctionCall.Arguments, store)

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
		messageHistory = append(messageHistory, toolResponse)
	}
	return messageHistory
}

func HandleToolCall(ctx context.Context, tool llms.Tool, args string, store *pinecone.Store) string {
	pineconeTool := tools.NewPineconeTool(store)

	// Crear una instancia de la herramienta.
	switch tool.Function.Name {
	case "getCurrentWeather":
		response, err := tools.ExecuteWeatherTool(args)
		if err != nil {
			log.Printf("Error executing weather tool: %v", err)
			return "Error fetching weather."
		}
		return response
	case "pineconeSearch":
		response, err := pineconeTool.Execute(ctx, args)
		if err != nil {
			log.Printf("Error executing Pinecone search: %v", err)
			return "Error executing Pinecone search."
		}
		return response
	default:
		return "Tool not supported."
	}
}
