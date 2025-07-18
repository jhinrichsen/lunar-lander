package main

import "testing"

var bad = []float64{
	0, 0, 0, 0, 0, 0, 0,
	170, 200, 200, 200, 200, 200,
	190, 0, 0, 0, 0, 0, 0, 20,
}
var good = []float64{
	0, 0, 0, 0, 0, 0, 0,
	170, 200, 200, 200, 200, 200,
	170, 0, 0, 30, 0, 8, 10, 9, 100,
}

func almost(x, y float64) bool { return (x-y) < 2 && (y-x) < 2 }

func TestBadLanding(t *testing.T) {
	v, _ := simulate(bad)
	if !almost(v, 102) {
		t.Fatalf("want 102 mph, got %.1f", v)
	}
}
func TestGoodLanding(t *testing.T) {
	v, fuel := simulate(good)
	if !almost(v, 21) {
		t.Fatalf("want 21 mph, got %.1f", v)
	}
	if fuel > 5 {
		t.Fatalf("fuel left %.1f lb, expected ≈0", fuel)
	}
}
