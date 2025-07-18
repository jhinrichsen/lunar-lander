package main

import (
	"fmt"
	"math"
)

// Physical constants (miles, seconds, pounds)
const (
	g       = 0.001 // lunar gravity (mi/s²) downward
	Z       = 1.8   // thrust factor (mi/s²) * (K / M)
	dt      = 10.0  // control tick [s]
	dryMass = 16500 // capsule & engine, no fuel [lb]
)

type Lander struct {
	A    float64 // altitude remaining [mi], starts 120 → 0
	V    float64 // velocity [mi/s] (downward positive)
	Mtot float64 // total mass [lb] (dry + fuel)
}

func NewLander() *Lander {
	return &Lander{
		A:    120.0, // 120 miles to surface
		V:    1.0,   // 3600 mph downward
		Mtot: dryMass + 16000,
	}
}

// step integrates for up to S seconds with constant burn K [lb/s].
func (l *Lander) step(K, S float64) {
	a := g - Z*K/l.Mtot // net acceleration (downward positive)
	displacement := l.V*S + 0.5*a*S*S
	l.A -= displacement // altitude decreases
	l.V += a * S        // new velocity
	l.Mtot -= K * S     // burn fuel
}

// simulate runs a sequence of burn-rates (one per 10 s tick).
// Returns |V| in mph at touchdown and fuel left [lb].
func simulate(burnSeq []float64) (mph, fuelLeft float64) {
	ln := NewLander()

	for _, Kcmd := range burnSeq {
		K := math.Max(0, math.Min(Kcmd, 200))

		// shorten step if fuel would empty mid-tick
		S := dt
		if K > 0 && ln.Mtot-K*dt < dryMass {
			S = (ln.Mtot - dryMass) / K
		}

		ln.step(K, S)

		if ln.A <= 0 { // landed / crashed
			break
		}
	}

	mph = ln.V * 3600 // downward mph (positive)
	fuelLeft = ln.Mtot - dryMass
	return
}

func main() {
	demo := make([]float64, 60) // 60 ticks of K=0
	v, fuel := simulate(demo)
	fmt.Printf("Demo free-fall: %.1f mph impact, fuel %.0f lb left\n", v, fuel)
}
