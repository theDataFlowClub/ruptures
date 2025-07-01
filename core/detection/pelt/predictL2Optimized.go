package pelt

import (
	"errors"
	"math"
	"sort"

	"github.com/theDataFlowClub/ruptures/core/cost"
)

// iinspired by: https://arxiv.org/pdf/1101.1438
// Optimal detection of changepoints with a linear computational cost
// Killick, R., Fearnhead, P. and Eckley, I.A.∗
//
// predictL2Optimized es una implementación optimizada del algoritmo PELT para costo L2.
// Utiliza sumas acumuladas para calcular el costo de segmento en O(1), optimizado para señales univariadas.
func (p *Pelt) predictL2Optimized(l2Cost *cost.CostL2, penalty float64) ([]int, error) {
	if len(p.signal) == 0 || len(p.signal[0]) != 1 {
		// Aseguramos que la señal sea univariada (una dimensión)
		return nil, errors.New("L2 optimized PELT requires univariate signal")
	}

	numSamples := p.nSamples
	// signalValues: Extraemos los valores de la señal univariada para un acceso más fácil
	signalValues := make([]float64, numSamples)
	for i := 0; i < numSamples; i++ {
		signalValues[i] = p.signal[i][0]
	}

	// prefixSums: Almacena las sumas acumuladas de los valores de la señal.
	// prefixSums[k] = sum(signalValues[0]...signalValues[k-1])
	prefixSums := make([]float64, numSamples+1)
	// prefixSquares: Almacena las sumas acumuladas de los cuadrados de los valores de la señal.
	// prefixSquares[k] = sum(signalValues[0]^2...signalValues[k-1]^2)
	prefixSquares := make([]float64, numSamples+1)

	// Calculamos las sumas acumuladas y sumas de cuadrados una vez.
	for i := 0; i < numSamples; i++ {
		prefixSums[i+1] = prefixSums[i] + signalValues[i]
		prefixSquares[i+1] = prefixSquares[i] + signalValues[i]*signalValues[i]
	}

	// minCostsToEnd: Almacena el costo mínimo acumulado para la señal hasta el índice actual.
	minCostsToEnd := make([]float64, numSamples+1)
	// optimalPrevBreakpoints: Almacena el índice del punto de cambio óptimo anterior para cada índice.
	optimalPrevBreakpoints := make([]int, numSamples+1)
	// pruningValues: Valores usados para la condición de poda en PELT.
	pruningValues := make([]float64, numSamples+1)

	// Inicialización de los arrays con valores "infinitos" o de inicio.
	for i := range minCostsToEnd {
		minCostsToEnd[i] = math.Inf(1)
		pruningValues[i] = math.Inf(1)
	}
	minCostsToEnd[0] = -penalty // Costo inicial, ajustado por la penalización
	pruningValues[0] = 0.0      // Valor inicial para la poda

	// firstValidCandidate: El índice del primer punto de cambio potencial que no ha sido podado.
	firstValidCandidate := 0

	// Bucle principal de PELT: itera a través de los posibles puntos finales 'currentEnd' de los segmentos.
	for currentEnd := p.MinSize; currentEnd <= numSamples; currentEnd++ {
		minCostsToEnd[currentEnd] = math.Inf(1) // Inicializa el costo mínimo para el `currentEnd`

		// Evaluar el primer candidato no podado.
		if firstValidCandidate <= currentEnd-p.MinSize {
			prevBreakpoint := firstValidCandidate
			// Calcula el costo L2 para el segmento [prevBreakpoint, currentEnd)
			segmentCost := calculateL2SegmentCostFromPrefixSums(
				prefixSums,
				prefixSquares,
				prevBreakpoint,
				currentEnd,
			)

			// Actualiza el valor de poda para el 'prevBreakpoint'
			pruningValues[prevBreakpoint] = minCostsToEnd[prevBreakpoint] + segmentCost

			// Calcula el costo total si 'prevBreakpoint' fuera el último punto de cambio óptimo.
			totalCostIfPrev := pruningValues[prevBreakpoint] + penalty

			// Si este es el mejor costo encontrado hasta ahora para 'currentEnd', actualiza.
			minCostsToEnd[currentEnd] = totalCostIfPrev
			optimalPrevBreakpoints[currentEnd] = prevBreakpoint
		}

		// Bucle para el resto de candidatos 'prevBreakpoint' (posibles puntos de cambio anteriores)
		// Recorre desde el siguiente candidato válido hasta el último punto de inicio posible
		// para un segmento de longitud `MinSize` que termina en `currentEnd`.
		for prevBreakpoint := firstValidCandidate + 1; prevBreakpoint <= currentEnd-p.MinSize; prevBreakpoint++ {
			// Calcula el costo L2 para el segmento [prevBreakpoint, currentEnd)
			segmentCost := calculateL2SegmentCostFromPrefixSums(
				prefixSums,
				prefixSquares,
				prevBreakpoint,
				currentEnd,
			)

			// Actualiza el valor de poda para el 'prevBreakpoint'
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

// calculateL2SegmentCostFromPrefixSums calcula el costo L2 para un segmento [startIdx, endIdx)
// utilizando sumas acumuladas y sumas de cuadrados precalculadas.
func calculateL2SegmentCostFromPrefixSums(
	prefixSums []float64,
	prefixSquares []float64,
	startIdx, endIdx int,
) float64 {
	segmentLength := float64(endIdx - startIdx)
	if segmentLength == 0 {
		return 0.0 // Un segmento vacío tiene costo cero
	}

	// Suma de los valores en el segmento [startIdx, endIdx)
	sumValues := prefixSums[endIdx] - prefixSums[startIdx]
	// Suma de los cuadrados de los valores en el segmento [startIdx, endIdx)
	sumSquares := prefixSquares[endIdx] - prefixSquares[startIdx]

	// Fórmula del costo L2 (varianza dentro del segmento, multiplicada por la longitud para ser consistente)
	// Costo = Sum(y_i^2) - (Sum(y_i))^2 / N
	cost := sumSquares - (sumValues*sumValues)/segmentLength
	return cost
}
