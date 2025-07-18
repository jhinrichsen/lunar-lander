package main

import (
	"math"
	"testing"
)

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

func approx(x, y float64) bool { return math.Abs(x-y) < 1.0 }

func TestBadLanding(t *testing.T) {
	vel, _ := simulate(burnBad)
	if !approx(vel, 102) {
		t.Fatalf("want ~102 mph, got %.1f", vel)
	}
}

func TestGoodLanding(t *testing.T) {
	vel, fuel := simulate(burnGood)
	if !approx(vel, 21) {
		t.Fatalf("want ~21 mph, got %.1f", vel)
	}
	if fuel > 5 {
		t.Fatalf("expect ≈0 lb fuel, got %.1f", fuel)
	}
}
