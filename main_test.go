package main

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"
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

	// Fuel rates to input
	kInputs := []string{
		"0", "0", "0", "0", "0", "0", "0",
		"164.31426784",
		"200", "200", "200", "200", "200", "200", "200",
	}
	inputIndex := 0

	var buffer bytes.Buffer

	for {
		char, err := reader.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("read error: %v", err)
		}
		buffer.WriteByte(char)
		output := buffer.String()

		// Show character-by-character output
		t.Logf(">> %q", string(char))

		// Respond when prompt appears
		if strings.Contains(output, "K=:") && inputIndex < len(kInputs) {
			k := kInputs[inputIndex]
			t.Logf("IN: %s", k)
			writer.WriteString(k + "\n")
			writer.Flush()
			inputIndex++
			buffer.Reset()
			continue
		}

		if strings.Contains(output, "ON THE MOON") {
			t.Log("✅ Landed")
			break
		}
		if strings.Contains(output, "IMPACT VELOCITY") {
			t.Log("✅ Impact data received")
		}
	}

	if err := cmd.Wait(); err != nil {
		t.Fatalf("retrofocal exited with error: %v", err)
	}
}
