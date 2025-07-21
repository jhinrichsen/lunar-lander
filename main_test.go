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

	kInputs := []string{
		"0", "0", "0", "0", "0", "0", "0",
		"164.31426784",
		"200", "200", "200", "200", "200", "200", "200",
	}
	inputIndex := 0

	var buffer bytes.Buffer
	var landedPrinted bool

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

		// Accumulate printable text until full line/prompt
		if strings.HasSuffix(current, "K=:") {
			if inputIndex < len(kInputs) {
				k := kInputs[inputIndex]
				t.Logf("PROMPT: %s  → IN: %s", strings.TrimSpace(current), k)
				writer.WriteString(k + "\n")
				writer.Flush()
				inputIndex++
				buffer.Reset()
				continue
			}
		}

		if strings.Contains(current, "ON THE MOON") && !landedPrinted {
			t.Log("✅ Landed")
			landedPrinted = true
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

	// Drain remaining data
	go io.Copy(io.Discard, reader)

	if err := cmd.Wait(); err != nil {
		t.Fatalf("retrofocal exited with error: %v", err)
	}
}
