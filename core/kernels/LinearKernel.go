// core/kernels/LinearKernel.go

package kernels

import (
	"errors"
	"fmt"

	"github.com/theDataFlowClub/ruptures/core/linalg" // Para Dot
	"github.com/theDataFlowClub/ruptures/core/types"
)

type LinearKernel struct{}

func NewLinearKernel() *LinearKernel {
	return &LinearKernel{}
}

func (lk *LinearKernel) Compute(x, y types.Vector) (float64, error) {
	if len(x) != len(y) {
		return 0, errors.New("linear kernel: input vectors must have the same dimension")
	}
	// Aquí se espera que linalg.Dot devuelva (float64, error)
	dot, err := linalg.Dot(x, y) // <--- Esta línea ya maneja el error
	if err != nil {
		return 0, fmt.Errorf("linear kernel: dot product failed: %w", err)
	}
	return dot, nil
}

func (lk *LinearKernel) Name() string {
	return "linear"
}
