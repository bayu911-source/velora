
package memory

import (
	"sync"
)

type MemoryManager struct {
	storage map[string]string
	mu      sync.Mutex
}

func NewMemoryManager() *MemoryManager {
	return &MemoryManager{
		storage: make(map[string]string),
	}
}

func (m *MemoryManager) Set(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.storage[key] = value
}

func (m *MemoryManager) Get(key string) (string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	value, ok := m.storage[key]
	return value, ok
}
