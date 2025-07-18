package main

import (
	"math"
	"testing"
)

// Sequences from prompt.
var (
	burnBad = []float64{
		0, 0, 0, 0, 0, 0, 0,
		170, 200, 200, 200, 200, 200,
		190, 0, 0, 0, 0, 0, 0, 20,
	}
	burnGood = []float64{
		0, 0, 0, 0, 0, 0, 0,
		170, 200, 200, 200, 200, 200,
		170, 0, 0, 30, 0, 8, 10, 9, 100,
	}
)

func almost(x, y float64) bool { return math.Abs(x-y) < 2.0 } // ±2 mph

func TestBadLanding(t *testing.T) {
	v, _ := simulate(burnBad)
	if !almost(v, 102.0) {
		t.Fatalf("expected ~102 mph crash, got %.1f", v)
	}
}

func TestGoodLanding(t *testing.T) {
	v, fuel := simulate(burnGood)
	if !almost(v, 21.0) {
		t.Fatalf("expected ~21 mph soft landing, got %.1f", v)
	}
	if fuel > 5 {
		t.Fatalf("expected ≈0 lb fuel left, got %.1f", fuel)
	}
}
