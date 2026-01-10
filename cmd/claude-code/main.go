// Package main implements a 1:1 port of the FOCAL lunar-lander.fc simulation.
// This Go implementation produces byte-identical output to the original FOCAL code
// when run with retrofocal.
package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

// Sim holds the simulation state
type Sim struct {
	A float64 // Altitude (miles)
	V float64 // Velocity (miles/sec)
	M float64 // Total mass (lbs)
	N float64 // Dry mass (lbs)
	G float64 // Gravity constant
	Z float64 // Thrust constant
	L float64 // Elapsed time (secs)
	T float64 // Time remaining in current interval
	K float64 // Fuel burn rate (lbs/sec)
	S float64 // Time step
	I float64 // New altitude (from subroutine 9)
	J float64 // New velocity (from subroutine 9)
	W float64 // Velocity in MPH (for landing)

	in  *bufio.Scanner
	out io.Writer
}

// NewSim creates a new simulation with the given input/output
func NewSim(in io.Reader, out io.Writer) *Sim {
	return &Sim{
		in:  bufio.NewScanner(in),
		out: out,
	}
}

// Run executes the simulation
func (s *Sim) Run() {
	for {
		s.intro()
		s.mainLoop()
		if !s.tryAgain() {
			break
		}
	}
}

// intro prints the introduction (lines 01.04-01.40)
func (s *Sim) intro() {
	fmt.Fprintln(s.out, "CONTROL CALLING LUNAR MODULE. MANUAL CONTROL IS NECESSARY")
	fmt.Fprintln(s.out, "YOU MAY RESET FUEL RATE K EACH 10 SECS TO 0 OR ANY VALUE")
	fmt.Fprintln(s.out, "BETWEEN 8 & 200 LBS/SEC. YOU'VE 16000 LBS FUEL. ESTIMATED")
	fmt.Fprintln(s.out, "FREE FALL IMPACT TIME-120 SECS. CAPSULE WEIGHT-32500 LBS")
	fmt.Fprintln(s.out, "FIRST RADAR CHECK COMING UP")
	fmt.Fprintln(s.out)
	fmt.Fprintln(s.out)
	fmt.Fprintln(s.out, "COMMENCE LANDING PROCEDURE")
	fmt.Fprintln(s.out, "TIME,SECS   ALTITUDE,MILES+FEET   VELOCITY,MPH   FUEL,LBS   FUEL RATE")

	// Line 01.50: Initialize variables
	s.A = 120
	s.V = 1
	s.M = 32500
	s.N = 16500
	s.G = 0.001
	s.Z = 1.8
	s.L = 0
}

// fitr returns the integer part of x (FOCAL's FITR function)
func fitr(x float64) float64 {
	return math.Trunc(x)
}

// printStatus prints the current status line (lines 02.10-02.20)
func (s *Sim) printStatus() {
	miles := fitr(s.A)
	feet := 5280 * (s.A - miles)
	velocity := 3600 * s.V
	fuel := s.M - s.N

	// FOCAL format from lines 02.10-02.20, empirically matched to retrofocal output:
	// "        0         120       0         3600.00      16000.0      K=:"
	// Time, Miles, Feet are fixed width
	// Velocity uses 9 spaces + %6.02 format, Fuel uses 6 spaces + %6.01 format
	fmt.Fprintf(s.out, "%9.0f%12.0f%8.0f         %6.2f      %6.1f      K=:",
		s.L, miles, feet, velocity, fuel)
}

// askK prompts for and validates the fuel rate K (lines 02.70-02.73)
func (s *Sim) askK() bool {
	for {
		if !s.in.Scan() {
			return false
		}
		line := strings.TrimSpace(s.in.Text())
		var k float64
		_, err := fmt.Sscanf(line, "%f", &k)
		if err != nil {
			// FOCAL interprets non-numeric input as 0
			k = 0
		}
		s.K = k
		s.T = 10

		// Line 02.70: Validate K
		// I (200-K)2.72 - if K > 200, goto 2.72
		// I (8-K)3.1,3.1 - if K < 8, goto 3.1 (valid if K >= 8)
		// I (K)2.72,3.1 - if K < 0 goto 2.72, if K = 0 goto 3.1, if K > 0 (already handled)
		if s.K > 200 {
			s.printNotPossible()
			continue
		}
		if s.K < 8 {
			if s.K == 0 {
				return true // K=0 is valid
			}
			// K is between 0 and 8 (exclusive), not valid
			s.printNotPossible()
			continue
		}
		// K is between 8 and 200 (inclusive), valid
		return true
	}
}

