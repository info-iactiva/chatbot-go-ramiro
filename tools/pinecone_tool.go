package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/pinecone"

	"langtools/globals"
	"langtools/utils"
)

type PineconeSearchArgs struct {
	Query  string `json:"query"`
	Limit  int    `json:"limit"`
	UserID string `json:"userId"`
}

type PineconeTool struct {
	store *pinecone.Store
}

func NewPineconeTool(store *pinecone.Store) *PineconeTool {
	return &PineconeTool{store: store}
}

// Constructor para PineconeTool
func NewPineconeToolLLM() llms.Tool {
	return llms.Tool{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "pineconeSearch",
			Description: "Busca en la db vectorial de la muebleria Dessa",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"query": map[string]any{
						"type":        "string",
						"description": "The query string to search in Pinecone",
					},
					"limit": map[string]any{
						"type":        "integer",
						"description": "Maximum number of results to return",
						"minimum":     1,
					},
					"userId": map[string]any{
						"type":        "string",
						"description": "The user ID to associate with the search results",
					},
				},
				"required": []string{"query", "limit", "userId"},
			},
		},
	}
}

// Ejecutar búsqueda con Pinecone
func (t *PineconeTool) Execute(ctx context.Context, argsJson string) (string, []schema.Document, error) {
	var args PineconeSearchArgs
	if err := json.Unmarshal([]byte(argsJson), &args); err != nil {
		return "", nil, fmt.Errorf("invalid arguments for Pinecone search: %v", err)
	}

	results, err := t.store.SimilaritySearch(ctx, args.Query, args.Limit)
	if err != nil {
		return "", nil, fmt.Errorf("error performing Pinecone search: %w", err)
	}

	if args.UserID == "" {
		args.UserID = "default"
	}

	utils.Info("Saving Pinecone results for user: " + args.UserID)
	globals.PineconeResultsCache[args.UserID] = results

	response, err := json.Marshal(results)
	if err != nil {
		return "", nil, fmt.Errorf("error marshaling results: %w", err)
	}

	return string(response), results, nil
}
