# Testing

## **Estrategia de Pruebas para `base` y `utils`**

Para `base` y `utils`, la estrategia de pruebas se centrará en dos tipos principales:

1.  **Pruebas unitarias:** Para cada función o método exportado (`SanityCheck`, `Pairwise`, `SumOfCosts`, etc.), probaremos su lógica individualmente, asegurando que se comporten como se espera en diferentes escenarios (casos base, casos límite, errores).
2.  **Pruebas de interfaz (para `base`):** Nos aseguraremos de que cualquier implementación futura de las interfaces `Estimator` y `CostFunction` cumpla con los contratos definidos. Aunque no podemos probar las interfaces directamente, sí podemos preparar el terreno para validar sus implementaciones.

-----

### **Implementación de Pruebas**

En Go, los archivos de prueba se colocan en el mismo paquete que el código que prueban, con el sufijo `_test.go`.

#### **1. Pruebas para `core/utils`**

Crearás un archivo llamado `pairwise_test.go` y `sanitycheck_test.go` (o puedes agruparlos en `utils_test.go` si prefieres, aunque separarlos por función es más granular).

**`core/utils/pairwise_test.go`**

```go
package utils_test // Note: _test suffix for package name when testing external functions

import (
	"reflect" // Used for deep comparison of slices and structs
	"testing"

	"github.com/theDataFlowClub/ruptures/core/utils" // Import the package being tested
)

func TestPairwise(t *testing.T) {
	// Test case 1: Empty slice
	t.Run("EmptySlice", func(t *testing.T) {
		input := []int{}
		expected := []utils.Pair[int, int](nil) // nil slice for expected empty result
		result := utils.Pairwise(input)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Pairwise(%v) = %v; want %v", input, result, expected)
		}
	})

	// Test case 2: Single element slice
	t.Run("SingleElementSlice", func(t *testing.T) {
		input := []int{1}
		expected := []utils.Pair[int, int](nil)
		result := utils.Pairwise(input)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Pairwise(%v) = %v; want %v", input, result, expected)
		}
	})

	// Test case 3: Standard slice
	t.Run("StandardSlice", func(t *T) {
		input := []int{1, 2, 3, 4}
		expected := []utils.Pair[int, int]{
			{First: 1, Second: 2},
			{First: 2, Second: 3},
			{First: 3, Second: 4},
		}
		result := utils.Pairwise(input)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Pairwise(%v) = %v; want %v", input, result, expected)
		}
	})

	// Test case 4: Slice with duplicate elements
	t.Run("DuplicateElements", func(t *testing.T) {
		input := []int{5, 5, 6, 6}
		expected := []utils.Pair[int, int]{
			{First: 5, Second: 5},
			{First: 5, Second: 6},
			{First: 6, Second: 6},
		}
		result := utils.Pairwise(input)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Pairwise(%v) = %v; want %v", input, result, expected)
		}
	})
}

func TestUnzip(t *testing.T) {
	// Test case 1: Empty slice of pairs
	t.Run("EmptyPairs", func(t *testing.T) {
		input := []utils.Pair[int, int]{}
		expectedA := []int{}
		expectedB := []int{}
		resultA, resultB := utils.Unzip(input)
		if !reflect.DeepEqual(resultA, expectedA) || !reflect.DeepEqual(resultB, expectedB) {
			t.Errorf("Unzip(%v) = (%v, %v); want (%v, %v)", input, resultA, resultB, expectedA, expectedB)
		}
	})

	// Test case 2: Standard slice of pairs
	t.Run("StandardPairs", func(t *testing.T) {
		input := []utils.Pair[int, int]{
			{First: 1, Second: 10},
			{First: 2, Second: 20},
			{First: 3, Second: 30},
		}
		expectedA := []int{1, 2, 3}
		expectedB := []int{10, 20, 30}
		resultA, resultB := utils.Unzip(input)
		if !reflect.DeepEqual(resultA, expectedA) || !reflect.DeepEqual(resultB, expectedB) {
			t.Errorf("Unzip(%v) = (%v, %v); want (%v, %v)", input, resultA, resultB, expectedA, expectedB)
		}
	})

	// Test case 3: Pairs with duplicate values
	t.Run("DuplicateValues", func(t *testing.T) {
		input := []utils.Pair[int, int]{
			{First: 7, Second: 7},
			{First: 8, Second: 9},
		}
		expectedA := []int{7, 8}
		expectedB := []int{7, 9}
		resultA, resultB := utils.Unzip(input)
		if !reflect.DeepEqual(resultA, expectedA) || !reflect.DeepEqual(resultB, expectedB) {
			t.Errorf("Unzip(%v) = (%v, %v); want (%v, %v)", input, resultA, resultB, expectedA, expectedB)
		}
	})
}
```

