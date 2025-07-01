# Cost Functions

## **`CostL1` (Least Absolute Deviation)**

La función de costo `CostL1`, también conocida como **Mínima Desviación Absoluta (Mean Absolute Error - MAE)** o **norma L1**, mide la homogeneidad de un segmento calculando la **suma de las diferencias absolutas** entre cada punto de datos del segmento y la **mediana** de ese segmento.

### **Entendiendo la Costo L1**

El costo L1 es una medida de error robusta, lo que significa que es **menos sensible a los valores atípicos (outliers)** en comparación con el costo L2. Esto se debe a que penaliza los errores linealmente, en lugar de cuadráticamente.

La fórmula para el costo L1 de un segmento $[start:end]$ es:

$$\text{Cost}_{\text{L1}}(segmento) = \sum_{i=start}^{end-1} ||\text{signal}[i] - \text{median}(\text{signal}[start:end])||_1$$

Donde $||.||\_1$ representa la **norma L1**, que para un vector de diferencias es la suma de los valores absolutos de sus componentes. Si la señal es multivariada, la mediana se calcula por cada característica (columna), y las desviaciones absolutas para cada característica se suman para obtener el costo total.

### **Características Clave y Uso**

  * **Robustez:** Es preferido cuando se espera que la señal contenga valores atípicos que no deben influir excesivamente en la detección de puntos de cambio.
  * **Modelo de Segmento:** Asume que cada segmento puede ser bien representado por un valor constante (su mediana). Los puntos de cambio se identifican donde la señal se desvía significativamente de esta constante.
  * **Aplicaciones:** Adecuado para señales donde los cambios se manifiestan como saltos abruptos o donde la presencia de ruido impulsivo es común.

### **Implementación en Go**

La implementación de `CostL1` en Go se caracteriza por:

1.  **Estructura `CostL1`:**
    ```go
    type CostL1 struct {
        Signal  types.Matrix // La señal ajustada. Forma (n_samples, n_features).
        MinSize int          // Tamaño mínimo requerido para un segmento válido. Por defecto es 2.
    }
    ```
2.  **Constructor `NewCostL1()`:** Crea una nueva instancia con `MinSize` inicializado a 2.
3.  **Método `Fit(signal types.Matrix) error`:** Almacena la señal de entrada en la instancia de `CostL1`. Esto es un paso crucial antes de calcular cualquier costo.
4.  **Método `Error(start, end int) (float64, error)`:**
      * Extrae el segmento de la señal `c.Signal[start:end]`.
      * Itera sobre cada **característica** (columna) de este segmento.
      * Para cada característica, calcula la **mediana** utilizando la función `stat.Median` de tu paquete `core/stat`.
      * Calcula la suma de las **desviaciones absolutas** de cada punto de esa característica con respecto a su mediana.
      * La suma total de estas desviaciones absolutas a través de todas las características es el costo final del segmento.
      * Incluye validaciones para asegurar que el segmento no sea demasiado corto (`MinSize`) o que los índices estén dentro de los límites válidos, retornando errores apropiados.

-----

## **`CostL2` (Least Squares Deviation)**

La función de costo `CostL2`, también conocida como **Mínimos Cuadrados (Mean Squared Error - MSE)** o **norma L2 al cuadrado**, evalúa la homogeneidad de un segmento cuantificando la **suma de los cuadrados de las diferencias** entre cada punto de datos del segmento y la **media** de ese segmento.

### **Entendiendo la Costo L2**

El costo L2 es una de las métricas de error más comunes y se utiliza ampliamente debido a sus propiedades matemáticas deseables (es diferenciable y convexa). Penaliza fuertemente los errores grandes debido al término cuadrático, lo que la hace **sensible a los valores atípicos**.

La fórmula para el costo L2 de un segmento $[start:end]$ es:

$$\text{Cost}_{\text{L2}}(segmento) = \sum_{i=start}^{end-1} ||\text{signal}[i] - \text{mean}(\text{signal}[start:end])||_2^2$$

Donde $||.||\_2^2$ es la **norma L2 al cuadrado**, que para un vector de diferencias es la suma de los cuadrados de sus componentes. Esta expresión es equivalente a:

$$\text{Cost}_{\text{L2}}(segmento) = (end - start) \times \text{variance}(\text{signal}[start:end])$$

Esto significa que el costo L2 es directamente proporcional a la **varianza del segmento** y a su longitud.

### **Características Clave y Uso**

  * **Sensibilidad a Outliers:** Al penalizar los errores grandes con más severidad, es muy sensible a los valores atípicos. Un solo outlier puede aumentar el costo significativamente.
  * **Modelo de Segmento:** Asume que cada segmento puede ser representado de manera óptima por un valor constante (su media). Los puntos de cambio se detectan donde la media de la señal cambia.
  * **Aplicaciones:** Ideal para señales donde los cambios de fase se manifiestan principalmente como **cambios en la media** y donde se asume que el ruido subyacente sigue una distribución normal (gaussiana).

### **Implementación en Go**

La implementación de `CostL2` en Go se articula en torno a:

