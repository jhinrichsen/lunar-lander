package main

import (
	"math"
	"testing"
)

// TestBadLanding tests the crash landing scenario from the FOCAL output
func TestBadLanding(t *testing.T) {
	lander := NewLunarLander()

	// Input sequence from the bad landing example
	// K values: 0,0,0,0,0,0,170,200,200,200,200,200,200,190,0,0,0,0,0,0,20
	inputs := []float64{0, 0, 0, 0, 0, 0, 170, 200, 200, 200, 200, 200, 200, 190, 0, 0, 0, 0, 0, 0, 20}

	landingTime, impactVelocity, fuelLeft := lander.SimulateWithInputs(inputs)

	// Expected values from the FOCAL output:
	// ON THE MOON AT   214.03 SECS
	// IMPACT VELOCITY OF   102.180 M.P.H.
	// FUEL LEFT:   319.47 LBS
	expectedTime := 214.03
	expectedVelocity := 102.180
	expectedFuel := 319.47

	// Test with tolerance for floating point precision (3 significant figures)
	tolerance := 0.1

	if math.Abs(landingTime-expectedTime) > tolerance {
		t.Errorf("Landing time mismatch: got %.2f, expected %.2f", landingTime, expectedTime)
	}

	if math.Abs(impactVelocity-expectedVelocity) > tolerance {
		t.Errorf("Impact velocity mismatch: got %.3f, expected %.3f", impactVelocity, expectedVelocity)
	}

	if math.Abs(fuelLeft-expectedFuel) > tolerance {
		t.Errorf("Fuel left mismatch: got %.2f, expected %.2f", fuelLeft, expectedFuel)
	}

	// Verify it's classified as a fatal crash (> 60 MPH)
	if impactVelocity <= 60 {
		t.Errorf("Expected fatal crash (>60 MPH), got %.3f MPH", impactVelocity)
	}
}

// TestGoodLanding tests the safe landing scenario from the FOCAL output
func TestGoodLanding(t *testing.T) {
	lander := NewLunarLander()

	// Input sequence from the good landing example
	// K values: 0,0,0,0,0,0,170,200,200,200,200,200,200,170,0,0,30,0,8,10,9,100
	inputs := []float64{0, 0, 0, 0, 0, 0, 170, 200, 200, 200, 200, 200, 200, 170, 0, 0, 30, 0, 8, 10, 9, 100}

	landingTime, impactVelocity, fuelLeft := lander.SimulateWithInputs(inputs)

	// Expected values from the FOCAL output:
	// ON THE MOON AT   226.12 SECS
	// IMPACT VELOCITY OF    21.36 M.P.H.
	// FUEL LEFT:     0.00 LBS
	expectedTime := 226.12
	expectedVelocity := 21.36
	expectedFuel := 0.00

	// Test with tolerance for floating point precision (3 significant figures)
	tolerance := 0.1

	if math.Abs(landingTime-expectedTime) > tolerance {
		t.Errorf("Landing time mismatch: got %.2f, expected %.2f", landingTime, expectedTime)
	}

	if math.Abs(impactVelocity-expectedVelocity) > tolerance {
		t.Errorf("Impact velocity mismatch: got %.3f, expected %.3f", impactVelocity, expectedVelocity)
	}

	if math.Abs(fuelLeft-expectedFuel) > tolerance {
		t.Errorf("Fuel left mismatch: got %.2f, expected %.2f", fuelLeft, expectedFuel)
	}

	// Verify it's classified as a poor landing (10-22 MPH)
	if impactVelocity <= 10 || impactVelocity > 22 {
		t.Errorf("Expected poor landing (10-22 MPH), got %.3f MPH", impactVelocity)
	}
}

// TestInitialConditions verifies the initial state of the lander
func TestInitialConditions(t *testing.T) {
	lander := NewLunarLander()

	// Verify initial conditions match the FOCAL program
	if lander.A != 120.0 {
		t.Errorf("Initial altitude should be 120.0, got %.1f", lander.A)
	}

	if lander.V != 1.0 {
		t.Errorf("Initial velocity should be 1.0, got %.1f", lander.V)
	}

	if lander.M != 32500.0 {
		t.Errorf("Initial total mass should be 32500.0, got %.1f", lander.M)
	}

	if lander.N != 16500.0 {
		t.Errorf("Initial capsule mass should be 16500.0, got %.1f", lander.N)
	}

	// Verify initial fuel is 16000 lbs (M-N)
	initialFuel := lander.M - lander.N
	if initialFuel != 16000.0 {
		t.Errorf("Initial fuel should be 16000.0, got %.1f", initialFuel)
	}

	if lander.G != 0.001 {
		t.Errorf("Gravity constant should be 0.001, got %.3f", lander.G)
	}

	if lander.Z != 1.8 {
		t.Errorf("Thrust coefficient should be 1.8, got %.1f", lander.Z)
	}
}