**`core/utils/sanitycheck_test.go`**

```go
package utils_test

import (
	"testing"

	"github.com/theDataFlowClub/ruptures/core/utils"
)

func TestSanityCheck(t *testing.T) {
	testCases := []struct {
		name     string
		nSamples int
		nBkps    int
		jump     int
		minSize  int
		expected bool
	}{
		{
			name:     "Valid_SimpleCase",
			nSamples: 100,
			nBkps:    1,
			jump:     1,
			minSize:  10,
			expected: true, // Segment 0-10, 10-100 (90 points) is possible
		},
		{
			name:     "Valid_MultipleBkps",
			nSamples: 100,
			nBkps:    3,
			jump:     1,
			minSize:  10,
			expected: true, // 3 bkps, 4 segments. Smallest possible: 3*10 + 10 = 40 points
		},
		{
			name:     "Invalid_TooManyBkps",
			nSamples: 50,
			nBkps:    5,
			jump:     1,
			minSize:  10,
			expected: false, // 5 bkps (6 segments) * 10 minSize = 60 points > 50 samples
		},
		{
			name:     "Invalid_TooManyBkpsWithJump",
			nSamples: 100,
			nBkps:    10, // Max admissible bkps for jump 10 is 100/10 = 10
			jump:     10,
			minSize:  5,
			expected: false, // 10 bkps * ceil(5/10)*10 + 5 = 10*1*10 + 5 = 105 > 100
		},
		{
			name:     "Valid_WithJumpConstraint",
			nSamples: 100,
			nBkps:    3,
			jump:     10,
			minSize:  10,
			expected: true, // 3 bkps, 4 segments. Smallest: 3*ceil(10/10)*10 + 10 = 3*10 + 10 = 40. OK.
		},
		{
			name:     "Invalid_MinSizeTooLarge",
			nSamples: 20,
			nBkps:    1,
			jump:     1,
			minSize:  15,
			expected: false, // 1 bkp, 2 segments. Smallest: 1*15 + 15 = 30 > 20
		},
		{
			name:     "EdgeCase_ZeroBkps",
			nSamples: 10,
			nBkps:    0,
			jump:     1,
			minSize:  1,
			expected: true, // 0 bkps, 1 segment (10 points). OK.
		},
		{
			name:     "EdgeCase_nSamplesLessThanMinSize",
			nSamples: 5,
			nBkps:    0,
			jump:     1,
			minSize:  10,
			expected: false, // 0 bkps, but minSize (10) > nSamples (5)
		},
		{
			name:     "EdgeCase_MinimumPossible",
			nSamples: 2,
			nBkps:    0,
			jump:     1,
			minSize:  2,
			expected: true, // 1 segment of size 2.
		},
		{
			name:     "EdgeCase_JumpGreaterThanMinSize",
			nSamples: 100,
			nBkps:    1,
			jump:     20,
			minSize:  10,
			expected: true, // 1 bkp, 2 segments. minPoints = 1*ceil(10/20)*20 + 10 = 1*1*20 + 10 = 30. OK.
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.SanityCheck(tc.nSamples, tc.nBkps, tc.jump, tc.minSize)
			if result != tc.expected {
				t.Errorf("SanityCheck(nSamples:%d, nBkps:%d, jump:%d, minSize:%d) = %v; want %v",
					tc.nSamples, tc.nBkps, tc.jump, tc.minSize, result, tc.expected)
			}
		})
	}
}
```

-----

#### **2. Pruebas para `core/base` (`SumOfCosts`)**

Crearás un archivo `sumofcosts_test.go` dentro de `core/base`. Para probar `SumOfCosts`, necesitaremos una implementación *mock* (simulada) de la interfaz `CostFunction`, ya que `SumOfCosts` depende de ella.

