// core/kernels/gaussian.go

package kernels

import (
	"errors"
	"fmt"
	"math"

	"github.com/theDataFlowClub/ruptures/core/linalg"
	"github.com/theDataFlowClub/ruptures/core/types"
)

type GaussianKernel struct {
	Gamma float64
}

// NewGaussianKernel crea una nueva instancia de GaussianKernel.
// Si gamma es nil, la heurística de la mediana generalmente se calcula
// en el contexto del ajuste (Fit) de un modelo que usa este kernel,
// no en la inicialización del kernel en sí. Aquí se espera un valor explícito.
func NewGaussianKernel(gamma float64) *GaussianKernel {
	return &GaussianKernel{
		Gamma: gamma,
	}
}

func (gk *GaussianKernel) Compute(x, y types.Vector) (float64, error) {
	if len(x) != len(y) {
		return 0, errors.New("gaussian kernel: input vectors must have the same dimension")
	}

	// --- CORRECCIÓN AQUÍ ---
	// Captura el valor y el error de linalg.SquaredEuclideanDistance
	squaredDistance, err := linalg.SquaredEuclideanDistance(x, y)
	if err != nil {
		return 0, fmt.Errorf("gaussian kernel: failed to compute squared Euclidean distance: %w", err)
	}
	// --- FIN CORRECCIÓN ---

	// Replicar el clipping del código C: exp(-clip(gamma * squared_distance, 0.01, 100))
	valToExp := gk.Gamma * squaredDistance
	// Manual clip, o ajustar linalg.ClipFloat si lo creas
	clippedVal := math.Max(0.01, math.Min(valToExp, 100.0))

	return math.Exp(-clippedVal), nil
}

func (gk *GaussianKernel) Name() string {
	return "gaussian"
}
