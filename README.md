# ruptures

A Go re-implementation of change point detection algorithms, inspired by the Python library [`ruptures`](https://github.com/deepcharles/ruptures).

> 🚧 This project is in early development. APIs and module structure are subject to change.

## ✨ Overview

This library provides efficient, modular implementations of change point detection algorithms such as:

- PELT (Pruned Exact Linear Time)
- Binary Segmentation
- Bottom-Up Segmentation
- Window-based approaches

It also includes cost functions (L1, L2, RBF) and metric evaluation utilities. The goal is to offer a fast, type-safe, and embeddable library in Go for time series segmentation and signal analysis.

## 📦 Features

- Idiomatic Go architecture
- Clean separation between cost functions and detection methods
- Fully compatible with custom input pipelines
- Suitable for CLI tools, web backends, or embedded systems
- MIT-licensed and open to commercial or academic use

## 🔧 Project Layout

```text

dxm/
├── core/
│   ├── base/        # Interfaces for CostFunction, Estimator
│   ├── cost/        # Cost functions: L2, L1, RBF, custom
│   ├── detection/   # Algorithms: PELT, Binseg, etc.
│   ├── metrics/     # Precision, recall, coverage
│   ├── utils/       # Helper functions (e.g., sanity checks)
│   └── exceptions/  # Typed error definitions
├── cmd/             # CLI entrypoints (if needed)
├── web/             # (Optional) Web rendering layer

```

## 🔍 Example (Coming Soon)

```go
// Example: PELT on simulated signal
signal := generatePiecewiseConstant(200)
pelt := changepoint.NewPelt(cost.L2{}, minSize=3, jump=5)
pelt.Fit(signal)
bkps := pelt.Predict(penalty=10.0)
fmt.Println("Breakpoints:", bkps)
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

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## 🤝 Contributions

Contributions are welcome! If you're interested in joining the effort or adapting statistical algorithms to Go, feel free to open an issue or pull request.

---