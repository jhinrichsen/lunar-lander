package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

// Test cases with input sequences and descriptions
// Each test must provide enough inputs to complete the simulation (land or crash)
var testCases = []struct {
	name   string
	inputs string
}{
	{
		name:   "immediate_crash_max_burn",
		inputs: "200\n200\n200\n200\n200\n200\n200\n200\n200\n200\n200\n200\nNO\n",
	},
	{
		name:   "coast_only",
		inputs: "0\n0\n0\n0\n0\n0\n0\n0\n0\n0\n0\n0\nNO\n",
	},
	{
		name:   "good_landing",
		inputs: "0\n0\n0\n0\n0\n0\n200\n200\n200\n200\n200\n0\n0\n100\n200\n200\n0\n0\n71\n0\nNO\n",
	},
	{
		name:   "perfect_landing",
		inputs: "0\n0\n0\n0\n0\n0\n200\n200\n200\n200\n200\n0\n0\n100\n200\n200\n0\n0\n71\n37\nNO\n",
	},
	{
		// Tests K=5 invalid (NOT POSSIBLE), then K=0 valid, then complete landing sequence
		name:   "invalid_then_valid",
		inputs: "5\n0\n0\n0\n0\n0\n0\n0\n0\n0\n0\n0\n0\nNO\n",
	},
	{
		// Tests K=8 boundary - slow burn needs many steps to land
		name:   "boundary_k_8",
		inputs: "8\n8\n8\n8\n8\n8\n8\n8\n8\n8\n8\n8\n8\n8\n8\n8\n8\n8\n8\n8\n8\n8\n8\n8\n8\nNO\n",
	},
	{
		// Tests K=200 then coast - completes with crash (needs 13 K values before landing)
		name:   "boundary_k_200",
		inputs: "200\n0\n0\n0\n0\n0\n0\n0\n0\n0\n0\n0\n0\nNO\n",
	},
	{
		// Alternating burns - needs enough steps to land
		name:   "mixed_burns",
		inputs: "0\n100\n0\n100\n0\n100\n0\n100\n0\n100\n0\n100\n0\n100\n0\nNO\n",
	},
	{
		// Full max burn exhausts fuel at 80 secs
		name:   "fuel_exhaustion",
		inputs: "200\n200\n200\n200\n200\n200\n200\n200\n200\n200\n200\n200\nNO\n",
	},
}

// runFOCAL runs the FOCAL simulation with retrofocal
func runFOCAL(t *testing.T, inputs string) string {
	cmd := exec.Command("retrofocal", "../../lunar-lander.fc")
	cmd.Stdin = strings.NewReader(inputs)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		// retrofocal may exit with error on some inputs, that's ok
		t.Logf("retrofocal exited with: %v", err)
	}
	return out.String()
}

// runGo runs the Go simulation
func runGo(t *testing.T, inputs string) string {
	var out bytes.Buffer
	sim := NewSim(strings.NewReader(inputs), &out)
	sim.Run()
	return out.String()
}

// TestIdenticalOutput tests that Go and FOCAL produce identical output
func TestIdenticalOutput(t *testing.T) {
	// Check if retrofocal is available
	if _, err := exec.LookPath("retrofocal"); err != nil {
		t.Skip("retrofocal not found in PATH, skipping comparison tests")
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			focalOut := runFOCAL(t, tc.inputs)
			goOut := runGo(t, tc.inputs)

			if focalOut != goOut {
				t.Errorf("Output mismatch for %s", tc.name)
				t.Logf("Inputs:\n%s", tc.inputs)

				// Show diff
				focalLines := strings.Split(focalOut, "\n")
				goLines := strings.Split(goOut, "\n")

				maxLines := len(focalLines)
				if len(goLines) > maxLines {
					maxLines = len(goLines)
				}

				for i := 0; i < maxLines; i++ {
					var fLine, gLine string
					if i < len(focalLines) {
						fLine = focalLines[i]
					}
					if i < len(goLines) {
						gLine = goLines[i]
					}

					if fLine != gLine {
						t.Logf("Line %d differs:", i+1)
						t.Logf("  FOCAL: %q", fLine)
						t.Logf("  Go:    %q", gLine)
						// Show byte-level diff for first difference
						if len(fLine) > 0 && len(gLine) > 0 {
							for j := 0; j < len(fLine) || j < len(gLine); j++ {
								var fb, gb byte
								if j < len(fLine) {
									fb = fLine[j]
								}
								if j < len(gLine) {
									gb = gLine[j]
								}
								if fb != gb {
									t.Logf("  First diff at char %d: FOCAL=%q(%d) Go=%q(%d)",
										j, string(fb), fb, string(gb), gb)
									break
								}
							}
						}
					}
				}
			}
		})
	}
}

