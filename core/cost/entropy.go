package cost

import (
	"errors"
	"fmt"
	"math"

	"github.com/theDataFlowClub/ruptures/core/base"
	"github.com/theDataFlowClub/ruptures/core/types"
)

// maxDiscreteValue define el tamaño máximo del alfabeto para el histograma.
const maxDiscreteValue = 256

// CostEntropy implementa base.CostFunction para el costo basado en entropía de Shannon.
type CostEntropy struct {
	signalData       types.Matrix
	prefixHistograms [][]int
}

// NewCostEntropy crea una nueva instancia de CostEntropy.
func NewCostEntropy() *CostEntropy {
	return &CostEntropy{}
}

// Fit prepara la función de costo con la señal y precalcula los histogramas de prefijo.
func (c *CostEntropy) Fit(signal types.Matrix) error {
	if signal == nil || len(signal) == 0 {
		return errors.New("CostEntropy: signal cannot be nil or empty")
	}
	if len(signal[0]) != 1 {
		return errors.New("CostEntropy: requires univariate signal (e.g., []float64{val}) representing discrete values")
	}

	c.signalData = signal

	numSamples := len(signal)
	// Inicializamos prefixHistograms. `numSamples+1` para incluir un histograma "vacío" en el índice 0.
	c.prefixHistograms = make([][]int, numSamples+1)

	// El histograma en el índice 0 es todo ceros (representa el prefijo antes de cualquier dato).
	c.prefixHistograms[0] = make([]int, maxDiscreteValue)

	// Llenamos prefixHistograms iterativamente.
	// currentHist acumula los conteos del segmento actual [0, i).
	currentHist := make([]int, maxDiscreteValue)

	for i := 0; i < numSamples; i++ {
		// Creamos el slice interno para c.prefixHistograms[i+1] ANTES de copiar en él
		c.prefixHistograms[i+1] = make([]int, maxDiscreteValue) // <-- ¡¡ESTO ES CLAVE!!

		// Copiamos el estado actual de currentHist (que representa el prefijo hasta i-1)
		copy(c.prefixHistograms[i+1], currentHist)

		// Obtenemos el valor de la muestra actual
		val := int(signal[i][0])

		// Validación de rango para asegurar que el valor entra en el histograma.
		if val < 0 || val >= maxDiscreteValue {
			return fmt.Errorf("CostEntropy: value %f at index %d out of expected discrete range [0, %d)", signal[i][0], i, maxDiscreteValue-1)
		}

		// Incrementamos el conteo para el valor actual en el histograma que se almacenará
		// en la posición i+1 (representa el prefijo hasta i, es decir, [0, i+1)).
		c.prefixHistograms[i+1][val]++

		// Actualizamos currentHist para la próxima iteración (es el histograma que se está construyendo)
		currentHist[val]++
	}
	return nil
}

// Error calcula el costo de entropía de Shannon para un segmento [start, end)
func (c *CostEntropy) Error(start, end int) (float64, error) {
	if c.prefixHistograms == nil || len(c.prefixHistograms) == 0 {
		return 0, errors.New("CostEntropy: Fit() must be called before Error()")
	}
	if start < 0 || end > len(c.signalData) || start >= end {
		return 0, fmt.Errorf("CostEntropy: invalid segment [%d, %d)", start, end)
	}

	segmentLength := float64(end - start)
	if segmentLength == 0 {
		return 0.0, nil
	}

	entropy := 0.0

	for val := 0; val < maxDiscreteValue; val++ {
		count := c.prefixHistograms[end][val] - c.prefixHistograms[start][val]
		if count > 0 {
			p := float64(count) / segmentLength
			entropy -= p * math.Log2(p)
		}
	}

	return segmentLength * entropy, nil
}

// Model devuelve el nombre del modelo de costo.
func (c *CostEntropy) Model() string {
	return "entropy"
}

// GetKernel es un método stub, no aplicable para CostEntropy.
func (c *CostEntropy) GetKernel() (interface{}, error) {
	return nil, errors.New("CostEntropy does not support GetKernel")
}

// init function is called automatically when the package is initialized.
func init() {
	RegisterCostFunction("entropy", func() base.CostFunction {
		return NewCostEntropy() // <-- ¡AHORA ESTO ES CORRECTO!
	})
}
