// core/kernels/cosine.go
package kernels

import (
	"errors"
	"fmt" // Necesario para fmt.Errorf

	"github.com/theDataFlowClub/ruptures/core/linalg"
	"github.com/theDataFlowClub/ruptures/core/types"
)

type CosineKernel struct{}

func NewCosineKernel() *CosineKernel {
	return &CosineKernel{}
}

func (ck *CosineKernel) Compute(x, y types.Vector) (float64, error) {
	if len(x) != len(y) {
		return 0, errors.New("cosine kernel: input vectors must have the same dimension")
	}

	// --- CORRECCIÓN AQUÍ ---
	// Captura el valor y el error de linalg.Dot
	dotProduct, err := linalg.Dot(x, y)
	if err != nil {
		return 0, fmt.Errorf("cosine kernel: failed to compute dot product: %w", err)
	}

	// Captura el valor y el error de linalg.VectorNorm
	normX, err := linalg.VectorNorm(x)
	if err != nil {
		return 0, fmt.Errorf("cosine kernel: failed to compute norm for vector x: %w", err)
	}

	// Captura el valor y el error de linalg.VectorNorm
	normY, err := linalg.VectorNorm(y)
	if err != nil {
		return 0, fmt.Errorf("cosine kernel: failed to compute norm for vector y: %w", err)
	}
	// --- FIN CORRECCIÓN ---

	denom := normX * normY
	if denom == 0 {
		// Si uno o ambos vectores son el vector cero, la similitud del coseno es indefinida.
		// Retornamos 0.0 o un error, dependiendo de la convención deseada.
		// ruptures parece retornar 0.0 en este caso.
		return 0.0, nil
	}

	return dotProduct / denom, nil
}

func (ck *CosineKernel) Name() string {
	return "cosine"
}
