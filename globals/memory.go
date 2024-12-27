package globals

import (
	"sync"

	"github.com/tmc/langchaingo/llms"
)

// Memoria de conversación por usuario
type Memory struct {
	mu            sync.RWMutex
	conversations map[string][]llms.MessageContent
}

func NewMemory() *Memory {
	return &Memory{
		conversations: make(map[string][]llms.MessageContent),
	}
}

// Obtener el historial de un usuario
func (m *Memory) GetHistory(userID string) []llms.MessageContent {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.conversations[userID]
}

// Actualizar el historial de un usuario
func (m *Memory) UpdateHistory(userID string, history []llms.MessageContent) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.conversations[userID] = history
}