// TestGoSimulation tests the Go simulation independently
func TestGoSimulation(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			sim := NewSim(strings.NewReader(tc.inputs), &out)
			sim.Run()

			output := out.String()
			if output == "" {
				t.Error("No output produced")
			}

			// Basic sanity checks
			if !strings.Contains(output, "CONTROL CALLING LUNAR MODULE") {
				t.Error("Missing intro text")
			}
			if !strings.Contains(output, "K=:") {
				t.Error("Missing K prompt")
			}
		})
	}
}

// TestPerfectLanding specifically tests the optimal sequence
func TestPerfectLanding(t *testing.T) {
	inputs := "0\n0\n0\n0\n0\n0\n200\n200\n200\n200\n200\n0\n0\n100\n200\n200\n0\n0\n71\n37\nNO\n"

	var out bytes.Buffer
	sim := NewSim(strings.NewReader(inputs), &out)
	sim.Run()

	output := out.String()

	if !strings.Contains(output, "PERFECT LANDING") {
		t.Error("Expected PERFECT LANDING message")
		t.Logf("Output:\n%s", output)
	}

	if !strings.Contains(output, "0.66") || !strings.Contains(output, "M.P.H.") {
		t.Error("Expected impact velocity around 0.66 M.P.H.")
	}
}

// TestMaxBurnCrash tests the immediate crash scenario
func TestMaxBurnCrash(t *testing.T) {
	inputs := "200\n200\n200\n200\n200\n200\n200\n200\n200\n200\n200\n200\nNO\n"

	var out bytes.Buffer
	sim := NewSim(strings.NewReader(inputs), &out)
	sim.Run()

	output := out.String()

	if !strings.Contains(output, "FUEL OUT AT") {
		t.Error("Expected FUEL OUT message")
	}

	if !strings.Contains(output, "NO SURVIVORS") {
		t.Error("Expected NO SURVIVORS message for high-speed crash")
		t.Logf("Output:\n%s", output)
	}
}

// TestInvalidInput tests invalid K value handling
func TestInvalidInput(t *testing.T) {
	// K=5 is invalid (must be 0 or 8-200)
	inputs := "5\n0\nNO\n"

	var out bytes.Buffer
	sim := NewSim(strings.NewReader(inputs), &out)
	sim.Run()

	output := out.String()

	if !strings.Contains(output, "NOT POSSIBLE") {
		t.Error("Expected NOT POSSIBLE message for invalid K=5")
		t.Logf("Output:\n%s", output)
	}
}

