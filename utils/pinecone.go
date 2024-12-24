package utils

import (
	"context"

)

type PineconeResult struct {
	Name  string
	Score float64
}

// Simula una búsqueda en Pinecone
func SearchInPinecone(ctx context.Context, query string, limit int) ([]PineconeResult, error) {
	// Simula resultados para la demostración
	results := []PineconeResult{
		{Name: "Document 1", Score: 0.95},
		{Name: "Document 2", Score: 0.89},
		{Name: "Document 3", Score: 0.85},
	}

	if limit > len(results) {
		limit = len(results)
	}

	return results[:limit], nil
}
