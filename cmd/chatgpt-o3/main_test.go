// main_test.go — validates “good” & “bad” sequences from the listing.
package main

import (
	"math"
	"testing"
)

// burn-rate sequences (lb/s) copied from the sample printouts.
// Each slice entry = K for the next 10-second tick.
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

func TestBadLanding(t *testing.T) {
	vel, fuel := simulate(burnBad)
	if math.Abs(vel-102.1) > vTolMph {
		t.Fatalf("expected ~102 mph crash, got %.2f mph", vel)
	}
	if fuel < 300 || fuel > 350 {
		t.Fatalf("fuel left mismatch: got %.0f lb", fuel)
	}
}

func TestGoodLanding(t *testing.T) {
	vel, fuel := simulate(burnGood)
	if math.Abs(vel-21.4) > vTolMph {
		t.Fatalf("expected ~21 mph soft landing, got %.2f mph", vel)
	}
	if fuel > 5 {
		t.Fatalf("expected ~0 lb fuel left, got %.0f", fuel)
	}
}
