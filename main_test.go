package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

func parseTelemetry(line string) (time int, miles int, feet int, velocity, fuel, k float64, err error) {
	fields := strings.Fields(line)
	if len(fields) < 12 {
		err = fmt.Errorf("unexpected field count: %d", len(fields))
		return
	}

	time, _ = strconv.Atoi(fields[0])
	miles, _ = strconv.Atoi(fields[1])
	feet, _ = strconv.Atoi(fields[2])
	velocity, _ = strconv.ParseFloat(fields[3], 64)
	fuel, _ = strconv.ParseFloat(fields[4], 64)
	k, _ = strconv.ParseFloat(fields[7], 64)
	return
}

func TestLunarLanderTelemetry(t *testing.T) {
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

	// K input values
	inputs := []string{
		"0", "0", "0", "0", "0", "0", "0",
		"164.31426784",
		"200", "200", "200", "200", "200", "200", "200",
	}
	inputIndex := 0

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("read error: %v", err)
		}
		line = strings.TrimSpace(line)
		t.Logf("OUT: %s", line)

		if strings.Contains(line, "K=:") && inputIndex < len(inputs) {
			k := inputs[inputIndex]
			t.Logf("IN: %s", k)
			writer.WriteString(k + "\n")
			writer.Flush()
			inputIndex++
		}

		if strings.HasPrefix(line, "ON THE MOON") {
			t.Log("✅ Landed successfully")
		}
		if strings.HasPrefix(line, "IMPACT VELOCITY OF") {
			// parse velocity
			fields := strings.Fields(line)
			if len(fields) < 5 {
				t.Errorf("bad impact line: %q", line)
				continue
			}
			v, err := strconv.ParseFloat(fields[3], 64)
			if err != nil {
				t.Errorf("failed to parse impact velocity: %v", err)
			} else if v > 5.0 {
				t.Errorf("🚨 Crash landing! Impact velocity too high: %.2f MPH", v)
			} else {
				t.Logf("🛬 Safe landing. Impact velocity: %.2f MPH", v)
			}
		}
		if strings.HasPrefix(line, "FUEL LEFT:") {
			fields := strings.Fields(line)
			if len(fields) < 3 {
				t.Errorf("bad fuel line: %q", line)
				continue
			}
			f, err := strconv.ParseFloat(fields[2], 64)
			if err != nil {
				t.Errorf("could not parse fuel: %v", err)
			} else if f < 0 {
				t.Errorf("negative fuel! %.2f lbs", f)
			} else {
				t.Logf("Fuel remaining: %.2f lbs", f)
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		t.Fatalf("simulation exited with error: %v", err)
	}
}
