# ruptures

A Go re-implementation of change point detection algorithms, inspired by the Python library [`ruptures`](https://github.com/deepcharles/ruptures).

> 🚧 This project is in early development. APIs and module structure are subject to change.

## ✨ Overview

This library provides a modular implementations of change point detection algorithms such as:

- PELT (Pruned Exact Linear Time)
- Binary Segmentation
- Bottom-Up Segmentation
- Window-based approaches

It also includes cost functions (L1, L2, RBF, entropy) and metric evaluation utilities. The goal is to offer a fast, type-safe, and embeddable library in Go for time series segmentation and signal analysis.

## 📦 Features

- Idiomatic Go architecture
- Clean separation between cost functions and detection methods
- Fully compatible with custom input pipelines
- Suitable for CLI tools, web backends, or embedded systems

## 🔧 Project Layout

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

## 🏁 Goals

* Port the core of `ruptures` to Go in idiomatic style
* Enable Go-native change point detection for scientific, industrial, and embedded applications
* Provide test coverage and validation against `ruptures` outputs
* Add statistical utilities and non-parametric tests in the long term

## 📚 Background

This project was inspired by [`ruptures`](https://github.com/deepcharles/ruptures), a Python library authored by [Charles Truong](https://github.com/deepcharles) and contributors, released under a BSD 3-Clause License.

While this project does not copy code from `ruptures`, it reimplements similar concepts and algorithmic structure. We gratefully acknowledge their work as the conceptual foundation of this library.

## 📄 License

See the [LICENSE](LICENSE) file for details.

## 🤝 Contributions

Contributions are welcome! If you're interested in joining the effort or adapting statistical algorithms to Go, feel free to open an issue or pull request.

---
