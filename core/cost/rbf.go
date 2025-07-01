package cost

import (
	"errors"
	"fmt"
	"math"

	"github.com/theDataFlowClub/ruptures/core/base"
	"github.com/theDataFlowClub/ruptures/core/exceptions"
	"github.com/theDataFlowClub/ruptures/core/kernels" // ¡NUEVO IMPORT!
	"github.com/theDataFlowClub/ruptures/core/linalg"
	"github.com/theDataFlowClub/ruptures/core/stat"
	"github.com/theDataFlowClub/ruptures/core/types"
)

// CostRbf implements the kernel cost function using an RBF kernel.

// CostRbf implements the kernel cost function using an RBF kernel.
type CostRbf struct {
	Signal         types.Matrix // The fitted signal (n_samples, n_features)
	minSegmentSize int          // Minimum segment size, default 1 // <-- CAMBIADO AQUÍ

	// Gamma for the RBF kernel. Use a pointer to allow nil for default heuristic calculation.
	Gamma *float64

	// _gram is the cached Gram matrix for the entire signal.
	// It's initialized to nil and computed on first access via GetGram().
	_gram types.Matrix

	// _kernel stores the actual RBF kernel instance used for computations.
	// This is what Pelt needs to call Compute().
	_kernel kernels.Kernel
}

// MinSize returns the minimum required length of a segment for this cost function.
// This prevents invalid segmentation where the mathematical assumptions (e.g. variance)
// or kernel requirements (e.g. Gram matrix rank) would fail.
func (c *CostRbf) MinSize() int {
	return c.minSegmentSize // ¡Aquí sí, con el receptor 'c.' para referirse al campo!
}

// NewCostRbf creates and returns a new instance of CostRbf.
// If gamma is nil, a default heuristic will be calculated during Fit() or the first access to GetGram().
func NewCostRbf(gamma *float64) *CostRbf {
	return &CostRbf{
		minSegmentSize: 1, // <-- CAMBIADO AQUÍ
		Gamma:          gamma,
		_gram:          nil,
		_kernel:        nil,
	}
}

// Model returns the name of the cost function model.
func (c *CostRbf) Model() string {
	return "rbf"
}

// GetKernel returns the underlying kernels.Kernel instance used by CostRbf.
// This is the method Pelt will call.
func (c *CostRbf) GetKernel() (kernels.Kernel, error) {
	if c._kernel == nil {
		_, err := c.GetGram()
		if err != nil {
			return nil, fmt.Errorf("CostRbf: failed to initialize kernel via GetGram(): %w", err)
		}
		if c.Gamma == nil {
			return nil, errors.New("CostRbf: gamma not set after GetGram() call, cannot create kernel")
		}
		c._kernel = kernels.NewGaussianKernel(*c.Gamma)
	}
	return c._kernel, nil
}

// GetGram calculates and returns the Gram matrix (lazy loading).
// This method is the Go equivalent of the Python @property def gram.
func (c *CostRbf) GetGram() (types.Matrix, error) {
	if c.Signal == nil || len(c.Signal) == 0 {
		return nil, errors.New("CostRbf: signal not fitted or empty, cannot compute Gram matrix. Call Fit() first")
	}

	if c._gram == nil {
		distSq, err := linalg.PdistSqEuclidean(c.Signal)
		if err != nil {
			return nil, fmt.Errorf("CostRbf: failed to calculate pairwise squared Euclidean distances: %w", err)
		}

		if c.Gamma == nil {
			if len(c.Signal) < 2 {
				gammaVal := 1.0
				c.Gamma = &gammaVal
			} else {
				medianDistSq, err := stat.Median(distSq)
				if err != nil {
					return nil, fmt.Errorf("CostRbf: error calculating median for gamma heuristic: %w", err)
				}
				if medianDistSq != 0 {
					gammaVal := 1.0 / medianDistSq
					c.Gamma = &gammaVal
				} else {
					gammaVal := 1.0
					c.Gamma = &gammaVal
				}
			}
		}

		if c._kernel == nil {
			c._kernel = kernels.NewGaussianKernel(*c.Gamma)
		}

		for i := range distSq {
			distSq[i] *= *c.Gamma
		}

		clippedDistSq := linalg.ClipSlice(distSq, 1e-2, 1e2)

		nSamples := len(c.Signal)
		tempGram, err := linalg.Squareform(clippedDistSq, nSamples)
		if err != nil {
			return nil, fmt.Errorf("CostRbf: failed to convert to square form: %w", err)
		}

		for i := range tempGram {
			for j := range tempGram[i] {
				tempGram[i][j] = math.Exp(-tempGram[i][j])
			}
		}

		c._gram = tempGram
	}
	return c._gram, nil
}

// Fit sets parameters of the instance and optionally calculates default gamma.
func (c *CostRbf) Fit(signal types.Matrix) error {
	if signal == nil || len(signal) == 0 || (len(signal) > 0 && len(signal[0]) == 0) {
		return exceptions.ErrNotEnoughPoints
	}
	c.Signal = signal

	if c.Gamma == nil {
		_, err := c.GetGram()
		if err != nil {
			return fmt.Errorf("CostRbf: failed to calculate default gamma during fit: %w", err)
		}
	} else {
		if c._kernel == nil {
			c._kernel = kernels.NewGaussianKernel(*c.Gamma)
		}
	}
	return nil
}

// Error calculates the RBF kernel cost on the segment [start:end].
func (c *CostRbf) Error(start, end int) (float64, error) {
	if c.Signal == nil {
		return 0.0, errors.New("CostRbf: signal not fitted, call Fit() first")
	}

	segmentLen := end - start
	if segmentLen < c.minSegmentSize { // <-- CAMBIADO AQUÍ
		return 0.0, exceptions.ErrNotEnoughPoints
	}

	fullGram, err := c.GetGram()
	if err != nil {
		return 0.0, fmt.Errorf("CostRbf: failed to get Gram matrix: %w", err)
	}

	subGram := make(types.Matrix, segmentLen)
	for i := 0; i < segmentLen; i++ {
		if start+i >= len(fullGram) || end > len(fullGram[0]) {
			return 0.0, errors.New("CostRbf: sub-segment indices out of bounds for Gram matrix")
		}
		subGram[i] = make([]float64, segmentLen)
		copy(subGram[i], fullGram[start+i][start:end])
	}

	diagSum, err := linalg.DiagonalSum(subGram)
	if err != nil {
		return 0.0, fmt.Errorf("CostRbf: error calculating diagonal sum: %w", err)
	}
	totalSum, err := linalg.Sum(subGram)
	if err != nil {
		return 0.0, fmt.Errorf("CostRbf: error calculating total sum: %w", err)
	}

	val := diagSum - (totalSum / float64(segmentLen))
	return val, nil
}

// GetCachedGramForTest returns the cached Gram matrix.
func (c *CostRbf) GetCachedGramForTest() types.Matrix {
	return c._gram
}

// SetCachedGramForTest sets the cached Gram matrix.
func (c *CostRbf) SetCachedGramForTest(gram types.Matrix) {
	c._gram = gram
}

// init function is called automatically when the package is initialized.
func init() {
	RegisterCostFunction("rbf", func() base.CostFunction {
		return NewCostRbf(nil)
	})
}
