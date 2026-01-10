package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func runRetrofocal(input string) (string, error) {
	// Locate lunar-lander.fc relative to the test file
	// Assuming test is run from cmd/antigravity/
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get wd: %w", err)
	}

	// Adjust path if running from root vs cmd/antigravity
	// simple heuristic: look for lunar-lander.fc up the tree
	fcPath := filepath.Join(wd, "../../lunar-lander.fc")
	if _, err := os.Stat(fcPath); os.IsNotExist(err) {
		// Try current dir? or panic?
		// For now assume standard layout
		return "", fmt.Errorf("lunar-lander.fc not found at %s", fcPath)
	}

	cmd := exec.Command("retrofocal", fcPath)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	// retrofocal might output to stderr too?
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("retrofocal execution failed: %w", err)
	}

	return out.String(), nil
}

func cleanLines(s string) []string {
	lines := strings.Split(s, "\n")
	var cleaned []string
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if trim == "" {
			continue
		}
		// Ignore interactive terminal prompts differences
		if strings.HasPrefix(trim, "(ANS.") {
			continue
		}
		if strings.HasPrefix(trim, "CONTROL OUT") {
			continue
		}
		if strings.HasPrefix(trim, "K=:") {
			continue
		}

		// Collapse spaces to handle variable whitespace in output columns
		cleaned = append(cleaned, strings.Join(strings.Fields(trim), " "))
	}
	return cleaned
}

func verifyOutput(t *testing.T, expected, actual string) {
	t.Helper()
	expectedNorm := cleanLines(expected)
	actualNorm := cleanLines(actual)

	match := true
	if len(expectedNorm) != len(actualNorm) {
		match = false
	} else {
		for i := range expectedNorm {
			if expectedNorm[i] != actualNorm[i] {
				// Special case: Fuel 0 vs 0.00 mismatch?
				// FOCAL might print 0 or 0.00. Go prints formatted.
				// Let's rely on strict match first.
				match = false
				break
			}
		}
	}

	if !match {
		t.Logf("Expected (Cleaned lines: %d)\n", len(expectedNorm))
		t.Logf("Actual (Cleaned lines: %d)\n", len(actualNorm))

		// Find first mismatch
		limit := len(expectedNorm)
		if len(actualNorm) < limit {
			limit = len(actualNorm)
		}
		for i := 0; i < limit; i++ {
			if expectedNorm[i] != actualNorm[i] {
				t.Logf("Mismatch at line %d:\nEXP: %q\nGOT: %q", i, expectedNorm[i], actualNorm[i])
				break
			}
		}

		t.Errorf("Output mismatch. See logs for detailed diff.")
	}
}

func TestFlightScenariosInternal(t *testing.T) {
	tests := []struct {
		name       string
		inputsFile string
	}{
		{
			name:       "Crash Landing",
			inputsFile: "testdata/inputs_crash.txt",
		},
		{
			name:       "Hard Landing",
			inputsFile: "testdata/inputs_hard.txt",
		},
		{
			name:       "Good Landing",
			inputsFile: "testdata/inputs_good.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputsBytes, err := os.ReadFile(tt.inputsFile)
			if err != nil {
				t.Fatalf("Failed to read input file %s: %v", tt.inputsFile, err)
			}
			inputStr := string(inputsBytes)

			// 1. Run Retrofocal (Ground Truth)
			expectedOutput, err := runRetrofocal(inputStr)
			if err != nil {
				t.Fatalf("Failed to run retrofocal: %v\nOutput: %s", err, expectedOutput)
			}

			// 2. Run Go Implementation
			var outputBuffer bytes.Buffer
			run(strings.NewReader(inputStr), &outputBuffer)
			actualOutput := outputBuffer.String()

			// 3. Verify
			verifyOutput(t, expectedOutput, actualOutput)
		})
	}
}

func FuzzLanding(f *testing.F) {
	// Seed corpus with known good inputs
	inputFiles := []string{
		"testdata/inputs_crash.txt",
		"testdata/inputs_good.txt",
		"testdata/inputs_hard.txt",
	}

	for _, file := range inputFiles {
		data, err := os.ReadFile(file)
		if err == nil {
			f.Add(string(data))
		}
	}

	f.Fuzz(func(t *testing.T, input string) {
		// Filter out inputs that invoke FOCAL's expression evaluation features (variables, math)
		// We only port numeric literal support.
		// Allow digits, whitespace, sign, dot, e/E for float.
		// If input contains anything else (letters, symbols), skip.
		for _, r := range input {
			if !strings.ContainsRune("0123456789.+-eE \n\r\t", r) {
				t.Skip("Skipping non-numeric input (FOCAL expression feature)")
			}
		}

		// 1. Run Retrofocal (Ground Truth)
		// Note: Retrofocal might crash or hang on garbage input?
		// We should probably rely on timeouts or assume robust error handling.
		// For Fuzzing, we mostly care that valid-ish inputs produce same outputs.

		expectedOutput, err := runRetrofocal(input)
		if err != nil {
			// If retrofocal fails, maybe our input was bad?
			// Or maybe Go should also fail similarly?
			// For now, treat retrofocal failure as "skip this test case"
			// t.Skip("Retrofocal failed on input")
			// Actually, let's verify Go output matches the failure message if possible,
			// or just return if it's a runtime error we can't replicate.
			return
		}

		// 2. Run Go Implementation
		var outputBuffer bytes.Buffer
		// Prevent Go panic from crashing the fuzzer
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Go implementation panicked: %v", r)
			}
		}()

		run(strings.NewReader(input), &outputBuffer)
		actualOutput := outputBuffer.String()

		// 3. Verify
		verifyOutput(t, expectedOutput, actualOutput)
	})
}
