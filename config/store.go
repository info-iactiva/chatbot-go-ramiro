package config

import (
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/vectorstores/pinecone"

	"langtools/utils"
)

// GetPineconeStore configura y retorna una instancia de Pinecone Store.
func GetPineconeStore(llm *openai.LLM) (*pinecone.Store, error) {
	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		utils.Error("error creating embedder: %w", err)
		return nil, err
	}

	// Configurar la tienda Pinecone con las variables de entorno.
	store, err := pinecone.New(
		pinecone.WithEmbedder(embedder),
		pinecone.WithAPIKey(utils.GetEnv("PINECONE_API_KEY", "")),
		pinecone.WithNameSpace(utils.GetEnv("PINECONE_NAMESPACE", "")),
		pinecone.WithHost(utils.GetEnv("PINECONE_HOST", "")),
	)

	if err != nil {
		utils.Error("error initializing Pinecone store: %w", err)
		return nil, err
	}

	utils.Info("Pinecone store initialized")

	return &store, nil
}
