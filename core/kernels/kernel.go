// core/kernels/kernel.go
package kernels

import (
	"errors"
	"fmt"

	"github.com/theDataFlowClub/ruptures/core/types"
) // Asumiendo que types.Vector es []float64

// Kernel es la interfaz que todos los tipos de kernel deben implementar.
type Kernel interface {
	Compute(x, y types.Vector) (float64, error) // Calcula el valor del kernel entre dos vectores
	Name() string                               // Retorna el nombre del kernel (ej. "linear", "gaussian")
	// Podríamos añadir un método para obtener parámetros si fuera necesario, ej.
	// Parameters() map[string]interface{}
}

// core/kernels/kernel.go (Solo la función NewKernelByName modificada)

// Helper para crear un kernel por nombre (similar a kernel_value_by_name en C)
// Ahora incluye un mapa de opciones para parámetros específicos del kernel.
func NewKernelByName(name string, opts map[string]float64) (Kernel, error) {
	switch name {
	case "linear":
		return NewLinearKernel(), nil
	case "gaussian":
		gamma, ok := opts["gamma"]
		if !ok {
			return nil, errors.New("gaussian kernel requires a 'gamma' parameter")
		}
		return NewGaussianKernel(gamma), nil
	case "cosine":
		return NewCosineKernel(), nil
	case "polynomial": // ¡NUEVA ENTRADA!
		scale, okS := opts["scale"]
		bias, okB := opts["bias"]
		degree, okD := opts["degree"]
		if !okS || !okB || !okD {
			return nil, errors.New("polynomial kernel requires 'scale', 'bias', and 'degree' parameters")
		}
		return NewPolynomialKernel(scale, bias, degree), nil
	default:
		return nil, fmt.Errorf("unknown kernel name: %s", name)
	}
}
