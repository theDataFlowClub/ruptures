# Entropía de Shannon como función de costo en PELT

**Objetivo:**
Integrar la **entropía binaria de Shannon** como función de costo dentro del algoritmo PELT, para detectar puntos de cambio donde cambia el patrón de incertidumbre/información en una señal de 0s y 1s.

---

### 🧠 Fundamento matemático

La entropía de Shannon mide la cantidad de incertidumbre en una distribución de probabilidad discreta. Para una variable binaria $X \in \{0,1\}$ con probabilidades $p_0, p_1$, su entropía es:

$$
H(X) = -p_0 \log_2 p_0 - p_1 \log_2 p_1
$$

Aplicado a un **segmento binario $[s, t)$** de una señal:

* Sean:

  * $n = t - s$: longitud del segmento
  * $n_0$: número de ceros en el segmento
  * $n_1 = n - n_0$: número de unos
* Entonces:

  $$
  p_0 = \frac{n_0}{n}, \quad p_1 = \frac{n_1}{n}
  $$
* La **entropía del segmento** es:

  $$
  H(s, t) = -p_0 \log_2 p_0 - p_1 \log_2 p_1
  $$

Y la **función de costo** es:

$$
\text{Cost}_{[s,t)} = n \cdot H(s, t) = -n \cdot \left( \frac{n_0}{n} \log_2 \frac{n_0}{n} + \frac{n_1}{n} \log_2 \frac{n_1}{n} \right)
$$

Que se simplifica a:

$$
\text{Cost}_{[s,t)} = -n_0 \log_2 \left( \frac{n_0}{n} \right) - n_1 \log_2 \left( \frac{n_1}{n} \right)
$$

Esto tiene una interpretación clara:

* El costo es **mínimo (0)** cuando el segmento es puro (todo 0s o todo 1s).
* El costo es **máximo** cuando el segmento es uniforme (mitad 0s, mitad 1s → entropía máxima de 1 bit).

---

### 🔎 Propiedades útiles

* Es una **función convexa** en $p$.
* Refleja la **incertidumbre estructural** del segmento.
* Es **simétrica**: no importa si cambias 0s por 1s.
* **Invariante a rotación** en señales binarias (posición no importa, solo proporción).

---

### 📦 Aplicación en PELT

* Dentro del marco de PELT, esta función se usa como `Cost.Error(start, end)`.
* No requiere estadísticas acumuladas sobre valores, sino **conteo acumulado de 0s y 1s**.

---

### ⚙️ Posible optimización computacional

Para cada índice $i \in [0, n]$, se pueden precomputar:

```go
prefixCount0[i] = número de ceros en y[0:i]
prefixCount1[i] = i - prefixCount0[i]
```

Entonces para cualquier $[s, t)$:

```go
n0 = prefixCount0[t] - prefixCount0[s]
n1 = (t - s) - n0
```

Y luego:

```go
p0 = float64(n0) / float64(t - s)
p1 = float64(n1) / float64(t - s)
cost = -float64(n0)*math.Log2(p0) - float64(n1)*math.Log2(p1)
```

Con protecciones para evitar `log2(0)`:

```go
if p0 > 0 { ... } else { 0.0 }
```

---

### ✅ Ventajas

* Interpretable: muestra dónde cambia la estructura informativa del stream binario.
* Eficiente: usando prefix sums, se puede calcular en $O(1)$ por segmento.
* Robusto a ruido binario, útil para secuencias de eventos o logs.

---

### ⚠️ Consideraciones

* Requiere que la señal esté codificada como binaria explícita (`0` o `1`).
* La métrica de entropía **no es diferenciable**, pero esto no afecta a PELT (que es discreto).
* Cuando $n_0 = 0$ o $n_1 = 0$, el término correspondiente debe omitirse para evitar `log2(0)`.

---

### 📌 Conclusión

La función de entropía binaria de Shannon puede usarse como una función de costo válida dentro del marco de PELT, para detectar cambios en patrones binarios. La implementación es directa usando contadores acumulativos, y puede extenderse a otros contextos discretos (multiclase, categóricos) con mayor complejidad.

---