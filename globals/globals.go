package globals

import (
	"github.com/tmc/langchaingo/schema"

)

// Cache global para almacenar los resultados de Pinecone
var PineconeResultsCache = make(map[string][]schema.Document)
