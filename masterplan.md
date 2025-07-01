# GOrupture plan

## 📂 Análisis de la estructura de `ruptures`

### Analisis de Módulos en `ruptures`:

| Módulo / Carpeta | Propósito                                                           |
| ---------------- | ------------------------------------------------------------------- |
| `base.py`        | Define interfaces y clases abstractas (`BaseCost`, `BaseEstimator`) |
| `cost`           | Funciones de costo como L1, L2, RBF...                              |
| `detection`      | Algoritmos como Pelt, Binseg, BottomUp, etc.                        |
| `datasets`       | Utilidades para generar señales de prueba.                          |
| `metrics`        | Evaluación (ej. precisión de detección).                            |
| `show`           | Renderizado (`matplotlib`, etc).                                    |
| `exceptions.py`  | Excepciones custom.                                                 |
| `utils`          | Funciones auxiliares como `sanity_check`                            |
| `__init__.py`    | Organiza qué se expone como API pública.                            |

---

## 📌 Objetivo: Organizar Reimplementación

Si empezamos desde las **técnicas de detección (`detection`) como paquetes separados** podremos tener unidades independientes que prueben conceptos concretos, Detectar dependencias comunes (como `BaseCost`, `sanity_check`, etc.) y Modularizar de abajo hacia arriba.

---

## ✅ Estructura para Go

Estructura propuesta en Go, inspirada en la modularidad de `ruptures`:

```text
rupture/
├── cmd/                  # CLI para probar, correr, exportar resultados
├── core/
│   ├── base/             # Interfaces clave: CostFunction, Estimator
│   ├── cmdutils/         # 
│   ├── cost/             # Costos L1, L2, RBF... testables de forma independiente
│   ├── datasets/         # Simulación de señales (¡útil para tests!)
│   ├── detection/        # Algoritmos separados por carpeta = 💯
│   └── exceptions/       # `ErrInvalidSegment`, `ErrIncompatibleCost`, etc.
│   └── kernels/          # Base para libreria autonoma: pykernels port
│   └── linalg/           # utilerias de Gonum / Algebra lineal : shotcuts, envoltorios
│   └── stat/             # utilerias de Gonum / Stats: shotcuts, envoltorio + scypy portss
│   ├── metrics/          # Precision, coverage, F1, etc.
│   ├── types/            # 
│   ├── utils/            # Funciones como `SanityCheck()`, padding, slicing...
├── docs/                 # Herramientas para generacion de documentacion / RAG
├── go.mod                # Manejo de dependencias limpio
├── README.md             # Main page
├── license               # LICENSE
├── plan.md               # Roadmap de implementación, dependencias, etc.
└── logbook.md            # Registro de decisiones, experimentos o insights

```

---

## 🚀 Plan de acción

### 🧱 **Fase 1: Fundamentos**

1. `base/`

   * `CostFunction` interface
   * `Estimator` interface

2. `exceptions/`

   * Definir errores como `ErrNotEnoughPoints`, `ErrBadSegmentationParams`

3. `utils/`

   * `SanityCheck(...)`
   * Otras funciones comunes

---

### ⚙️ **Fase 2: Cost Functions**

4. `cost/l2.go`
5. `cost/l1.go`
6. `cost/rbf.go`
7. `cost/factory.go`

Con cada uno implementando `.Fit(signal)` y `.Error(start, end)`

---

### 🔍 **Fase 3: Algoritmos de Detección**

8. `detection/pelt/`

   * Implementar `Pelt`
9. `detection/binseg/`
10. `detection/dynp/`
11. `detection/window/`

(Con cada uno implementando `Estimator`)

---

### 📊 **Fase 4: Métricas**

12. `metrics/`

* Precision
* Hausdorff
* Coverage

---

### 🧪 **Fase 5: Datasets y pruebas**

13. `datasets/`

* `pw_constant`, `pw_linear`, etc.

14. Pruebas: comparar contra salidas de `ruptures` en Python

---

## ✅ Enfoque

1. Estructura tipica en Go, modular y mantenible.
2. Comenzar por fundamentos (`base`, `cost`), luego algoritmos (`detection`).
3. Pruebas paralelas contra Python (`ruptures`) para validar equivalencia.

---