**`core/base/sumofcosts_test.go`**

```go
package base_test // Note: _test suffix for package name when testing external functions

import (
	"testing"

	"github.com/theDataFlowClub/ruptures/core/base" // Import the package being tested
	"github.com/theDataFlowClub/ruptures/core/types" // Import types for mock
)

// MockCostFunction is a simple mock implementation of the base.CostFunction interface
// for testing purposes. It returns a predefined error value for any segment.
type MockCostFunction struct {
	costPerSegment float64
}

func (m *MockCostFunction) Fit(signal types.Matrix) error {
	// No-op for this mock
	return nil
}

func (m *MockCostFunction) Error(start, end int) float64 {
	// For testing, return a fixed cost per segment.
	// You could make this more complex if needed, e.g., based on segment length.
	return m.costPerSegment
}

func (m *MockCostFunction) Model() string {
	return "mock_cost"
}

func TestSumOfCosts(t *testing.T) {
	// Create a mock cost function that returns a fixed cost for each segment.
	mockCost := &MockCostFunction{costPerSegment: 10.0}

	testCases := []struct {
		name     string
		bkps     []int
		expected float64
	}{
		{
			name:     "NoBreakpoints_EmptySlice",
			bkps:     []int{},
			expected: 0.0, // As per your current implementation
		},
		{
			name:     "OneSegment_NoBkpsProvided",
			bkps:     []int{100}, // Represents a signal from 0 to 100, one segment
			expected: 10.0,       // 1 segment * 10.0 cost/segment
		},
		{
			name:     "TwoSegments",
			bkps:     []int{50, 100}, // Segments 0-50, 50-100
			expected: 20.0,           // 2 segments * 10.0 cost/segment
		},
		{
			name:     "MultipleSegments",
			bkps:     []int{25, 50, 75, 100}, // Segments 0-25, 25-50, 50-75, 75-100
			expected: 40.0,                    // 4 segments * 10.0 cost/segment
		},
		{
			name:     "BkpsNotSorted_ShouldStillWorkIfLogicHandlesIt", // Although bkps should typically be sorted
			bkps:     []int{100, 50, 75}, // This would typically not be sorted in a real scenario, but SumOfCosts doesn't sort it.
			expected: 30.0,               // It will still process pairs (0,100), (100,50), (50,75) and sum them.
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := base.SumOfCosts(mockCost, tc.bkps)
			if result != tc.expected {
				t.Errorf("SumOfCosts(mockCost, %v) = %f; want %f", tc.bkps, result, tc.expected)
			}
		})
	}
}
```

-----

### **Cómo Ejecutar las Pruebas**

1.  **Navega a la raíz de tu proyecto:** Abre tu terminal y ve a la carpeta `ruptures/`.
2.  **Ejecuta todas las pruebas:**
    ```bash
    go test ./...
    ```
    Esto buscará y ejecutará todos los archivos `_test.go` en tu proyecto y sus subdirectorios.
3.  **Ejecutar pruebas de un paquete específico:**
    ```bash
    go test ./core/utils
    go test ./core/base
    ```
4.  **Ver el resultado de cobertura de código:**
    ```bash
    go test -cover ./...
    ```
    Esto te mostrará un porcentaje de cuánto de tu código está cubierto por las pruebas, lo cual es muy útil.

-----

### **Puntos Clave para Recordar**

  * **Nomenclatura:** Las funciones de prueba siempre empiezan con `Test` seguido del nombre de la función que se va a probar (ej. `TestPairwise`, `TestSanityCheck`). Toman un argumento `*testing.T`.
  * **Subtests (`t.Run`):** Usar `t.Run` es una excelente práctica para agrupar pruebas relacionadas y darles nombres descriptivos. Esto mejora la legibilidad de la salida de las pruebas.
  * **`reflect.DeepEqual`:** Para comparar slices o structs, no uses `==`. Go proporciona `reflect.DeepEqual` para comparaciones de contenido.
  * **Mocks:** Para probar funciones que dependen de interfaces (como `SumOfCosts` que depende de `CostFunction`), crea implementaciones *mock* de esas interfaces. Esto aísla la lógica que estás probando y evita dependencias externas.

Con estas pruebas unitarias bien establecidas, se tendrá una base sólida y confiable para seguir construyendo la librería.
