// Lunar-Lander 1969 (FOCAL) – Go 1.24 translation
package main

import (
	"fmt"
	"math"
)

// Physical constants (miles, seconds, pounds).
const (
	g        = 0.001 // lunar gravity (mi/s²), downward
	Z        = 1.8   // thrust coeff (mi/s²) * (K / M)
	tick     = 10.0  // pilot enters K every 10 s
	dt       = 0.05  // sub-integration step [s]
	dryMass  = 16500 // capsule + engine (lb)
	initFuel = 16000 // starting fuel (lb)
)

// State of the lander.
type Lander struct {
	A    float64 // altitude to surface [mi], positive
	V    float64 // velocity [mi/s], downward positive
	Mtot float64 // total mass [lb] (dry + fuel)
}

// newLander returns initial state (matches PDP-8 listing).
func newLander() *Lander {
	return &Lander{
		A:    120.0,
		V:    1.0, // 3600 mph downward
		Mtot: dryMass + initFuel,
	}
}

// subStep integrates one tiny dt with constant K.
func (l *Lander) subStep(K float64) {
	if l.Mtot <= dryMass { // no fuel
		K = 0
	}
	a := g - Z*K/l.Mtot // net acceleration
	l.A -= l.V*dt + 0.5*a*dt*dt
	l.V += a * dt
	l.Mtot -= K * dt
}

// runTick executes one 10-second command interval.
func (l *Lander) runTick(Kcmd float64) {
	K := math.Max(0, math.Min(Kcmd, 200))
	steps := int(tick / dt)
	for i := 0; i < steps && l.A > 0; i++ {
		l.subStep(K)
	}
}

// simulate runs burn-rate commands until touchdown.
// If the sequence ends early, it keeps descending with K = 0.
func simulate(seq []float64) (impactMph, fuelLeft float64) {
	L := newLander()

	for tickIndex := 0; L.A > 0; tickIndex++ {
		var K float64
		if tickIndex < len(seq) {
			K = seq[tickIndex]
		} else {
			K = 0 // after commands run out: free-fall
		}
		L.runTick(K)
	}

	impactMph = L.V * 3600                 // downward mph
	fuelLeft = math.Max(0, L.Mtot-dryMass) // lb
	return
}

func main() {
	dummy := make([]float64, 60) // K=0 free-fall
	v, fuel := simulate(dummy)
	fmt.Printf("Demo: free-fall impact %.1f mph, fuel %.0f lb left\n", v, fuel)
}