// printNotPossible prints the "NOT POSSIBLE" message with dots (line 02.72-02.73)
func (s *Sim) printNotPossible() {
	fmt.Fprint(s.out, "NOT POSSIBLE")
	for x := 1; x <= 51; x++ {
		fmt.Fprint(s.out, ".")
	}
	fmt.Fprint(s.out, "K=:")
}

// subroutine9 calculates new velocity (J) and altitude (I) (lines 09.10-09.40)
func (s *Sim) subroutine9() {
	Q := s.S * s.K / s.M
	Q2 := Q * Q
	Q3 := Q2 * Q
	Q4 := Q3 * Q
	Q5 := Q4 * Q
	// J = V + G*S + Z*(-Q - Q^2/2 - Q^3/3 - Q^4/4 - Q^5/5)
	s.J = s.V + s.G*s.S + s.Z*(-Q-Q2/2-Q3/3-Q4/4-Q5/5)
	// I = A - G*S*S/2 - V*S + Z*S*(Q/2 + Q^2/6 + Q^3/12 + Q^4/20 + Q^5/30)
	s.I = s.A - s.G*s.S*s.S/2 - s.V*s.S + s.Z*s.S*(Q/2+Q2/6+Q3/12+Q4/20+Q5/30)
}

// subroutine6 updates state variables (line 06.10)
func (s *Sim) subroutine6() {
	s.L = s.L + s.S
	s.T = s.T - s.S
	s.M = s.M - s.S*s.K
	s.A = s.I
	s.V = s.J
}

// State constants for control flow
const (
	stateGetInput = iota
	stateLoop31
	stateLoop71
	stateLoop81
	stateFuelOut
	stateLanding
	stateDone
)

