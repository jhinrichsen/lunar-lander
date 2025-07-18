// main.go — Lunar Lander (FOCAL 1969) re-implementation
package main

import (
	"fmt"
	"math"
)

// Constants from the original program (miles & seconds)
const (
	g       = 0.001 // lunar gravity (mi / s²)
	Z       = 1.8   // thrust constant (mi / s² per unit K/M)
	dt      = 10.0  // control interval [s]
	vTolMph = 0.05  // tolerance used in tests
)

// Lander holds the simulation state.
type Lander struct {
	A float64 // altitude (mi)
	V float64 // velocity (mi/s), + up
	M float64 // total mass (lb)
}

func NewLander() *Lander {
	return &Lander{
		A: 120,     // 120 mi ≈ 634 kft (as in listing)
		V: 0,       // starting at rest
		M: 32500.0, // capsule + full fuel
	}
}

// step integrates one constant-K interval S ≤ 10 s.
func (l *Lander) step(k float64, S float64) {
	a := -g + Z*k/l.M
	l.A += l.V*S + 0.5*a*S*S
	l.V += a * S
	l.M -= k * S
}

// simulate runs until lander reaches the surface (A ≤ 0)
// or burnSeq is exhausted. Returns final velocity [mph] and
// remaining fuel [lb].
func simulate(burnSeq []float64) (velMph float64, fuel float64) {
	l := NewLander()

	for i := 0; l.A > 0 && i < len(burnSeq); i++ {
		k := burnSeq[i]
		// limit K between 0 … 200
		if k < 0 {
			k = 0
		} else if k > 200 {
			k = 200
		}
		// if fuel would go negative, shorten step
		S := dt
		if l.M-k*dt < 0 {
			S = l.M / k // burn until zero fuel
		}
		l.step(k, S)
	}

	velMph = math.Abs(l.V) * 3600
	fuel = l.M - 32500 + 16000 // remaining fuel (initial fuel = 16000)
	return
}

func main() {
	// quick interactive demo: constant K = 0
	v, fuel := simulate(make([]float64, 50)) // 50 × 10 s = 500 s
	fmt.Printf("Impact velocity: %.2f mph, fuel left: %.0f lb\n", v, fuel)
}
