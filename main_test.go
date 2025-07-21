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

		// Extract landing stats
		if strings.Contains(current, "ON THE MOON AT") {
			landingTime = extractFloatAfter(current, "ON THE MOON AT", "SEC")
			t.Logf("✅ Landed at %.2f seconds", landingTime)
			buffer.Reset()
		}
		if strings.Contains(current, "IMPACT VELOCITY OF") {
			impactVelocity = extractFloatAfter(current, "IMPACT VELOCITY OF", "M.P.H.")
			t.Logf("🛬 Impact velocity: %.2f MPH", impactVelocity)
			buffer.Reset()
		}
		if strings.Contains(current, "FUEL LEFT:") {
			fuelLeft = extractFloatAfter(current, "FUEL LEFT:", "LBS")
			t.Logf("⛽ Fuel remaining: %.2f lbs", fuelLeft)
			buffer.Reset()
		}

		if strings.Contains(current, "(ANS. YES OR NO)") {
			t.Log("↪️ Responding NO to (ANS. YES OR NO)")
			writer.WriteString("NO\n")
			writer.Flush()
			buffer.Reset()
			continue
		}
	}

	go io.Copy(io.Discard, reader)

	if err := cmd.Wait(); err != nil {
		t.Fatalf("retrofocal exited with error: %v", err)
	}

	// Final checks
	if landingTime == 0 || impactVelocity == 0 {
		t.Error("❌ Did not extract final landing statistics")
	}
}

// Helper to extract float after a label and before an end keyword
func extractFloatAfter(s, after, before string) float64 {
	start := strings.Index(s, after)
	if start == -1 {
		return 0
	}
	start += len(after)
	end := strings.Index(s[start:], before)
	if end == -1 {
		return 0
	}
	field := strings.TrimSpace(s[start : start+end])
	v, err := strconv.ParseFloat(field, 64)
	if err != nil {
		return 0
	}
	return v
}
