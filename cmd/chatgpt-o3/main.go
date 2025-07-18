package main

import (
	"fmt"
	"math"
)

// Constants (miles, seconds).
const (
	g  = 0.001 // lunar gravity, acts downward (positive)
	Z  = 1.8   // thrust coefficient (mi/s²) * (K/M)
	dt = 10.0  // control step [s]
)

// Lander holds state.
type Lander struct {
	A    float64 // altitude [mi], positive downward from 0 → 120
	V    float64 // velocity [mi/s], positive downward
	Mtot float64 // total mass [lb] (dry + fuel)
}

func NewLander() *Lander {
	return &Lander{
		A:    0,      // 0 at surface, 120 at start
		V:    1.0,    // 3600 mph downward
		Mtot: 32500., // dry + fuel
	}
}

// step integrates for S seconds with constant K (lb/s).
func (l *Lander) step(k, S float64) {
	a := g - Z*k/l.Mtot // thrust acts upward (reduces positive V)

	l.A += l.V*S + 0.5*a*S*S
	l.V += a * S
	l.Mtot -= k * S
}

// simulate returns impact velocity (mph) and remaining fuel (lb).
func simulate(burn []float64) (velMph, fuelLeft float64) {
	l := NewLander()
	const fuel0 = 16000.0
	for i := 0; l.A < 120 && i < len(burn); i++ {
		k := math.Max(0, math.Min(burn[i], 200))
		S := dt
		if l.Mtot-k*dt < 16500 { // fuel would run out; shorten
			S = (l.Mtot - 16500) / k
		}
		l.step(k, S)
	}
	velMph = l.V * 3600 // positive mph downward
	fuelLeft = l.Mtot - 16500
	return
}

func main() {
	fmt.Println("Lunar Lander demo – constant K=0")
	v, fuel := simulate(make([]float64, 60))
	fmt.Printf("Impact %.1f mph, fuel left %.0f lb\n", v, fuel)
}