// TestByteIdentical runs a comparison and reports if outputs are byte-identical
func TestByteIdentical(t *testing.T) {
	if _, err := exec.LookPath("retrofocal"); err != nil {
		t.Skip("retrofocal not found in PATH")
	}

	inputs := "0\n0\n0\n0\n0\n0\n200\n200\n200\n200\n200\n0\n0\n100\n200\n200\n0\n0\n71\n37\nNO\n"

	focalOut := runFOCAL(t, inputs)
	goOut := runGo(t, inputs)

	if focalOut == goOut {
		t.Log("SUCCESS: Outputs are byte-identical!")
	} else {
		t.Errorf("Outputs differ by %d bytes", abs(len(focalOut)-len(goOut)))
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// FuzzIdenticalOutput uses fuzzing to find inputs where Go and FOCAL outputs differ
// Note: Very long simulations (>20 steps) may accumulate floating-point precision
// differences between Go and retrofocal. This is expected for edge cases.
func FuzzIdenticalOutput(f *testing.F) {
	// Check if retrofocal is available
	if _, err := exec.LookPath("retrofocal"); err != nil {
		f.Skip("retrofocal not found in PATH, skipping fuzz tests")
	}

	// Seed with known working cases
	for _, tc := range testCases {
		f.Add([]byte(tc.inputs))
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		// Convert random bytes to valid K values
		inputs := bytesToKSequence(data)
		if inputs == "" {
			t.Skip("no valid inputs generated")
		}

		// Run both simulators
		focalOut := runFOCALFuzz(inputs)
		goOut := runGoFuzz(inputs)

		// Compare outputs
		if focalOut != goOut {
			t.Errorf("Output mismatch!\nInputs: %q\nFOCAL len: %d\nGo len: %d",
				inputs, len(focalOut), len(goOut))

			// Find first difference
			minLen := len(focalOut)
			if len(goOut) < minLen {
				minLen = len(goOut)
			}
			for i := 0; i < minLen; i++ {
				if focalOut[i] != goOut[i] {
					start := i - 20
					if start < 0 {
						start = 0
					}
					end := i + 20
					if end > minLen {
						end = minLen
					}
					t.Errorf("First diff at byte %d:\n  FOCAL: %q\n  Go:    %q",
						i, focalOut[start:end], goOut[start:end])
					break
				}
			}
		}
	})
}

// FuzzShortSequences focuses on short burn sequences that complete quickly
// These are more likely to produce byte-identical output
func FuzzShortSequences(f *testing.F) {
	if _, err := exec.LookPath("retrofocal"); err != nil {
		f.Skip("retrofocal not found in PATH, skipping fuzz tests")
	}

	// Seed with patterns that lead to quick landings/crashes
	// Each seed will be padded to 20 K values in bytesToShortSequence
	f.Add([]byte{255, 255, 255, 255, 255, 255, 255, 255}) // max burn crash
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0, 0})                 // coast only

	f.Fuzz(func(t *testing.T, data []byte) {
		inputs := bytesToShortSequence(data)
		if inputs == "" {
			t.Skip("no valid inputs generated")
		}

		focalOut := runFOCALFuzz(inputs)
		goOut := runGoFuzz(inputs)

		if focalOut != goOut {
			t.Errorf("Output mismatch!\nInputs: %q", inputs)
		}
	})
}

// bytesToShortSequence creates sequences using only 0 and 200
// Padded to 35 K values to ensure simulation completes before "NO"
func bytesToShortSequence(data []byte) string {
	if len(data) < 1 {
		return ""
	}

	var sb strings.Builder
	kCount := 0
	maxK := 8 // Use first 8 bytes for K pattern

	for i := 0; i < len(data) && kCount < maxK; i++ {
		// Only use 0 or 200 to minimize precision differences
		if data[i] < 128 {
			sb.WriteString("0\n")
		} else {
			sb.WriteString("200\n")
		}
		kCount++
	}

	// Pad with zeros to ensure landing completes before "NO"
	// 50 K values covers worst-case slow coasting scenarios (e.g., 7x200 burn)
	for kCount < 50 {
		sb.WriteString("0\n")
		kCount++
	}

	sb.WriteString("NO\n")
	return sb.String()
}

// bytesToKSequence converts random bytes to a valid K input sequence
func bytesToKSequence(data []byte) string {
	if len(data) < 1 {
		return ""
	}

	var sb strings.Builder
	kCount := 0
	maxK := 30 // Limit iterations to ensure simulation completes

	for i := 0; i < len(data) && kCount < maxK; i++ {
		b := data[i]
		var k int

		// Map byte to valid K values: 0 or 8-200
		if b < 32 {
			k = 0 // ~12.5% chance of coast
		} else {
			// Map 32-255 to 8-200
			k = 8 + int(b-32)*(200-8)/(255-32)
			if k > 200 {
				k = 200
			}
		}

		sb.WriteString(fmt.Sprintf("%d\n", k))
		kCount++
	}

	// Ensure we have enough K values for a complete simulation
	// The simulation can need up to ~50 K values for slow coasting scenarios
	// Adding extra zeros ensures "NO" is only consumed by TRY AGAIN prompt
	for kCount < 50 {
		sb.WriteString("0\n")
		kCount++
	}

	sb.WriteString("NO\n")
	return sb.String()
}

// runFOCALFuzz runs retrofocal without test context
func runFOCALFuzz(inputs string) string {
	cmd := exec.Command("retrofocal", "../../lunar-lander.fc")
	cmd.Stdin = strings.NewReader(inputs)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Run() // Ignore errors
	return out.String()
}

// runGoFuzz runs Go simulation without test context
func runGoFuzz(inputs string) string {
	var out bytes.Buffer
	sim := NewSim(strings.NewReader(inputs), &out)
	sim.Run()
	return out.String()
}
