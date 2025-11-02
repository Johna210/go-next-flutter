package core

import (
	"fmt"
	"sync"
)

// EntityProvider defines the interface for modules to provide their entities
type EntityProvider interface {
	Entities() []any
	ModuleName() string
}

// SchemaManager manages all entities from different modules
type SchemaManager struct {
	mu        sync.RWMutex
	providers map[string]EntityProvider
}

// NewSchemaManager creates a new schema manager
func NewSchemaManager() *SchemaManager {
	return &SchemaManager{
		providers: make(map[string]EntityProvider),
	}
}

// RegisterProvider registers an entity provider (module)
func (sm *SchemaManager) RegisterProvider(provider EntityProvider) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	moduleName := provider.ModuleName()

	if _, exists := sm.providers[moduleName]; exists {
		return fmt.Errorf("module '%s' already registered", moduleName)
	}

	sm.providers[moduleName] = provider
	return nil
}

// GetAllEntities returns all registered entities from all modules
func (sm *SchemaManager) GetAllEntities() []interface{} {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var allEntities []interface{}
	for _, provider := range sm.providers {
		allEntities = append(allEntities, provider.Entities()...)
	}

	return allEntities
}

// GetEntitiesByModules returns entities for specific modules
func (sm *SchemaManager) GetEntitiesByModules(moduleNames ...string) []interface{} {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var entities []interface{}
	for _, name := range moduleNames {
		if provider, exists := sm.providers[name]; exists {
			entities = append(entities, provider.Entities()...)
		}
	}
	return entities
}

// ListModules returns all registered module names
func (sm *SchemaManager) ListModules() []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	names := make([]string, 0, len(sm.providers))
	for name := range sm.providers {
		names = append(names, name)
	}
	return names
}

// GetModuleInfo returns detailed info about all modules
func (sm *SchemaManager) GetModuleInfo() map[string]int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	info := make(map[string]int)
	for name, provider := range sm.providers {
		info[name] = len(provider.Entities())
	}

	return info
}
