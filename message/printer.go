package message

import (
	"encoding/json"
	"log"

	"github.com/tmc/langchaingo/llms"

)

// Imprimir historial de mensajes en formato JSON
func PrintHistory(messageHistory []llms.MessageContent) {
	historyJSON, err := json.MarshalIndent(messageHistory, "", "    ")
	if err != nil {
		log.Fatalf("Error marshaling message history: %v", err)
	}
	log.Println("Message History:\n", string(historyJSON))
}
