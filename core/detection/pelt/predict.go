package pelt

import (
	"errors"
	"fmt"

	"github.com/theDataFlowClub/ruptures/core/cost"
)

// Predict es la función principal que selecciona la implementación optimizada
// basada en el tipo de CostFunction.
func (p *Pelt) Predict(penalty float64) ([]int, error) {
	// Validaciones generales antes de cualquier implementación específica
	if p.signal == nil || p.nSamples == 0 {
		return nil, errors.New("Pelt: detector not fitted. Call Fit() first.")
	}
	if penalty <= 0 {
		return nil, errors.New("Pelt: penalty must be greater than 0.")
	}
	if p.MinSize < 1 {
		return nil, errors.New("Pelt: min_size must be at least 1.")
	}

	// Selecciona la función de predicción optimizada basada en el tipo de CostFunction
	switch concreteCost := p.Cost.(type) {
	case *cost.CostRbf:
		fmt.Println("Pelt: Usando implementación optimizada para CostRbf.")
		return p.predictRbfOptimized(concreteCost, penalty)
	case *cost.CostL1:
		fmt.Println("Pelt: Usando implementación optimizada para CostL1 (por implementar).")
		return p.predictL1Optimized(concreteCost, penalty) // Llama a la función específica de L1
	case *cost.CostL2:
		fmt.Println("Pelt: Usando implementación optimizada para CostL2 (por implementar).")
		return p.predictL2Optimized(concreteCost, penalty) // Llama a la función específica de L2
	case *cost.CostEntropy: // ¡NUEVO CASO!
		fmt.Println("Pelt: Usando implementación genérica para CostEntropy.")
		// Para Entropy, inicialmente usas la función genérica,
		// ya que la optimización O(1) con prefix histograms es más compleja.
		// podrías renombrar 'predictGeneric' a 'predictBasic' o similar si te gusta.
		return p.predictEntropyOptimized(concreteCost, penalty) // Necesitarás implementar predictGeneric si aún no lo tienes.
	default:
		// En caso de que se pase una función de costo no reconocida o no optimizada
		return nil, fmt.Errorf("Pelt: la función de costo '%s' no tiene una implementación Predict optimizada. Considera añadirla o usar una función genérica.", p.Cost.Model())
	}
}
