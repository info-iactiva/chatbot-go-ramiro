package tools

import (
	"github.com/tmc/langchaingo/llms"

)

// Registrar herramientas disponibles
func RegisterTools() []llms.Tool {
	return []llms.Tool{
		NewPineconeToolLLM(),
		NewMongoDBToolLLM(),
	}
}

// Buscar herramienta por nombre
func GetToolByName(name string) llms.Tool {
	registeredTools := RegisterTools() // Actualiza con las herramientas registradas
	for _, tool := range registeredTools {
		if tool.Function.Name == name {
			return tool
		}
	}
	return llms.Tool{}
}
