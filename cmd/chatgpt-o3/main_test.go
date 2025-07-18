package main

import (
	"math"
	"testing"
)

var badSeq = []float64{
	0, 0, 0, 0, 0, 0, 0,
	170, 200, 200, 200, 200, 200,
	190, 0, 0, 0, 0, 0, 0, 20,
}

var goodSeq = []float64{
	0, 0, 0, 0, 0, 0, 0,
	170, 200, 200, 200, 200, 200,
	170, 0, 0, 30, 0, 8, 10, 9, 100,
}

func almost(a, b float64) bool { return math.Abs(a-b) < 1.5 }

func TestBadLanding(t *testing.T) {
	vel, _ := simulate(badSeq)
	if !almost(vel, 102) {
		t.Fatalf("expected ~102 mph, got %.1f", vel)
	}
}

func TestGoodLanding(t *testing.T) {
	vel, fuel := simulate(goodSeq)
	if !almost(vel, 21) {
		t.Fatalf("expected ~21 mph, got %.1f", vel)
	}
	if fuel > 5 {
		t.Fatalf("expected ≈0 lb fuel, got %.1f", fuel)
	}
}
