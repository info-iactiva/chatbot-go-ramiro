package globals

import (
	"github.com/tmc/langchaingo/schema"

)

// Cache global para almacenar los resultados de Pinecone
var PineconeResultsCache = make(map[string][]schema.Document)

// GlobalMemory es una memoria global para almacenar el historial de mensajes
var GlobalMemory = NewMemory()


