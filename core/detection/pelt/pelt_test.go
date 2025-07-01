package pelt_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/theDataFlowClub/ruptures/core/cost"           // Para CostRbf
	"github.com/theDataFlowClub/ruptures/core/detection/pelt" // Tu implementación de PELT
	"github.com/theDataFlowClub/ruptures/core/types"          // Para types.Matrix
)

// Helper para crear una señal de prueba simple
func createSignal(data []float64, dims int) types.Matrix {
	signal := make(types.Matrix, len(data)/dims)
	for i := 0; i < len(data)/dims; i++ {
		signal[i] = make([]float64, dims)
		copy(signal[i], data[i*dims:(i+1)*dims])
	}
	return signal
}

func TestPeltBasicDetection(t *testing.T) {
	// --- Test 1: Señal constante (sin puntos de cambio esperados) ---
	t.Run("ConstantSignal", func(t *testing.T) {
		signalData := []float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0}
		signal := createSignal(signalData, 1) // Señal unidimensional
		gamma := 0.1                          // Un gamma pequeño para que el kernel se parezca a un promedio
		cRbf := cost.NewCostRbf(&gamma)

		p := pelt.NewPelt(cRbf, 1, 1) // min_size=1, jump=1 (sin subsampling)
		err := p.Fit(signal)
		if err != nil {
			t.Fatalf("Fit failed: %v", err)
		}

		penalty := 1.0 // Una penalización lo suficientemente alta
		bkps, err := p.Predict(penalty)
		if err != nil {
			t.Fatalf("Predict failed: %v", err)
		}

		expectedBkps := []int{10} // Solo el final de la señal, sin cambios
		if !reflect.DeepEqual(bkps, expectedBkps) {
			t.Errorf("For constant signal, expected %v, got %v", expectedBkps, bkps)
		}
	})

	// --- Test 2: Señal con un punto de cambio claro (escalón) ---
	t.Run("SingleChangePoint", func(t *testing.T) {
		// Signal: [0,0,0,0,10,10,10,10]
		signalData := []float64{0.0, 0.0, 0.0, 0.0, 10.0, 10.0, 10.0, 10.0}
		signal := createSignal(signalData, 1)
		gamma := 0.1
		cRbf := cost.NewCostRbf(&gamma)

		p := pelt.NewPelt(cRbf, 2, 1) // min_size=2 para segmentos
		err := p.Fit(signal)
		if err != nil {
			t.Fatalf("Fit failed: %v", err)
		}

		penalty := 2.0 // Ajusta la penalización. Un valor medio debería detectar el cambio.
		bkps, err := p.Predict(penalty)
		if err != nil {
			t.Fatalf("Predict failed: %v", err)
		}

		// Ruptures en Python con signal = np.array([0,0,0,0,10,10,10,10]), model="rbf", pen=2
		// bkps = [4, 8]
		expectedBkps := []int{4, 8} // El cambio ocurre en el índice 4, el final en el 8
		sort.Ints(bkps)             // Asegurar que estén ordenados para la comparación
		if !reflect.DeepEqual(bkps, expectedBkps) {
			t.Errorf("For single change point, expected %v, got %v", expectedBkps, bkps)
		}
	})

	// --- Test 3: Señal con múltiples puntos de cambio (ejemplo más complejo) ---
	t.Run("MultipleChangePoints", func(t *testing.T) {
		// Ejemplo inspirado en ruptures: np.array([0]*10 + [5]*10 + [0]*10)
		signalData := make([]float64, 30)
		for i := 0; i < 10; i++ {
			signalData[i] = 0.0
		}
		for i := 10; i < 20; i++ {
			signalData[i] = 5.0
		}
		for i := 20; i < 30; i++ {
			signalData[i] = 0.0
		}
		signal := createSignal(signalData, 1)
		gamma := 0.1
		cRbf := cost.NewCostRbf(&gamma)

		p := pelt.NewPelt(cRbf, 1, 1)
		err := p.Fit(signal)
		if err != nil {
			t.Fatalf("Fit failed: %v", err)
		}

		penalty := 1.5 // Un valor de penalización que debería detectar 2 cambios
		bkps, err := p.Predict(penalty)
		if err != nil {
			t.Fatalf("Predict failed: %v", err)
		}

		// Ruptures en Python con señal similar, model="rbf", pen=1.5
		// expected: [10, 20, 30]
		expectedBkps := []int{10, 20, 30}
		sort.Ints(bkps)
		if !reflect.DeepEqual(bkps, expectedBkps) {
			t.Errorf("For multiple change points, expected %v, got %v", expectedBkps, bkps)
		}
	})

	// --- Test 4: Manejo de errores (señal no ajustada, penalización inválida) ---
	t.Run("ErrorHandling", func(t *testing.T) {
		cRbf := cost.NewCostRbf(nil)
		p := pelt.NewPelt(cRbf, 1, 1)

		// Test Predict sin Fit previo
		_, err := p.Predict(1.0)
		if err == nil {
			t.Error("Predict should fail if detector not fitted")
		}
		if err != nil && err.Error() != "Pelt: detector not fitted. Call Fit() first." {
			t.Errorf("Expected 'detector not fitted' error, got: %v", err)
		}

		// Test Predict con penalización inválida
		signalData := []float64{1, 2, 3, 4, 5}
		signal := createSignal(signalData, 1)
		p.Fit(signal)           // Ajustar primero para que solo falle la penalización
		_, err = p.Predict(0.0) // Penalización <= 0
		if err == nil {
			t.Error("Predict should fail with non-positive penalty")
		}
		if err != nil && err.Error() != "Pelt: penalty must be greater than 0." {
			t.Errorf("Expected 'penalty must be greater than 0' error, got: %v", err)
		}
	})

	// Puedes añadir más tests aquí para diferentes dimensiones, min_size, etc.
}
