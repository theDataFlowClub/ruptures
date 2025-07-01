package cmdutils

import (
	"fmt"
	"strconv"
)

// ParametersFromArgs es una estructura para almacenar los parámetros parseados.
type ParametersFromArgs struct {
	CostFuncName string
	Penalty      float64
}

// parseArgs analiza los argumentos de línea de comandos y devuelve los parámetros.
// Asume que los argumentos son: [nombre_programa] [cost_func_name] [penalty_value]
func ParseArgs(args []string) ParametersFromArgs {
	params := ParametersFromArgs{
		CostFuncName: "rbf", // Valor por defecto
		Penalty:      5.0,   // Valor por defecto
	}

	if len(args) > 1 {
		params.CostFuncName = args[1]
	}
	if len(args) > 2 {
		p, err := strconv.ParseFloat(args[2], 64)
		if err == nil {
			params.Penalty = p
		} else {
			fmt.Printf("Advertencia: No se pudo parsear la penalización '%s', usando valor por defecto: %.2f\n", args[2], params.Penalty)
		}
	}
	return params
}