// mainLoop runs the main simulation loop using state machine
func (s *Sim) mainLoop() {
	state := stateGetInput

	for state != stateDone {
		switch state {
		case stateGetInput:
			// Line 02.10-02.20: Print status and get K
			s.printStatus()
			if !s.askK() {
				return
			}
			state = stateLoop31

		case stateLoop31:
			// Line 03.10: Check fuel and time
			if s.M-s.N < 0.001 {
				state = stateFuelOut
				continue
			}
			if s.T < 0.001 {
				state = stateGetInput
				continue
			}

			s.S = s.T

			// Line 03.40: Check if enough fuel for burn
			if s.N+s.S*s.K > s.M {
				s.S = (s.M - s.N) / s.K
			}

			// Line 03.50: D 9 (call subroutine 9)
			s.subroutine9()

			// I (I)7.1,7.1 - if I <= 0, goto 7.1
			if s.I <= 0 {
				state = stateLoop71
				continue
			}

			// I (V)3.8,3.8 - if V <= 0, goto 3.8
			if s.V <= 0 {
				// Line 03.80: D 6; G 3.1
				s.subroutine6()
				state = stateLoop31
				continue
			}

			// I (J)8.1 - if J < 0, goto 8.1
			if s.J < 0 {
				state = stateLoop81
				continue
			}

			// Line 03.80: D 6; G 3.1
			s.subroutine6()
			state = stateLoop31

		case stateLoop71:
			// Line 07.10: I (S-.005)5.1
			if s.S < 0.005 {
				state = stateLanding
				continue
			}
			// S = 2*A / (V + FSQT(V*V + 2*A*(G - Z*K/M)))
			s.S = 2 * s.A / (s.V + math.Sqrt(s.V*s.V+2*s.A*(s.G-s.Z*s.K/s.M)))
			// Line 07.30: D 9; D 6; G 7.1
			s.subroutine9()
			s.subroutine6()
			// Stay in loop71

		case stateLoop81:
			// Line 08.10
			s.W = (1 - s.M*s.G/(s.Z*s.K)) / 2
			s.S = s.M*s.V/(s.Z*s.K*(s.W+math.Sqrt(s.W*s.W+s.V/s.Z))) + 0.05
			s.subroutine9()

			// Line 08.30: I (I)7.1,7.1
			if s.I <= 0 {
				state = stateLoop71
				continue
			}

			// D 6
			s.subroutine6()

			// I (-J)3.1,3.1 - if J >= 0, goto 3.1
			if s.J >= 0 {
				state = stateLoop31
				continue
			}

			// I (V)3.1,3.1,8.1 - if V <= 0, goto 3.1; if V > 0, goto 8.1
			if s.V <= 0 {
				state = stateLoop31
				continue
			}
			// V > 0: stay in loop81

		case stateFuelOut:
			// Line 04.10 - FOCAL format matches %9.2f for time
			fmt.Fprintf(s.out, "FUEL OUT AT%9.2f SECS\n", s.L)
			// Line 04.40: Free fall calculation
			s.S = (math.Sqrt(s.V*s.V+2*s.A*s.G) - s.V) / s.G
			s.V = s.V + s.G*s.S
			s.L = s.L + s.S
			state = stateLanding

		case stateLanding:
			// Lines 05.10-05.83
			// FOCAL line 05.20 format - uses the last active format spec
			fmt.Fprintf(s.out, "ON THE MOON AT%9.2f SECS\n", s.L)
			s.W = 3600 * s.V
			fmt.Fprintf(s.out, "IMPACT VELOCITY OF%9.2fM.P.H.\n", s.W)
			// FOCAL adds a leading space for positive numbers: min width 9, or 2+numLen if >= 8 chars
			fuelStr := fmt.Sprintf("%.2f", s.M-s.N)
			padWidth := 9
			if len(fuelStr) >= 8 {
				padWidth = 2 + len(fuelStr)
			}
			fmt.Fprintf(s.out, "FUEL LEFT:%*s LBS\n", padWidth, fuelStr)

			if s.W <= 1 {
				fmt.Fprintln(s.out, "PERFECT LANDING !-(LUCKY)")
			} else if s.W <= 10 {
				fmt.Fprintln(s.out, "GOOD LANDING-(COULD BE BETTER)")
			} else if s.W <= 22 {
				fmt.Fprintln(s.out, "CONGRATULATIONS ON A POOR LANDING")
			} else if s.W <= 40 {
				fmt.Fprintln(s.out, "CRAFT DAMAGE. GOOD LUCK")
			} else if s.W <= 60 {
				fmt.Fprintln(s.out, "CRASH LANDING-YOU'VE 5 HRS OXYGEN")
			} else {
				fmt.Fprintln(s.out, "SORRY,BUT THERE WERE NO SURVIVORS-YOU BLEW IT!")
				fmt.Fprintf(s.out, "IN FACT YOU BLASTED A NEW LUNAR CRATER%9.2f FT.DEEP\n", s.W*0.277777)
			}
			state = stateDone
		}
	}
}

// tryAgain prompts for retry (lines 05.90-05.98)
func (s *Sim) tryAgain() bool {
	// FOCAL line 05.90: T !!!!"TRY AGAIN?"! = 4 newlines then TRY AGAIN? then newline
	fmt.Fprintln(s.out)
	fmt.Fprintln(s.out)
	fmt.Fprintln(s.out)
	fmt.Fprintln(s.out)
	fmt.Fprintln(s.out, "TRY AGAIN?")

	for {
		fmt.Fprint(s.out, "(ANS. YES OR NO):")
		if !s.in.Scan() {
			fmt.Fprintln(s.out, "CONTROL OUT")
			fmt.Fprintln(s.out)
			fmt.Fprintln(s.out)
			return false
		}
		line := strings.ToUpper(strings.TrimSpace(s.in.Text()))
		if line == "NO" {
			fmt.Fprintln(s.out, "CONTROL OUT")
			fmt.Fprintln(s.out)
			fmt.Fprintln(s.out)
			return false
		}
		if line == "YES" {
			return true
		}
		// Invalid answer, ask again
	}
}

func main() {
	sim := NewSim(os.Stdin, os.Stdout)
	sim.Run()
}
