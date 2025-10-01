package sensors

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

// Registry manages sensor factory functions and provides sensor creation
type Registry struct {
	factories map[string]SensorFactory
	mutex     sync.RWMutex // Protects concurrent access to factories map
}

// NewRegistry creates a new sensor registry
func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[string]SensorFactory),
	}
}

// Register adds a sensor factory function for the specified type
func (r *Registry) Register(sensorType string, factory SensorFactory) error {
	if sensorType == "" {
		return errors.New("sensor type cannot be empty")
	}

	if factory == nil {
		return errors.New("factory function cannot be nil")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check for duplicate registration
	if _, exists := r.factories[sensorType]; exists {
		return fmt.Errorf("sensor type '%s' is already registered", sensorType)
	}

	r.factories[sensorType] = factory
	return nil
}

// Create creates a new sensor instance of the specified type
func (r *Registry) Create(sensorType string) (SensorInterface, error) {
	r.mutex.RLock()
	factory, exists := r.factories[sensorType]
	r.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("unknown sensor type: %s", sensorType)
	}

	// Create new sensor instance
	sensor := factory()
	if sensor == nil {
		return nil, fmt.Errorf("factory for sensor type '%s' returned nil", sensorType)
	}

	return sensor, nil
}

// IsRegistered checks if a sensor type is registered
func (r *Registry) IsRegistered(sensorType string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	_, exists := r.factories[sensorType]
	return exists
}

// List returns a sorted list of all registered sensor types
func (r *Registry) List() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	types := make([]string, 0, len(r.factories))
	for sensorType := range r.factories {
		types = append(types, sensorType)
	}

	sort.Strings(types)
	return types
}

// Count returns the number of registered sensor types
func (r *Registry) Count() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return len(r.factories)
}

// Unregister removes a sensor type from the registry
func (r *Registry) Unregister(sensorType string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.factories[sensorType]; !exists {
		return fmt.Errorf("sensor type '%s' is not registered", sensorType)
	}

	delete(r.factories, sensorType)
	return nil
}

// Clear removes all registered sensor types
func (r *Registry) Clear() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.factories = make(map[string]SensorFactory)
}

// GetFactory returns the factory function for a sensor type (for advanced usage)
func (r *Registry) GetFactory(sensorType string) (SensorFactory, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	factory, exists := r.factories[sensorType]
	if !exists {
		return nil, fmt.Errorf("unknown sensor type: %s", sensorType)
	}

	return factory, nil
}

// RegistryInfo provides information about the registry state
type RegistryInfo struct {
	RegisteredTypes []string
	Count           int
	Built_insLoaded bool
}

// GetInfo returns information about the current registry state
func (r *Registry) GetInfo() RegistryInfo {
	types := r.List()

	// Check if built-in types are loaded
	builtInsLoaded := false
	expectedBuiltIns := []string{"numeric", "categorical", "text", "spatial"}
	if len(types) >= len(expectedBuiltIns) {
		builtInsLoaded = true
		for _, expected := range expectedBuiltIns {
			found := false
			for _, registered := range types {
				if registered == expected {
					found = true
					break
				}
			}
			if !found {
				builtInsLoaded = false
				break
			}
		}
	}

	return RegistryInfo{
		RegisteredTypes: types,
		Count:           len(types),
		Built_insLoaded: builtInsLoaded,
	}
}

// CreateMultiple creates multiple sensor instances of different types
func (r *Registry) CreateMultiple(sensorTypes []string) (map[string]SensorInterface, error) {
	sensors := make(map[string]SensorInterface)

	for _, sensorType := range sensorTypes {
		sensor, err := r.Create(sensorType)
		if err != nil {
			// Clean up already created sensors (if they implement cleanup)
			return nil, fmt.Errorf("failed to create sensor '%s': %v", sensorType, err)
		}
		sensors[sensorType] = sensor
	}

	return sensors, nil
}

// ValidateRegistration tests that a factory function works correctly
func (r *Registry) ValidateRegistration(sensorType string) error {
	// Create a test instance
	sensor, err := r.Create(sensorType)
	if err != nil {
		return fmt.Errorf("failed to create test sensor: %v", err)
	}

	// Validate basic interface compliance
	metadata := sensor.Metadata()
	if metadata.Type == "" {
		return fmt.Errorf("sensor metadata has empty type")
	}

	if metadata.SDRWidth <= 0 {
		return fmt.Errorf("sensor metadata has invalid SDR width: %d", metadata.SDRWidth)
	}

	if metadata.MaxInputSize <= 0 {
		return fmt.Errorf("sensor metadata has invalid max input size: %d", metadata.MaxInputSize)
	}

	// Test validation method
	_ = sensor.Validate()
	// It's okay for validate to fail if sensor is not configured yet
	// This just tests that the method exists and doesn't panic

	return nil
}

// Global registry instance for convenience
var globalRegistry *Registry
var globalRegistryOnce sync.Once

// GetGlobalRegistry returns the global sensor registry instance
func GetGlobalRegistry() *Registry {
	globalRegistryOnce.Do(func() {
		globalRegistry = NewRegistry()
	})
	return globalRegistry
}

// RegisterGlobal registers a sensor type in the global registry
func RegisterGlobal(sensorType string, factory SensorFactory) error {
	return GetGlobalRegistry().Register(sensorType, factory)
}

// CreateGlobal creates a sensor from the global registry
func CreateGlobal(sensorType string) (SensorInterface, error) {
	return GetGlobalRegistry().Create(sensorType)
}

// ListGlobal lists all sensor types in the global registry
func ListGlobal() []string {
	return GetGlobalRegistry().List()
}

// RegistryManager provides higher-level registry management
type RegistryManager struct {
	registries map[string]*Registry
	mutex      sync.RWMutex
}

// NewRegistryManager creates a new registry manager
func NewRegistryManager() *RegistryManager {
	return &RegistryManager{
		registries: make(map[string]*Registry),
	}
}

// CreateRegistry creates a new named registry
func (rm *RegistryManager) CreateRegistry(name string) (*Registry, error) {
	if name == "" {
		return nil, errors.New("registry name cannot be empty")
	}

	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	if _, exists := rm.registries[name]; exists {
		return nil, fmt.Errorf("registry '%s' already exists", name)
	}

	registry := NewRegistry()
	rm.registries[name] = registry
	return registry, nil
}

// GetRegistry retrieves a named registry
func (rm *RegistryManager) GetRegistry(name string) (*Registry, error) {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	registry, exists := rm.registries[name]
	if !exists {
		return nil, fmt.Errorf("registry '%s' not found", name)
	}

	return registry, nil
}

// ListRegistries returns all registry names
func (rm *RegistryManager) ListRegistries() []string {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	names := make([]string, 0, len(rm.registries))
	for name := range rm.registries {
		names = append(names, name)
	}

	sort.Strings(names)
	return names
}

// RemoveRegistry removes a named registry
func (rm *RegistryManager) RemoveRegistry(name string) error {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	if _, exists := rm.registries[name]; !exists {
		return fmt.Errorf("registry '%s' not found", name)
	}

	delete(rm.registries, name)
	return nil
}
