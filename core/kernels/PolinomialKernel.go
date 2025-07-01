// core/kernels/polynomial.go
package kernels

import (
	"errors"
	"fmt" // Para fmt.Errorf si se usa en Compute
	"math"

	"github.com/theDataFlowClub/ruptures/core/linalg"
	"github.com/theDataFlowClub/ruptures/core/types"
)

// PolynomialKernel implementa el kernel Polinomial.
// K(x, y) = (scale * <x, y> + bias)^degree
type PolynomialKernel struct {
	Scale  float64
	Bias   float64
	Degree float64 // El grado puede ser flotante, aunque comúnmente es un entero.
}

// NewPolynomialKernel crea y retorna una nueva instancia de PolynomialKernel.
// Se pueden especificar los parámetros de escala, bias y grado.
func NewPolynomialKernel(scale, bias, degree float64) *PolynomialKernel {
	return &PolynomialKernel{
		Scale:  scale,
		Bias:   bias,
		Degree: degree,
	}
}

// Compute calcula el valor del kernel Polinomial entre dos vectores.
func (pk *PolynomialKernel) Compute(x, y types.Vector) (float64, error) {
	if len(x) != len(y) {
		return 0, errors.New("polynomial kernel: input vectors must have the same dimension")
	}

	// Calcula el producto punto <x, y>
	dotProduct, err := linalg.Dot(x, y)
	if err != nil {
		return 0, fmt.Errorf("polynomial kernel: failed to compute dot product: %w", err)
	}

	// Aplica la transformación afín: (scale * <x, y> + bias)
	baseVal := (pk.Scale * dotProduct) + pk.Bias

	// Eleva a la potencia 'degree': baseVal ^ degree
	// math.Pow (base, exponent)
	return math.Pow(baseVal, pk.Degree), nil
}

// Name retorna el nombre del kernel.
func (pk *PolynomialKernel) Name() string {
	return "polynomial"
}