1.  **Estructura `CostL2`:**
    ```go
    type CostL2 struct {
        Signal  types.Matrix // La señal ajustada. Forma (n_samples, n_features).
        MinSize int          // Tamaño mínimo requerido para un segmento válido. Por defecto es 1.
    }
    ```
2.  **Constructor `NewCostL2()`:** Crea una nueva instancia con `MinSize` inicializado a 1.
3.  **Método `Fit(signal types.Matrix) error`:** Almacena la señal de entrada en la instancia de `CostL2` para su uso posterior.
4.  **Método `Error(start, end int) (float64, error)`:**
      * Extrae el segmento de la señal `c.Signal[start:end]`.
      * Itera sobre cada **característica** (columna) del segmento.
      * Para cada característica, calcula la **varianza** utilizando la función `stat.Variance` de tu paquete `core/stat`.
      * Suma las varianzas de todas las características.
      * Finalmente, multiplica esta suma por la **longitud del segmento** (`end - start`) para obtener el costo L2 total.
      * Contiene validaciones para la longitud mínima del segmento (`MinSize`) y para que los índices de inicio y fin estén dentro de los límites de la señal.

---

## CostRbf

La función de costo RBF mide la similitud dentro de un segmento utilizando un **kernel de base radial (RBF)**. Un cambio en la media o en la varianza de la señal se detectaría como un cambio en la estructura del kernel. Es particularmente útil para señales donde los cambios no son solo en la media, sino en patrones o distribuciones más complejas.

La fórmula del kernel RBF (también conocido como kernel Gaussiano) para dos puntos $x\_i$ y $x\_j$ es:

$k(x\_i, x\_j) = \\exp(-\\gamma ||x\_i - x\_j||^2)$

donde $\\gamma$ (gamma) es un parámetro de escala (inverso de la varianza del kernel) y $||x\_i - x\_j||^2$ es la distancia euclidiana al cuadrado entre $x\_i$ y $x\_j$.

El costo para un segmento dado se calcula a partir de la matriz de Gram de ese segmento.

### Desafíos Clave para `CostRbf` en Go

1.  **Operaciones Numéricas Eficientes**: Go no tiene `NumPy` nativo. Necesitamos un enfoque para manejar vectores y matrices de manera eficiente.
      * **Opción A (Recomendada inicialmente):** Implementar las operaciones necesarias manualmente en tu paquete `core/utils` o `core/linalg` (álgebra lineal). Esto te da control total y evita dependencias externas al principio.
      * **Opción B:** Evaluar una librería de álgebra lineal en Go. Hay opciones como `gonum/matrix/mat64` o `gorgonia/tensor`. Sin embargo, introducir una dependencia grande al principio podría complicar el proyecto. Para `pdist` y `squareform` específicos, la implementación manual de los bucles para calcular distancias puede ser más directa.
2.  **Parámetro `gamma`**: `CostRbf` tendrá un parámetro `gamma` que puede ser establecido por el usuario. Necesitarás añadirlo a la estructura `CostRbf` en Go. En `ruptures`, `gamma` a menudo se calcula automáticamente si no se especifica. Podríamos implementarlo así:
      * Si `gamma` no se proporciona, se calcula como `1.0 / n_features` (donde `n_features` es la dimensionalidad de la señal). Esto es un valor heurístico común.
3.  **La Lógica de `pdist` y `squareform`**:
      * `pdist`: Iterar sobre todos los pares de puntos en el segmento y calcular la distancia euclidiana al cuadrado.
      * `squareform`: Tomar esas distancias y colocarlas en una matriz cuadrada simétrica.
4.  **Cálculo del Costo Final**: Una vez que la matriz de Gram (`K`) está construida, el costo se calcula de una manera específica (relacionada con la suma de los elementos de `K` o la suma de los cuadrados de los elementos). En `ruptures`, es la suma total de los elementos de la matriz de kernel.

-----

### **Propuesta de Implementación para `CostRbf`**

Considerando los puntos anteriores, podríamos estructurarlo de la siguiente manera:

1.  **Archivo:** `core/cost/rbf.go`
2.  **Estructura `CostRbf`:**
    ```go
    type CostRbf struct {
        Signal  types.Matrix // The fitted signal
        MinSize int
        Gamma   float64      // The RBF kernel parameter
    }
    ```
3.  **`NewCostRbf(gamma float64)`:** Un constructor que inicialice `MinSize` (típicamente 2) y `Gamma`. Si `gamma` es 0 o negativo, podemos usar la heurística `1.0 / n_features` después de `Fit`.
4.  **`Fit(signal types.Matrix) error`:** Igual que las otras funciones de costo. Aquí es donde podrías calcular el `gamma` por defecto si no se proporcionó.
5.  **`Error(start, end int) (float64, error)`:**
      * Extrae el segmento `signal[start:end]`.
      * **Implementa `pdist_squared`:** Una función que calcule las distancias euclidianas al cuadrado entre todos los pares de filas del segmento. Esto te dará un slice de `float64`.
      * **Implementa la construcción de la Matriz de Kernel:** Usando esas distancias, aplica la función `exp(-gamma * dist_sq)` para cada par y construye la matriz de kernel `K`.
      * **Suma los elementos de `K`:** El costo final será la suma de todos los elementos de la matriz de kernel.

---
