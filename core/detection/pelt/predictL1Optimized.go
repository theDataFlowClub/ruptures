package pelt

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/theDataFlowClub/ruptures/core/cost"
)

// predictL1Optimized es una implementación pragmática del algoritmo PELT para costo L1.
// Funciona para señales univariadas, calculando la mediana de cada segmento.
func (p *Pelt) predictL1Optimized(l1Cost *cost.CostL1, penalty float64) ([]int, error) {
	if len(p.signal) == 0 || len(p.signal[0]) != 1 {
		// Aseguramos que la señal sea univariada (una dimensión)
		return nil, errors.New("L1 optimized PELT requires univariate signal")
	}

	// Extraemos la señal univariada para facilitar el acceso
	numSamples := p.nSamples
	signalValues := make([]float64, numSamples)
	for i := 0; i < numSamples; i++ {
		signalValues[i] = p.signal[i][0] // Asume señal[i] es []float64{valor}
	}

	// Arrays principales del algoritmo PELT
	// minCostsToEnd: Almacena el costo mínimo acumulado hasta el índice 'i'.
	minCostsToEnd := make([]float64, numSamples+1)
	// optimalPrevBreakpoints: Almacena el índice del último punto de cambio óptimo antes del índice 'i'.
	optimalPrevBreakpoints := make([]int, numSamples+1)
	// pruningValues: Valores usados para la condición de poda.
	pruningValues := make([]float64, numSamples+1)

	// Inicialización de los arrays con valores "infinitos" o de inicio.
	for i := range minCostsToEnd {
		minCostsToEnd[i] = math.Inf(1)
		pruningValues[i] = math.Inf(1)
	}
	minCostsToEnd[0] = -penalty // El costo del punto inicial (punto de referencia ficticio)
	pruningValues[0] = 0.0      // Valor de poda inicial

	// s_min: El primer índice candidato que no ha sido podado.
	firstValidCandidate := 0

	// Bucle principal de PELT: itera a través de los posibles puntos finales 't' de los segmentos.
	for currentEnd := p.MinSize; currentEnd <= numSamples; currentEnd++ {
		minCostsToEnd[currentEnd] = math.Inf(1) // Inicializamos el costo mínimo para el `currentEnd`

		// Evaluar el primer candidato no podado
		if firstValidCandidate <= currentEnd-p.MinSize {
			prevBreakpoint := firstValidCandidate
			// Calcula el costo L1 para el segmento [prevBreakpoint, currentEnd)
			segmentCost, err := l1Cost.Error(prevBreakpoint, currentEnd)
			if err != nil {
				// Manejo de errores si l1Cost.Error falla (aunque con l1SegmentCost no debería)
				return nil, fmt.Errorf("Pelt (L1): error calculating segment cost for [%d, %d): %w", prevBreakpoint, currentEnd, err)
			}

			// Actualiza el valor de poda para el 'prevBreakpoint'
			pruningValues[prevBreakpoint] = minCostsToEnd[prevBreakpoint] + segmentCost

			// Calcula el costo total si 'prevBreakpoint' fuera el último punto de cambio
			totalCostIfPrev := pruningValues[prevBreakpoint] + penalty

			// Actualiza el costo mínimo y el camino si encontramos uno mejor
			minCostsToEnd[currentEnd] = totalCostIfPrev
			optimalPrevBreakpoints[currentEnd] = prevBreakpoint
		}

		// Bucle para el resto de candidatos 's' (posibles puntos de cambio anteriores)
		// Recorre desde el siguiente candidato válido hasta el último punto de inicio posible
		// para un segmento de longitud `MinSize` que termina en `currentEnd`.
		for prevBreakpoint := firstValidCandidate + 1; prevBreakpoint <= currentEnd-p.MinSize; prevBreakpoint++ {
			// Calcula el costo L1 para el segmento [prevBreakpoint, currentEnd)
			segmentCost, err := l1Cost.Error(prevBreakpoint, currentEnd)
			if err != nil {
				// Manejo de errores
				return nil, fmt.Errorf("Pelt (L1): error calculating segment cost for [%d, %d): %w", prevBreakpoint, currentEnd, err)
			}

			// Actualiza el valor de poda para el 'prevBreakpoint'
			pruningValues[prevBreakpoint] = minCostsToEnd[prevBreakpoint] + segmentCost

			// Calcula el costo total si 'prevBreakpoint' fuera el último punto de cambio
			totalCostIfPrev := pruningValues[prevBreakpoint] + penalty

			// Si este costo total es menor que el mínimo encontrado hasta ahora para 'currentEnd', actualizamos.
			if totalCostIfPrev < minCostsToEnd[currentEnd] {
				minCostsToEnd[currentEnd] = totalCostIfPrev
				optimalPrevBreakpoints[currentEnd] = prevBreakpoint
			}
		}

		// --- Lógica de Poda (Pruning) ---
		// Avanza 'firstValidCandidate' mientras sus valores podados sean peores
		// que el costo mínimo actual para 'currentEnd'.
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
	currentIdx := numSamples
	for currentIdx != 0 {
		currentIdx = optimalPrevBreakpoints[currentIdx]
		if currentIdx != 0 { // No añadimos el punto inicial (0) como un punto de cambio real
			changePoints = append(changePoints, currentIdx)
		}
	}
	sort.Ints(changePoints) // Ordena los puntos de cambio de forma ascendente

	return changePoints, nil
}

// l1SegmentCost calcula el costo L1 para un segmento univariado [startIdx, endIdx).
// Usa ordenamiento para encontrar la mediana y calcular la suma de desviaciones absolutas.
func l1SegmentCost(y []float64, startIdx, endIdx int) float64 {
	// Extrae el subsegmento. Una copia es necesaria porque el slice original
	// 'y' no debe ser modificado por el ordenamiento.
	segment := make([]float64, endIdx-startIdx)
	copy(segment, y[startIdx:endIdx])

	sort.Float64s(segment) // Ordena el segmento para encontrar la mediana

	numElements := len(segment)
	if numElements == 0 {
		return 0.0 // Un segmento vacío tiene costo cero
	}

	medianValue := 0.0
	if numElements%2 == 0 {
		// Mediana para un número par de elementos: promedio de los dos del medio
		medianValue = (segment[numElements/2-1] + segment[numElements/2]) / 2
	} else {
		// Mediana para un número impar de elementos: el elemento del medio
		medianValue = segment[numElements/2]
	}

	sumOfAbsoluteDeviations := 0.0
	for _, val := range segment {
		sumOfAbsoluteDeviations += math.Abs(val - medianValue) // Suma de los valores absolutos de las diferencias
	}
	return sumOfAbsoluteDeviations
}
