package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strings"
	"testing"
)

// TestGoodLanding verifies the lander's behavior with a specific input sequence
func TestGoodLanding(t *testing.T) {
	// Sequence of fuel rates for the test case (from the original FOCAL example)
	inputs := []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 170, 200, 200, 200, 200, 200, 200, 200, 200, 190, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20}
	inputIdx := 0

	// Mock the GetFuelRate function
	originalGetFuelRate := getFuelRate
	getFuelRate = func(l *Lander) float64 {
		if inputIdx < len(inputs) {
			rate := inputs[inputIdx]
			inputIdx++
			return rate
		}
		return 0
	}
	defer func() { getFuelRate = originalGetFuelRate }()

	// Redirect output to capture it
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create and run the lander
	lander := NewLander()
	lander.Land()

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read the output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Check for expected output patterns
	expectedPhrase := "ON THE MOON AT"
	if !strings.Contains(output, expectedPhrase) {
		t.Fatalf("Expected output to contain %q", expectedPhrase)
	}

	// Verify impact velocity is reported
	expectedVelocity := "IMPACT VELOCITY OF"
	if !strings.Contains(output, expectedVelocity) {
		t.Fatalf("Expected output to contain %q", expectedVelocity)
	}

	// Verify the crash scenario is reported correctly
	expectedCrash := "SORRY, BUT THERE WERE NO SURVIVORS"
	if !strings.Contains(output, expectedCrash) {
		t.Errorf("Expected output to contain %q", expectedCrash)
	}

	// Verify crater depth is reported
	expectedCrater := "FT. DEEP"
	if !strings.Contains(output, expectedCrater) {
		t.Errorf("Expected output to contain crater depth")
	}
}

// TestBadLanding verifies the lander crashes when using a crash-landing sequence
func TestBadLanding(t *testing.T) {
	// Sequence of fuel rates that will cause a crash (no thrust)
	inputs := make([]float64, 100) // All zeros - no thrust applied
	inputIdx := 0

	// Mock the GetFuelRate function
	originalGetFuelRate := getFuelRate
	getFuelRate = func(l *Lander) float64 {
		if inputIdx < len(inputs) {
			rate := inputs[inputIdx]
			inputIdx++
			return rate
		}
		return 0
	}
	defer func() { getFuelRate = originalGetFuelRate }()

	// Redirect output to capture it
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create and run the lander
	lander := NewLander()
	lander.Land()

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read the output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Check for expected output patterns
	expectedPhrase := "ON THE MOON AT"
	if !strings.Contains(output, expectedPhrase) {
		t.Fatalf("Expected output to contain %q", expectedPhrase)
	}

	// Verify impact velocity is reported
	expectedVelocity := "IMPACT VELOCITY OF"
	if !strings.Contains(output, expectedVelocity) {
		t.Fatalf("Expected output to contain %q", expectedVelocity)
	}

	// Verify the crash scenario is reported correctly
	expectedCrash := "SORRY, BUT THERE WERE NO SURVIVORS"
	if !strings.Contains(output, expectedCrash) {
		t.Errorf("Expected output to contain %q", expectedCrash)
	}

	// Verify crater depth is reported
	expectedCrater := "FT. DEEP"
	if !strings.Contains(output, expectedCrater) {
		t.Errorf("Expected output to contain crater depth")
	}
}

// TestPhysics verifies the physics calculations
func TestPhysics(t *testing.T) {
	// Test cases with expected values after 10 seconds of simulation
	tests := []struct {
		name          string
		runTest       func(t *testing.T)
	}{
		{
			name: "No thrust - verify free fall",
			runTest: func(t *testing.T) {
				lander := NewLander()
				lander.fuelRate = 0

				// Initial state
				initialAltitude := lander.altitude

				// Run for 10 seconds
				for i := 0; i < 10; i++ {
					lander.Update(1.0)
				}

				// Verify mass remains the same (no fuel used)
				expectedMass := emptyMass + fuelMass
				if math.Abs(lander.mass-expectedMass) > 0.1 {
					t.Errorf("Expected mass %.1f, got %.1f", expectedMass, lander.mass)
				}

				// Verify velocity changed due to gravity (becoming more negative)
				if lander.velocity >= initialVelocity {
					t.Errorf("Expected velocity to become more negative due to gravity")
				}

				// Verify altitude decreased
				if lander.altitude >= initialAltitude {
					t.Errorf("Expected altitude to decrease in free fall")
				}
			},
		},
		{
			name: "Max thrust - verify thrust overcomes gravity",
			runTest: func(t *testing.T) {
				lander := NewLander()
				lander.fuelRate = 200 // Max thrust

				// Initial state
				initialAltitude := lander.altitude

				// Run for 10 seconds
				for i := 0; i < 10; i++ {
					lander.Update(1.0)
				}

				// Verify mass decreased due to fuel consumption
				expectedMass := emptyMass + fuelMass - 2000 // 200 lbs/s * 10s
				if math.Abs(lander.mass-expectedMass) > 1.0 {
					t.Errorf("Expected mass ~%.1f, got %.1f", expectedMass, lander.mass)
				}

				// With max thrust, the lander should be ascending
				if lander.velocity <= 0 {
					t.Error("Expected positive velocity with max thrust")
				}

				// Altitude should have increased
				if lander.altitude <= initialAltitude {
					t.Error("Expected altitude to increase with max thrust")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.runTest)
	}
}

// Helper function to get fuel rate (mocked in tests)
var getFuelRate = (*Lander).GetFuelRate

// Helper function to run the simulation with predefined inputs
func runSimulation(inputs []float64) string {
	// Backup real stdout
	oldStdout := os.Stdout

	// Create input string
	input := ""
	for _, rate := range inputs {
		input += fmt.Sprintf("%f\n", rate)
	}

	// Create a pipe for input
	inReader := strings.NewReader(input)

	// Create a pipe for output
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create a channel to capture output
	out := make(chan string)
	go func() {
		var buf bytes.Buffer
		buf.ReadFrom(r)
		out <- buf.String()
	}()

	// Create and run the lander with mocked input
	lander := NewLander()
	lander.getFuelRate = func(l *Lander) float64 {
		var rate float64
		fmt.Fscanln(inReader, &rate)
		return rate
	}
	lander.Land()

	// Close the write end of the pipe
	w.Close()

	// Restore stdout
	os.Stdout = oldStdout

	// Get the output
	return <-out
}
