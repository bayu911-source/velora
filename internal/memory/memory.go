
package memory

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

// MemoryManager handles short-term and long-term memory.
type MemoryManager struct {
	shortTerm map[string]string
	longTerm  string // Path to the JSON file
	mu        sync.RWMutex
}

// NewMemoryManager creates a new MemoryManager.
func NewMemoryManager(longTermPath string) *MemoryManager {
	return &MemoryManager{
		shortTerm: make(map[string]string),
		longTerm:  longTermPath,
	}
}

// SaveMemory saves a key-value pair to both short-term and long-term memory.
func (m *MemoryManager) SaveMemory(key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.shortTerm[key] = value

	return m.persist()
}

// GetMemory retrieves a value from memory, checking short-term first, then long-term.
func (m *MemoryManager) GetMemory(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if value, ok := m.shortTerm[key]; ok {
		return value, true
	}

	data, err := m.load()
	if err != nil {
		return "", false
	}

	if value, ok := data[key]; ok {
		return value, true
	}

	return "", false
}

// persist saves the entire short-term memory to the long-term JSON file.
func (m *MemoryManager) persist() error {
	data, err := json.MarshalIndent(m.shortTerm, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(m.longTerm, data, 0644)
}

// load loads the long-term memory from the JSON file into short-term memory.
func (m *MemoryManager) load() (map[string]string, error) {
	if _, err := os.Stat(m.longTerm); os.IsNotExist(err) {
		return make(map[string]string), nil
	}

	data, err := ioutil.ReadFile(m.longTerm)
	if err != nil {
		return nil, err
	}

	var result map[string]string
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}
