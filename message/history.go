package message

import (
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

func UpdateHistory(messageHistory []llms.MessageContent, resp *llms.ContentResponse) []llms.MessageContent {
	respChoice := resp.Choices[0]
	assistantResponse := llms.TextParts(llms.ChatMessageTypeAI, respChoice.Content)
	for _, tc := range respChoice.ToolCalls {
		assistantResponse.Parts = append(assistantResponse.Parts, tc)
	}
	return append(messageHistory, assistantResponse)
}

// Simulación del historial (en producción, usa una base de datos o caché)
var history []string

// ProcessMessage procesa un mensaje entrante y actualiza el historial
func ProcessMessage(msg string) string {
	history = append(history, msg)
	response := fmt.Sprintf("Mensaje recibido: %s", msg)

	// Aquí puedes integrar lógica adicional (e.g., herramientas)
	return response
}