// TestFuelValidation tests the fuel burn rate validation
func TestFuelValidation(t *testing.T) {
	lander := NewLunarLander()

	// Test with a single valid input
	inputs := []float64{50} // Valid fuel burn rate

	landingTime, impactVelocity, fuelLeft := lander.SimulateWithInputs(inputs)

	// Should complete simulation without errors
	if landingTime <= 0 {
		t.Errorf("Landing time should be positive, got %.2f", landingTime)
	}

	if impactVelocity < 0 {
		t.Errorf("Impact velocity should be positive, got %.3f", impactVelocity)
	}

	if fuelLeft < 0 {
		t.Errorf("Fuel left should not be negative, got %.2f", fuelLeft)
	}
}

// TestFirstStep tests the first simulation step to debug the logic
func TestFirstStep(t *testing.T) {
	lander := NewLunarLander()

	// Test first step with K=0 (no fuel burn)
	K := 0.0
	lander.T = 10.0
	lander.S = lander.T

	t.Logf("Initial state:")
	t.Logf("  Altitude: %.6f miles", lander.A)
	t.Logf("  Velocity: %.6f", lander.V)
	t.Logf("  Mass: %.1f", lander.M)
	t.Logf("  Fuel: %.1f", lander.M-lander.N)

	// Run one step
	lander.subroutine9(K)
	lander.L = lander.L + lander.S
	lander.T = lander.T - lander.S
	lander.M = lander.M - lander.S*K
	lander.A = lander.I
	lander.V = lander.J

	// Expected values from FOCAL output at t=10:
	// TIME=10, ALTITUDE=109 miles 5016 feet, VELOCITY=3636.88 MPH, FUEL=16000.0
	expectedA := 109.0 + 5016.0/5280.0 // Convert feet to miles
	expectedV := 3636.88 / 3600.0       // Convert MPH to miles per time unit
	expectedFuel := 16000.0

	t.Logf("After first step:")
	t.Logf("  Time: %.2f", lander.L)
	t.Logf("  Altitude: %.6f miles (expected %.6f)", lander.A, expectedA)
	t.Logf("  Velocity: %.6f (expected %.6f)", lander.V, expectedV)
	t.Logf("  Fuel: %.1f (expected %.1f)", lander.M-lander.N, expectedFuel)

	// Check if values are reasonably close (within 1%)
	altTolerance := 0.01 * expectedA
	if math.Abs(lander.A-expectedA) > altTolerance {
		t.Errorf("Altitude mismatch: got %.6f, expected %.6f", lander.A, expectedA)
	}
}

// TestNumericalIntegration tests the core physics calculation
func TestNumericalIntegration(t *testing.T) {
	lander := NewLunarLander()

	// Store initial state
	initialA := lander.A

	// Apply numerical integration with known values
	K := 100.0 // Fuel burn rate
	lander.S = 10.0  // Time step

	lander.subroutine9(K)

	// Verify that I and J have been calculated
	if lander.I == 0 && lander.J == 0 {
		t.Error("Subroutine 9 should have calculated I and J values")
	}

	// Verify altitude would decrease (I should be less than initial A)
	if lander.I >= initialA {
		t.Errorf("Calculated altitude should decrease, was %.3f, calculated %.3f", initialA, lander.I)
	}
}

