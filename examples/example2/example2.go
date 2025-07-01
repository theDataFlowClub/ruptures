package main

import (
	"fmt"
	"log"
	"os" // Para manejar argumentos de línea de comandos si queremos flexibilidad de penalización

	// Para convertir strings a float64
	"github.com/theDataFlowClub/ruptures/core/base"
	"github.com/theDataFlowClub/ruptures/core/cmdutils"
	"github.com/theDataFlowClub/ruptures/core/cost"
	"github.com/theDataFlowClub/ruptures/core/detection/pelt"
	"github.com/theDataFlowClub/ruptures/core/types"
)

// createSignal es una función de ayuda para convertir un slice de float64 en types.Matrix.
func createSignal(data []float64, dims int) types.Matrix {
	signal := make(types.Matrix, len(data)/dims)
	for i := 0; i < len(data)/dims; i++ {
		signal[i] = make([]float64, dims)
		copy(signal[i], data[i*dims:(i+1)*dims])
	}
	return signal
}

func main() {
	log.Println("Iniciando la traducción del ejercicio de Python a Go...")

	// --- 1. Definir la señal de entrada (equivalente a np.array en Python) ---
	// La señal es unidimensional (dims=1).
	signalData := []float64{
		58.6443735, 58.33650448, 58.7098672, 58.89852615, 58.82356058,
		58.69664464, 58.49310189, 58.43034263, 58.03767246, 58.2923039,
		58.7524563, 58.87493025, 57.79134702, 57.44300875, 57.35328436,
		57.37789301, 56.54201116, 56.08067459, 55.52050122, 56.16019169,
		56.32794111, 56.62233421,
	}
	signal := createSignal(signalData, 1) // Signal unidimensional

	fmt.Printf("Señal de entrada (longitud %d): %v\n", len(signal), signal)

	// --- 2. Crear un objeto Detector (equivalente a rpt.Pelt().fit(signal)) ---
	// Por defecto, ruptures.py usa "l2" para Pelt, pero permite "rbf".
	// Si no especificamos el modelo en Python, usa "l2" por defecto.
	// Tu última implementación de PELT usó RBF en el ejemplo, así que usaremos RBF aquí por consistencia.
	// Si deseas L2, asegúrate de que CostL2 esté implementada y registrada.

	// Vamos a permitir que se pase el modelo y la penalización como argumentos de línea de comandos.
	costFuncName := "rbf" // Modelo por defecto para el ejemplo
	penalty := 5.0        // Penalización por defecto (como en tu script Python)

	// --- 2. Obtener parámetros de los argumentos de línea de comandos ---
	params := cmdutils.ParseArgs(os.Args) // ¡Llama a nuestra nueva función!

	// Obtener la función de costo basada en el nombre
	// --- 3. Crear un objeto Detector ---
	var selectedCostFunc base.CostFunction
	var err error // Declara la variable para la interfaz

	// ¡CAMBIO CLAVE AQUÍ! Usar cost.NewCost
	selectedCostFunc, err = cost.NewCost(costFuncName) // <-- ESTO ES LO QUE NECESITAS
	if err != nil {
		log.Fatalf("Error al obtener la función de costo '%s': %v. Asegúrate de que esté implementada y registrada.", params.CostFuncName, err)
	}

	// Si el modelo es RBF, podemos establecer el gamma si es necesario (aquí lo dejamos a nil para la heurística por defecto)
	if rbfCost, ok := selectedCostFunc.(*cost.CostRbf); ok {
		var gamma *float64 // nil para usar la heurística de CostRbf
		rbfCost.Gamma = gamma
		fmt.Println("Usando CostRbf con gamma heurístico.")
	} else {
		fmt.Printf("Usando función de costo: %s\n", selectedCostFunc.Model())
	}

	// Configuración de Pelt. Usamos min_size=10 y jump=1 como en tu comentario original de Python,
	// o los valores por defecto de tu Pelt Go si no se especifican.
	// Si quieres replicar `algo = rpt.Pelt().fit(signal)` exactamente, usaríamos los defaults de ruptures.
	// Los defaults de ruptures para min_size son 1, y jump es 1.
	minSize := 10 // Replicando el min_size comentado en Python
	jump := 1     // Replicando el jump comentado en Python

	// Si quieres usar los defaults de tu NewPelt sin especificar min_size o jump en la creación,
	// necesitarías un constructor de Pelt que no los pida, o pasar los mismos defaults de Go.
	// Por ahora, usamos los valores explícitos para replicar tu código Python.
	peltDetector := pelt.NewPelt(selectedCostFunc, minSize, jump)

	// Ajustar el detector a la señal (Fit)
	err = peltDetector.Fit(signal)
	if err != nil {
		log.Fatalf("Error al ajustar el detector PELT: %v", err)
	}

	fmt.Printf("Detector PELT ajustado con modelo '%s', min_size=%d, jump=%d.\n", costFuncName, minSize, jump)

	// --- 3. Predecir los breakpoints (equivalente a algo.predict(pen=5)) ---
	fmt.Printf("Prediciendo puntos de cambio con penalización (pen): %.2f\n", penalty)
	result, err := peltDetector.Predict(penalty)
	if err != nil {
		log.Fatalf("Error al predecir los puntos de cambio: %v", err)
	}

	// --- 4. Imprimir los resultados (equivalente a print(result)) ---
	// La librería ruptures en Python retorna los índices del final de cada segmento,
	// incluyendo el índice final de la señal.
	fmt.Println("\n--- Puntos de cambio detectados (índices de fin de segmento) ---")
	fmt.Printf("Go Output: %v\n", result)
	fmt.Println("--------------------------------------------------")

	// Si quieres comparar con el resultado de Python, ejecuta tu script Python y compara el array.
	// En Python: [13 22] para este caso con pen=5 y modelo "l2".
	// Para replicar "l2", necesitas implementar CostL2 y usar `go run main.go l2 5.0`
}
