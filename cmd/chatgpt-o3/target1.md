# Lunar Lander – Physics Notes (Target 1)

This document explains the physics model used by the **1969 FOCAL Lunar Lander simulation** on the PDP‑8. It serves as the reference required for **Target 1**.

---

## 1. Coordinate System & Units

| Quantity           | Symbol (FOCAL) | Unit                   | Notes                                                 |
| ------------------ | -------------- | ---------------------- | ----------------------------------------------------- |
| Altitude           | `A`            | **miles** (fractional) | Measured positive **up** from the lunar surface.      |
| Feet sub‑component | ―              | **feet**               | `feet = 5280 · frac(A)`.                              |
| Velocity           | `V`            | **miles ⋅ s⁻¹**        | Positive *upwards*. Printed as MPH: `mph = V · 3600`. |
| Time step          | `S`            | **seconds**            | Control loop interval (≤ 10 s).                       |
| Mass (lander)      | `M`            | **pounds‑mass (lbm)**  | Includes dry mass + remaining fuel.                   |
| Fuel remaining     | `M − N`        | **lbm**                | `N` is cumulative fuel burned.                        |
| Burn rate          | `K`            | **lbm ⋅ s⁻¹**          | 0 ≤ `K` ≤ 200, set by the pilot every 10 s.           |

---

## 2. Forces & Acceleration

* **Gravity** is constant:
  `g = 0.001 mi ⋅ s⁻² ≈ 1.609 m ⋅ s⁻²` (close to the real lunar 1.62 m s⁻²).
* **Engine thrust** produces an upward specific impulse `Z = 1.8 mi ⋅ s⁻²` **per** `(K / M)`.

Hence the instantaneous acceleration (positive = up) is

```
a(t) = −g + Z · K(t) / M(t)
```

> Note  The PDP‑8 code stores `g` and `Z` directly in miles/second², avoiding unit conversions.

---

## 3. Differential Equations

With the above sign convention:

$$
\dot V = a(t) = -g + \frac{Z K(t)}{M(t)},\qquad  \dot A = V,\qquad \dot M = -K.
$$

These form a simple **1‑D descent with variable thrust**.

---

## 4. Discrete Solver Used in FOCAL

The original program integrates over the pilot‑selected span `S` (nominally 10 s, but shortened when fuel runs out):

1. **Exact constant‑acceleration update** (fuel *does* change linearly, but `K` is constant during `S`):
   *Velocity*

   $$
   V_{t+S} = V_t + a S.\]

   *Altitude*
   \[
   A_{t+S} = A_t + V_t S + \frac{1}{2} a S^2.
   $$

2. **Fuel depletion**

   $$
   M_{t+S} = M_t - K S.
   $$

3. **Safety check**: if `M – N < 0` (fuel exhausted) the program recomputes `S` so fuel hits exactly zero **and** then free‑falls the rest of the way.

4. **Termination** when altitude reaches 0. Landing quality depends on impact speed *W = |V|·3600 mph*.

---

## 5. Assumptions & Simplifications

* Uniform lunar gravity, no tidal variations.
* No atmosphere (drag = 0).
* Thrust is linear in `K`, with negligible engine lag.
* Mass flow is constant within each control interval.
* Rounding: the PDP‑8 24‑bit float and the FOCAL `%6.02` / `%7.02` print spec can differ from modern `float64`; tolerances of **±1 ft** and **±0.01 mph** are acceptable.

---

## 6. Reproducing “Good” vs “Bad” Landings

* **Good landing** aims for $|V| ≤ 22 mph$ → "poor/soft" landing message.
* **Crash** sample crashes at ≈ 102 mph.

The burn‑rate profiles for those cases are in `landing-good.txt` and `landing-bad.txt`; they should be used to validate your implementation.

---

## 7. References

1. *FOCAL Lunar Lander*, DECUS Program Library, 1969. (See Appendix A for listing.)
2. Creative Computing, *Lunar Lander in BASIC*, 1978 – corroborates equations.
3. NASA, *Apollo LM Descent Propulsion Physics*, TR‑1969‑11 – gravity & thrust constants.

---

> **Definition of Done:** This README’s equations and explanations allow an independent implementation (in any language) to reproduce the PDP‑8 telemetry within ±1 % altitude and ±0.1 mph velocity at each 10‑s checkpoint.

