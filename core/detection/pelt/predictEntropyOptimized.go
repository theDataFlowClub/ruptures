package pelt

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/theDataFlowClub/ruptures/core/cost" // Asegúrate de que el import sea correcto
)

// ... (Pelt struct, NewPelt, Fit, Predict - no cambian en esta revisión)
// ... (predictL1Optimized, predictL2Optimized, predictRbfOptimized - ya implementadas)

// predictEntropyOptimized es una implementación optimizada del algoritmo PELT para costo de Entropía.
// Requiere que CostEntropy precalcule histogramas de prefijo para un cálculo de costo O(AlphabetSize).
func (p *Pelt) predictEntropyOptimized(entropyCost *cost.CostEntropy, penalty float64) ([]int, error) {
	if len(p.signal) == 0 || len(p.signal[0]) != 1 {
		// Aseguramos que la señal sea univariada (una dimensión)
		return nil, errors.New("Entropy optimized PELT requires univariate signal")
	}

	// numSamples: Número total de muestras en la señal.
	numSamples := p.nSamples

	// PELT core arrays:
	// minCostsToEnd: Almacena el costo mínimo acumulado para la señal hasta el índice actual `t`.
	minCostsToEnd := make([]float64, numSamples+1)
	// optimalPrevBreakpoints: Almacena el índice del punto de cambio óptimo anterior para cada `t`.
	optimalPrevBreakpoints := make([]int, numSamples+1)
	// pruningValues: Valores usados para la condición de poda de PELT.
	pruningValues := make([]float64, numSamples+1)

	// Inicialización de los arrays con valores "infinitos" o de inicio.
	for i := range minCostsToEnd {
		minCostsToEnd[i] = math.Inf(1)
		pruningValues[i] = math.Inf(1)
	}
	minCostsToEnd[0] = -penalty // El costo del punto inicial (punto de referencia ficticio)
	pruningValues[0] = 0.0      // Valor inicial para la poda

	// firstValidCandidate: El índice del primer punto de cambio potencial que no ha sido podado.
	firstValidCandidate := 0

	// Bucle principal de PELT: itera a través de los posibles puntos finales 'currentEnd' de los segmentos.
	for currentEnd := p.MinSize; currentEnd <= numSamples; currentEnd++ {
		minCostsToEnd[currentEnd] = math.Inf(1) // Inicializa el costo mínimo para el `currentEnd`

		// Evaluar el primer candidato no podado.
		// Este bloque es común a todas las implementaciones PELT optimizadas.
		if firstValidCandidate <= currentEnd-p.MinSize {
			prevBreakpoint := firstValidCandidate
			// Calcula el costo de entropía para el segmento [prevBreakpoint, currentEnd).
			// Aquí es crucial que `entropyCost.Error()` sea eficiente (O(AlphabetSize) o O(1)).
			segmentCost, err := entropyCost.Error(prevBreakpoint, currentEnd)
			if err != nil {
				return nil, fmt.Errorf("Pelt (Entropy): error calculating segment cost for [%d, %d): %w", prevBreakpoint, currentEnd, err)
			}

			// Actualiza el valor de poda para el 'prevBreakpoint'.
			pruningValues[prevBreakpoint] = minCostsToEnd[prevBreakpoint] + segmentCost

			// Calcula el costo total si 'prevBreakpoint' fuera el último punto de cambio óptimo.
			totalCostIfPrev := pruningValues[prevBreakpoint] + penalty

			// Si este es el mejor costo encontrado hasta ahora para 'currentEnd', actualiza.
			minCostsToEnd[currentEnd] = totalCostIfPrev
			optimalPrevBreakpoints[currentEnd] = prevBreakpoint
		}

		// Bucle para el resto de candidatos 'prevBreakpoint' (posibles puntos de cambio anteriores).
		// Recorre desde el siguiente candidato válido hasta el último punto de inicio posible
		// para un segmento de longitud `MinSize` que termina en `currentEnd`.
		for prevBreakpoint := firstValidCandidate + 1; prevBreakpoint <= currentEnd-p.MinSize; prevBreakpoint++ {
			// Calcula el costo de entropía para el segmento [prevBreakpoint, currentEnd).
			segmentCost, err := entropyCost.Error(prevBreakpoint, currentEnd)
			if err != nil {
				return nil, fmt.Errorf("Pelt (Entropy): error calculating segment cost for [%d, %d): %w", prevBreakpoint, currentEnd, err)
			}

			// Actualiza el valor de poda para el 'prevBreakpoint'.
			pruningValues[prevBreakpoint] = minCostsToEnd[prevBreakpoint] + segmentCost

			// Calcula el costo total si 'prevBreakpoint' fuera el último punto de cambio óptimo.
			totalCostIfPrev := pruningValues[prevBreakpoint] + penalty

			// Si este costo total es menor que el mínimo encontrado hasta ahora para 'currentEnd', actualizamos.
			if totalCostIfPrev < minCostsToEnd[currentEnd] {
				minCostsToEnd[currentEnd] = totalCostIfPrev
				optimalPrevBreakpoints[currentEnd] = prevBreakpoint
			}
		}

		// --- Lógica de Poda (Pruning) ---
		// Avanza 'firstValidCandidate' (el primer índice candidato no podado)
		// mientras sus valores podados sean peores que el costo mínimo actual para 'currentEnd'.
		// Esta lógica es idéntica en todas las implementaciones PELT optimizadas.
		for (firstValidCandidate < currentEnd-p.MinSize+1) && (pruningValues[firstValidCandidate] >= minCostsToEnd[currentEnd]) {
			if firstValidCandidate == 0 {
				// Si el primer candidato es 0 y se poda, saltamos al tamaño mínimo de segmento.
				firstValidCandidate = p.MinSize
			} else {
				firstValidCandidate++
			}
		}
	}

	// --- Reconstrucción de los puntos de cambio ---
	// Empezamos desde el final y retrocedemos usando los 'optimalPrevBreakpoints'.
	// Esta parte también es común a todas las implementaciones PELT.
	changePoints := []int{numSamples} // El último punto de la señal es siempre un punto de cambio
	currentChangePoint := numSamples
	for currentChangePoint != 0 {
		currentChangePoint = optimalPrevBreakpoints[currentChangePoint]
		if currentChangePoint != 0 { // No añadimos el punto inicial (0) como un punto de cambio real
			changePoints = append(changePoints, currentChangePoint)
		}
	}
	sort.Ints(changePoints) // Ordena los puntos de cambio de forma ascendente

	return changePoints, nil
}
