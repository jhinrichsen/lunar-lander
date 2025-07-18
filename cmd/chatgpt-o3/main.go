package main

import (
	"fmt"
	"math"
)

// Physical constants (miles, seconds, pounds)
const (
	g        = 0.001 // lunar gravity (mi/s²)
	Z        = 1.8   // thrust coefficient
	dryMass  = 16500 // lb
	initFuel = 16000 // lb
	cmdTick  = 10.0  // seconds between pilot inputs
)

// Lander state ▸  A = miles to surface (positive), V = mi/s downward (+)
type Lander struct{ A, V, M float64 }

func newLander() *Lander { return &Lander{A: 120, V: 1, M: dryMass + initFuel} }

// oneTick reproduces the PDP-8 five-term series for a constant-K interval S
func (l *Lander) oneTick(K, S float64) {
	if K < 0 {
		K = 0
	}
	if K > 200 {
		K = 200
	}
	if l.M <= dryMass || K == 0 { // out of fuel ⇒ free-fall
		l.A -= l.V*S + .5*g*S*S
		l.V += g * S
		return
	}
	if l.M-K*S < dryMass {
		S = (l.M - dryMass) / K
	} // shorten so fuel hits zero

	Q := S * K / l.M // dimensionless impulse
	// update velocity (five-term)
	l.V += g*S - Z*(Q+Q*Q/2+Q*Q*Q/3+Q*Q*Q*Q/4+Q*Q*Q*Q*Q/5)
	// update altitude (five-term)
	l.A -= l.V*S + .5*g*S*S - Z*S*(Q/2+Q*Q/6+Q*Q*Q/12+Q*Q*Q*Q/20+Q*Q*Q*Q*Q/30)
	l.M -= K * S
}

// simulate executes pilot commands (10-s apart) until A ≤ 0
func simulate(seq []float64) (impactMph, fuelLeft float64) {
	ln := newLander()
	for t := 0; ln.A > 0; t++ {
		K := 0.0
		if t < len(seq) {
			K = seq[t]
		}
		ln.oneTick(K, cmdTick)
	}
	return ln.V * 3600, math.Max(0, ln.M-dryMass)
}

func main() {
	v, f := simulate(make([]float64, 60)) // free-fall demo
	fmt.Printf("Free-fall impact %.1f mph, fuel %.0f lb left\n", v, f)
}
