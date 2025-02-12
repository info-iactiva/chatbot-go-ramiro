package config

import (
	"os"

	"github.com/tmc/langchaingo/llms/openai"

)

var (
	OpenAIModel      = os.Getenv("OPENAI_API_KEY")
	PineconeAPIKey   = os.Getenv("PINECONE_API_KEY")
	PineconeIndexURL = os.Getenv("PINECONE_INDEX_URL")
)

func InitOpenAI() (*openai.LLM, error) {
	return openai.New(openai.WithModel(OpenAIModel))
}

// LoadPrompt carga el contenido del archivo estático prompt.txt
func LoadPrompt(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
