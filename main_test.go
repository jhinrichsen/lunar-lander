package main

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

func TestLunarLanderInteractive(t *testing.T) {
	cmd := exec.Command("retrofocal", "lunar-lander.fc")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	reader := bufio.NewReader(stdout)
	writer := bufio.NewWriter(stdin)

	kInputs := []string{
		"0", "0", "0", "0", "0", "0", "0",
		"164.31426784",
		"200", "200", "200", "200", "200", "200", "200",
	}
	inputIndex := 0
	var buffer bytes.Buffer

	var landingTime, impactVelocity, fuelLeft float64

	for {
		char, err := reader.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("read error: %v", err)
		}

		buffer.WriteByte(char)
		current := buffer.String()

		if strings.HasSuffix(current, "K=:") && inputIndex < len(kInputs) {
			k := kInputs[inputIndex]
			t.Logf("PROMPT: K=:  → IN: %s", k)
			writer.WriteString(k + "\n")
			writer.Flush()
			inputIndex++
			buffer.Reset()
			continue
		}

		// Extract final values
		if strings.Contains(current, "ON THE MOON AT") {
			landingTime = extractFloatFromLine(current)
			t.Logf("✅ Landed at %.2f seconds", landingTime)
			buffer.Reset()
		}
		if strings.Contains(current, "IMPACT VELOCITY OF") {
			impactVelocity = extractFloatFromLine(current)
			t.Logf("🛬 Impact velocity: %.2f MPH", impactVelocity)
			buffer.Reset()
		}
		if strings.Contains(current, "FUEL LEFT:") {
			fuelLeft = extractFloatFromLine(current)
			t.Logf("⛽ Fuel remaining: %.2f lbs", fuelLeft)
			buffer.Reset()
		}

		// Respond to YES/NO input prompt
		if strings.Contains(current, "(ANS. YES OR NO)") {
			t.Log("↪️ Responding NO to (ANS. YES OR NO)")
			writer.WriteString("NO\n")
			writer.Flush()
			buffer.Reset()
			continue
		}
	}

	// Clean shutdown
	go io.Copy(io.Discard, reader)

	if err := cmd.Wait(); err != nil {
		t.Fatalf("retrofocal exited with error: %v", err)
	}

	// Final checks
	if landingTime == 0 || impactVelocity == 0 {
		t.Error("❌ Did not extract final landing statistics")
	} else {
		t.Logf("🏁 Final Stats — Time: %.2f sec | Velocity: %.2f MPH | Fuel: %.2f lbs",
			landingTime, impactVelocity, fuelLeft)
	}
}

// Robust float extractor: gets last float from line
func extractFloatFromLine(s string) float64 {
	fields := strings.Fields(s)
	for i := len(fields) - 1; i >= 0; i-- {
		if v, err := strconv.ParseFloat(strings.Trim(fields[i], ":;"), 64); err == nil {
			return v
		}
	}
	return 0
}
