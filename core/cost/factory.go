// Package cost provides implementations of various cost functions used in change point detection,
// along with a factory for dynamic creation of these functions.
package cost

import (
	"fmt"
	"sync" // For thread-safe map access

	"github.com/theDataFlowClub/ruptures/core/base" // For the CostFunction interface
)

// costFactoryRegistry holds a map of model names to functions that construct CostFunction instances.
// This allows for dynamic registration and creation of different cost functions.
var costFactoryRegistry = make(map[string]func() base.CostFunction)

// mu protects costFactoryRegistry from concurrent access.
var mu sync.RWMutex

// RegisterCostFunction registers a CostFunction constructor with the factory.
// This function is typically called in an init() block of each cost function implementation
// to make it available via the factory.
// It is safe for concurrent use.
func RegisterCostFunction(model string, constructor func() base.CostFunction) {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := costFactoryRegistry[model]; exists {
		panic(fmt.Sprintf("cost function model '%s' already registered", model))
	}
	costFactoryRegistry[model] = constructor
	fmt.Printf("DEBUG: Cost function '%s' registered.\n", model) // <-- Añade esta línea
}

// NewCost creates and returns a new instance of a CostFunction based on its model name.
// This is the primary entry point for users to obtain a cost function dynamically.
//
// Parameters:
//
//	model: The string identifier for the desired cost function (e.g., "l2", "l1", "rbf", "entropy").
//
// Returns:
//
//	base.CostFunction: A new instance of the requested cost function.
//	error: An error if the specified model is not found in the registry.
func NewCost(model string) (base.CostFunction, error) {
	mu.RLock()
	defer mu.RUnlock()
	constructor, ok := costFactoryRegistry[model]
	if !ok {
		return nil, fmt.Errorf("no such cost function model: %s", model)
	}
	return constructor(), nil
}
