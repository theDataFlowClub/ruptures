package pelt

import (
	"fmt"
	"math"
	"sort"

	"github.com/theDataFlowClub/ruptures/core/cost"
)

// predictRbfOptimized es la implementación de PELT optimizada para CostRbf.
// Esta es tu función `Predict` original, renombrada.
func (p *Pelt) predictRbfOptimized(rbfCost *cost.CostRbf, penalty float64) ([]int, error) {
	// --- Obtener el kernel de la CostFunctionRbf ---
	// Esta sección es específica de RBF y está bien aquí.
	currentKernel, err := rbfCost.GetKernel()
	if err != nil {
		return nil, fmt.Errorf("Pelt (RBF): failed to get kernel from CostRbf: %w", err)
	}
	// --- FIN: Obtener el kernel ---

	// Arrays auxiliares (específicos de RBF)
	M_V := make([]float64, p.nSamples+1)
	M_path := make([]int, p.nSamples+1)
	D := make([]float64, p.nSamples+1)
	S := make([]float64, p.nSamples+1)
	M_pruning := make([]float64, p.nSamples+1)

	// Inicialización
	for t := 0; t <= p.nSamples; t++ {
		D[t] = 0.0
		S[t] = 0.0
		M_V[t] = math.Inf(1)
		M_path[t] = 0
		M_pruning[t] = math.Inf(1)
	}
	M_V[0] = -penalty
	M_pruning[0] = 0.0

	var (
		t, s                    int
		s_min                   = 0
		c_cost, c_cost_sum, c_r float64
	)

	// Bucle inicial para t < 2 * min_size
	for t = 1; t < 2*p.MinSize && t <= p.nSamples; t++ {
		diag_element_val, err := currentKernel.Compute(p.signal[t-1], p.signal[t-1])
		if err != nil {
			return nil, fmt.Errorf("Pelt (RBF): error computing diagonal kernel element at t=%d: %w", t, err)
		}
		D[t] = D[t-1] + diag_element_val

		c_r = 0.0
		for s = t - 1; s >= 0; s-- {
			val, err := currentKernel.Compute(p.signal[s], p.signal[t-1])
			if err != nil {
				return nil, fmt.Errorf("Pelt (RBF): error computing kernel element for S at s=%d, t-1=%d: %w", s, t-1, err)
			}
			c_r += val
			S[s] += 2*c_r - diag_element_val
		}

		if t > 0 {
			c_cost = (D[t] - D[0]) - (S[0] / float64(t))
		} else {
			c_cost = 0.0
		}
		M_V[t] = c_cost + penalty
		M_path[t] = 0
	}

	// Bucle de computación principal (PELT)
	for t = 2 * p.MinSize; t <= p.nSamples; t++ {
		diag_element_val, err := currentKernel.Compute(p.signal[t-1], p.signal[t-1])
		if err != nil {
			return nil, fmt.Errorf("Pelt (RBF): error computing diagonal kernel element at t=%d in main loop: %w", t, err)
		}
		D[t] = D[t-1] + diag_element_val

		c_r = 0.0
		for s = t - 1; s >= s_min; s-- {
			val, err := currentKernel.Compute(p.signal[s], p.signal[t-1])
			if err != nil {
				return nil, fmt.Errorf("Pelt (RBF): error computing kernel element for S in main loop at s=%d, t-1=%d: %w", s, t-1, err)
			}
			c_r += val
			S[s] += 2*c_r - diag_element_val
		}

		M_V[t] = math.Inf(1)

		if s_min <= t-p.MinSize {
			segmentLen := float64(t - s_min)
			c_cost = (D[t] - D[s_min]) - (S[s_min] / segmentLen)
			c_cost_sum = M_V[s_min] + c_cost
			M_pruning[s_min] = c_cost_sum
			c_cost_sum += penalty

			M_V[t] = c_cost_sum
			M_path[t] = s_min
		}

		for s = s_min + 1; s <= t-p.MinSize; s++ {
			segmentLen := float64(t - s)
			c_cost = (D[t] - D[s]) - (S[s] / segmentLen)
			c_cost_sum = M_V[s] + c_cost
			M_pruning[s] = c_cost_sum

			c_cost_sum += penalty

			if M_V[t] > c_cost_sum {
				M_V[t] = c_cost_sum
				M_path[t] = s
			}
		}

		for (s_min < t-p.MinSize+1) && (M_pruning[s_min] >= M_V[t]) {
			if s_min == 0 {
				s_min = p.MinSize
			} else {
				s_min++
			}
		}
	}

	// Reconstruir los puntos de cambio
	resultBkps := []int{p.nSamples}
	current := p.nSamples
	for current != 0 {
		current = M_path[current]
		if current != 0 {
			resultBkps = append(resultBkps, current)
		}
	}
	sort.Ints(resultBkps)

	return resultBkps, nil
}
