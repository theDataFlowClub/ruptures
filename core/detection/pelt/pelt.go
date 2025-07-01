package pelt

import (

	// Asegúrate de importar sort

	"github.com/theDataFlowClub/ruptures/core/base"       // ¡Necesario para el casting a CostRbf, CostL1, CostL2!
	"github.com/theDataFlowClub/ruptures/core/exceptions" // Necesario para el tipo kernels.Kernel (solo para RBF)
	"github.com/theDataFlowClub/ruptures/core/types"
)

type Pelt struct {
	Cost     base.CostFunction // La función de costo (ej. CostRbf, CostL1, CostL2)
	MinSize  int               // Tamaño mínimo de un segmento
	Jump     int               // Salto de subsampling (opcional, implementar más tarde si es necesario)
	nSamples int               // Número de muestras en la señal
	signal   types.Matrix      // La señal ajustada
}

// NewPelt crea una nueva instancia de Pelt.
func NewPelt(costFunc base.CostFunction, minSize int, jump int) *Pelt {
	return &Pelt{
		Cost:    costFunc,
		MinSize: minSize,
		Jump:    jump,
	}
}

// Fit establece la señal en el detector y la función de costo.
func (p *Pelt) Fit(signal types.Matrix) error {
	if signal == nil || len(signal) == 0 {
		return exceptions.ErrInvalidSignal
	}
	p.signal = signal
	p.nSamples = len(signal)
	return p.Cost.Fit(signal) // Asegura que la función de costo se ajuste a la señal
}