// TestDebugSimulation traces the first few steps of simulation
func TestDebugSimulation(t *testing.T) {
	lander := NewLunarLander()

	// Test first few steps with K=0 (no fuel burn) like the FOCAL examples
	inputs := []float64{0, 0, 0}

	t.Logf("Initial state:")
	t.Logf("  Time: %.2f, Altitude: %.6f, Velocity: %.6f, Mass: %.1f, Fuel: %.1f", 
		lander.L, lander.A, lander.V, lander.M, lander.M-lander.N)

	for i, K := range inputs {
		t.Logf("\n--- Step %d with K=%.0f ---", i+1, K)

		lander.T = 10.0
		lander.S = lander.T

		// Check fuel availability
		if lander.M-lander.N <= 0.001 {
			t.Logf("Fuel out!")
			break
		}

		// Adjust time step if not enough fuel
		if lander.N+lander.S*K > lander.M {
			lander.S = (lander.M - lander.N) / K
			t.Logf("Adjusted time step to %.3f", lander.S)
		}

		// Call subroutine 9
		lander.subroutine9(K)
		t.Logf("After subroutine9: I=%.6f, J=%.6f", lander.I, lander.J)

		// Check conditions
		if lander.I <= 0 {
			t.Logf("I <= 0, would go to landing sequence")
			break
		}
		if lander.V < 0 {
			t.Logf("V < 0, would call subroutine6 and continue")
		}
		if lander.J <= 0 {
			t.Logf("J <= 0, would go to special velocity handling")
		}

		// Normal case: call subroutine6
		lander.subroutine6(K)

		t.Logf("After subroutine6:")
		t.Logf("  Time: %.2f, Altitude: %.6f, Velocity: %.6f, Mass: %.1f, Fuel: %.1f", 
			lander.L, lander.A, lander.V, lander.M, lander.M-lander.N)

		// Expected values from FOCAL output at t=10, 20, 30:
		// t=10: ALT=109 5016, VEL=3636.88, FUEL=16000
		// t=20: ALT=99 4224, VEL=3672.80, FUEL=16000  
		// t=30: ALT=89 2904, VEL=3788.88, FUEL=16000
		if i == 0 {
			expectedA := 109.0 + 5016.0/5280.0
			expectedV := 3636.88 / 3600.0
			t.Logf("Expected at t=10: Alt=%.6f, Vel=%.6f", expectedA, expectedV)
		}
	}
}

// TestDebugBadLanding traces the bad landing sequence step by step
func TestDebugBadLanding(t *testing.T) {
	lander := NewLunarLander()

	// Bad landing sequence from FOCAL output
	inputs := []float64{0, 0, 0, 0, 0, 0, 170, 200, 200, 200, 200, 200, 200, 190, 0, 0, 0, 0, 0, 0, 20}

	// Expected values from FOCAL output at key time points
	expectedValues := map[int]struct {
		alt, vel, fuel float64
	}{
		0:   {120.0, 3600.0, 16000.0},
		10:  {109.0 + 5016.0/5280.0, 3636.88, 16000.0},
		70:  {47.0 + 2904.0/5280.0, 3852.80, 16000.0},
		80:  {37.0 + 1474.0/5280.0, 3539.86, 14300.0},
		140: {0.0 + 5040.0/5280.0, 556.96, 2300.0},
	}

	t.Logf("Initial state:")
	t.Logf("  Time: %.2f, Altitude: %.6f, Velocity: %.6f, Mass: %.1f, Fuel: %.1f", 
		lander.L, lander.A, lander.V, lander.M, lander.M-lander.N)

	for i, K := range inputs {
		if i > 15 { // Stop after a reasonable number of steps for debugging
			break
		}

		t.Logf("\n--- Step %d (t=%.0f) with K=%.0f ---", i+1, lander.L+10, K)

		lander.T = 10.0
		lander.S = lander.T

		// Check fuel availability
		if lander.M-lander.N <= 0.001 {
			t.Logf("Fuel out!")
			break
		}

		// Adjust time step if not enough fuel
		if lander.N+lander.S*K > lander.M {
			lander.S = (lander.M - lander.N) / K
			t.Logf("Adjusted time step to %.3f", lander.S)
		}

		// Call subroutine 9
		lander.subroutine9(K)
		t.Logf("After subroutine9: I=%.6f, J=%.6f", lander.I, lander.J)

		// Check conditions
		if lander.I <= 0 {
			t.Logf("I <= 0, would go to landing sequence")
			break
		}
		if lander.V < 0 {
			t.Logf("V < 0, would call subroutine6 and continue")
		}
		if lander.J <= 0 {
			t.Logf("J <= 0, would go to special velocity handling")
		}

		// Normal case: call subroutine6
		lander.subroutine6(K)

		t.Logf("After subroutine6:")
		t.Logf("  Time: %.2f, Altitude: %.6f, Velocity: %.6f, Mass: %.1f, Fuel: %.1f", 
			lander.L, lander.A, lander.V, lander.M, lander.M-lander.N)

		// Compare with expected values at key time points
		timeKey := int(lander.L)
		if expected, exists := expectedValues[timeKey]; exists {
			t.Logf("Expected at t=%d: Alt=%.6f, Vel=%.2f, Fuel=%.1f", 
				timeKey, expected.alt, expected.vel, expected.fuel)

			altDiff := math.Abs(lander.A - expected.alt)
			velDiff := math.Abs(3600*lander.V - expected.vel)
			fuelDiff := math.Abs((lander.M-lander.N) - expected.fuel)

			t.Logf("Differences: Alt=%.6f, Vel=%.2f, Fuel=%.1f", altDiff, velDiff, fuelDiff)
		}
	}
}
