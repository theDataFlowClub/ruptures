package main

import (
	"fmt"
	"log"
	"os"

	// Necesario para convertir string a float64
	// Para la interfaz CostFunction
	"github.com/theDataFlowClub/ruptures/core/cmdutils"       // Donde pusimos la función ParseArgs
	"github.com/theDataFlowClub/ruptures/core/cost"           // Para la fábrica de funciones de costo (NewCost) y CostRbf
	"github.com/theDataFlowClub/ruptures/core/detection/pelt" // Tu implementación de PELT
	"github.com/theDataFlowClub/ruptures/core/types"          // Para types.Matrix
)

// createSignal es una función de ayuda para convertir un slice de float64 en types.Matrix.
// Esto simula cómo tus datos de señal se representarían en Go.
func createSignal(data []float64, dims int) types.Matrix {
	signal := make(types.Matrix, len(data)/dims)
	for i := 0; i < len(data)/dims; i++ {
		signal[i] = make([]float64, dims)
		copy(signal[i], data[i*dims:(i+1)*dims])
	}
	return signal
}

func main() {
	log.Println("Iniciando ejemplo de detección de puntos de cambio con PELT...")

	// --- 1. Definir la señal de entrada ---
	// Esta es una señal de ejemplo con dos cambios claros para la demostración.
	signalData := []float64{}
	for i := 0; i < 20; i++ { // Primer segmento (valor 0.5)
		signalData = append(signalData, 0.5)
	}
	for i := 0; i < 30; i++ { // Segundo segmento (valor 3.0)
		signalData = append(signalData, 3.0)
	}
	for i := 0; i < 20; i++ { // Tercer segmento (valor 1.0)
		signalData = append(signalData, 1.0)
	}
	signal := createSignal(signalData, 1) // Creamos una señal unidimensional

	fmt.Printf("Señal de entrada (primeros 10 puntos): %v...\n", signal[:10])
	fmt.Printf("Longitud total de la señal: %d\n", len(signal))

	// --- 2. Obtener parámetros de los argumentos de línea de comandos ---
	// La función ParseArgs está en el paquete cmdutils y se encarga de leer
	// el nombre del modelo y la penalización de los argumentos de la terminal.
	params := cmdutils.ParseArgs(os.Args)
	fmt.Printf("DEBUG: Parámetro de función de costo recibido: '%s'\n", params.CostFuncName) // <-- Línea de depuración 1

	// --- 3. Crear un objeto de la función de costo seleccionada ---
	// Utilizamos la fábrica NewCost del paquete 'cost' para obtener una instancia
	// de la función de costo (RBF, L1, L2) basada en el nombre que se obtuvo de los argumentos.
	selectedCostFunc, err := cost.NewCost(params.CostFuncName)
	if err != nil {
		log.Fatalf("Error al obtener la función de costo '%s': %v. Asegúrate de que esté implementada y registrada.", params.CostFuncName, err)
	}

	// Si la función de costo seleccionada es RBF, podemos configurar su parámetro Gamma.
	// La implementación de CostRbf usa una heurística por defecto si Gamma es nil.
	if rbfCost, ok := selectedCostFunc.(*cost.CostRbf); ok {
		var gamma *float64 // Lo dejamos como nil para que CostRbf use su heurística por defecto
		// Podrías añadir lógica aquí para leer 'gamma' desde los argumentos si lo necesitaras:
		// if len(os.Args) > 3 { /* leer gamma de os.Args[3] */ }
		rbfCost.Gamma = gamma
		fmt.Println("Usando CostRbf con gamma heurístico.")
	} else {
		fmt.Printf("Usando función de costo: %s\n", selectedCostFunc.Model())
	}

	// --- 4. Configurar y ajustar el detector PELT ---
	// 'minSize' es el tamaño mínimo de un segmento detectado.
	// 'jump' es el paso para el subsampling (1 significa sin subsampling, cada punto).
	minSize := 2
	jump := 1
	peltDetector := pelt.NewPelt(selectedCostFunc, minSize, jump)

	// Ajustamos el detector a nuestra señal de entrada. Esto prepara el algoritmo para la detección.
	err = peltDetector.Fit(signal)
	if err != nil {
		log.Fatalf("Error al ajustar el detector PELT a la señal: %v", err)
	}
	log.Println("Detector PELT ajustado correctamente a la señal.")

	// --- 5. Predecir los puntos de cambio ---
	// La penalización (penalty) es un parámetro crucial:
	// Un valor más alto resulta en menos puntos de cambio detectados.
	// Un valor más bajo resulta en más puntos de cambio detectados.
	fmt.Printf("Prediciendo puntos de cambio con penalización (pen): %.2f\n", params.Penalty)
	changePoints, err := peltDetector.Predict(params.Penalty)
	if err != nil {
		log.Fatalf("Error al predecir los puntos de cambio: %v", err)
	}

	// --- 6. Mostrar los resultados ---
	fmt.Println("\n--- Resultados de Detección de Puntos de Cambio ---")
	if len(changePoints) <= 1 {
		fmt.Println("No se detectaron puntos de cambio significativos para la penalización dada.")
	} else {
		// Los puntos de cambio son los índices donde termina un segmento y comienza el siguiente.
		// El último punto siempre es la longitud total de la señal.
		fmt.Printf("Puntos de cambio detectados: %v\n", changePoints)
		fmt.Println("Interpretación (índices de los segmentos):")
		prevBkp := 0
		for _, bkp := range changePoints {
			if bkp == 0 { // Asegurarse de no procesar un 0 inicial si aparece por alguna lógica.
				continue
			}
			if bkp > prevBkp { // Asegurar que el segmento sea válido (inicio < fin)
				fmt.Printf("  Segmento: [%d, %d)\n", prevBkp, bkp)
			}
			prevBkp = bkp // El fin del segmento actual es el inicio del siguiente
		}
	}
	fmt.Println("--------------------------------------------------")
}
