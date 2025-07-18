package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// LunarLander represents the state of the lunar lander simulation
type LunarLander struct {
	A float64 // Altitude above lunar surface (miles)
	V float64 // Velocity (miles per unit time)
	M float64 // Total spacecraft mass (lbs)
	N float64 // Remaining fuel mass (lbs)
	L float64 // Mission elapsed time (seconds)
	T float64 // Time step for numerical integration (seconds)
	G float64 // Lunar gravitational acceleration constant
	Z float64 // Thrust efficiency coefficient
	S float64 // Adaptive time step
	I float64 // Temporary altitude variable
	J float64 // Temporary velocity variable
	Q float64 // Temporary calculation variable
	W float64 // Impact velocity in MPH
}

// NewLunarLander creates a new lunar lander with initial conditions
func NewLunarLander() *LunarLander {
	return &LunarLander{
		A: 120.0,   // Starting altitude: 120 miles
		V: 1.0,     // Starting velocity: 1 mile/time unit
		M: 32500.0, // Total current mass: 32,500 lbs (matches FOCAL M=32500)
		N: 16500.0, // Capsule dry mass: 16,500 lbs (matches FOCAL N=16500)
		L: 0.0,     // Mission elapsed time starts at 0
		T: 10.0,    // Time step: 10 seconds
		G: 0.001,   // Lunar gravity constant
		Z: 1.8,     // Thrust efficiency coefficient
	}
}

// RunSimulation runs the main simulation loop following FOCAL logic
func (ll *LunarLander) RunSimulation() {
	scanner := bufio.NewScanner(os.Stdin)

	// 01.04-01.20: Initial messages
	fmt.Println("CONTROL CALLING LUNAR MODULE. MANUAL CONTROL IS NECESSARY")
	fmt.Println("YOU MAY RESET FUEL RATE K EACH 10 SECS TO 0 OR ANY VALUE")
	fmt.Println("BETWEEN 8 & 200 LBS/SEC. YOU'VE 16000 LBS FUEL. ESTIMATED")
	fmt.Println("FREE FALL IMPACT TIME=120 SECS. CAPSULE WEIGHT=32500 LBS")
	fmt.Println("FIRST RADAR CHECK COMING UP")
	fmt.Println()
	fmt.Println("COMMENCE LANDING PROCEDURE")
	fmt.Println("TIME,SECS   ALTITUDE,MILES+FEET   VELOCITY,MPH   FUEL,LBS   FUEL RATE")

	for {
		// 02.10-02.20: Display status and get input
		altitudeMiles := int(ll.A)
		altitudeFeet := int(5280 * (ll.A - float64(altitudeMiles)))
		velocityMPH := 3600 * ll.V
		fuelRemaining := ll.M - ll.N

		fmt.Printf("%8.0f %13d %6d %12.2f %9.1f       K=:",
			ll.L, altitudeMiles, altitudeFeet, velocityMPH, fuelRemaining)

		// Get fuel burn rate input
		var K float64
		if scanner.Scan() {
			input := strings.TrimSpace(scanner.Text())
			if input == "" {
				K = 0
			} else {
				var err error
				K, err = strconv.ParseFloat(input, 64)
				if err != nil {
					K = 0
				}
			}
		}

		ll.T = 10 // Reset time step

		// 02.70-02.73: Validate fuel burn rate
		if K > 200 || (K > 0 && K < 8) {
			fmt.Print("NOT POSSIBLE")
			for i := 0; i < 51; i++ {
				fmt.Print(".")
			}
			fmt.Print("K=:")
			continue
		}

		// 03.10: Check fuel availability
		if ll.M-ll.N <= 0.001 {
			// 04.10: Fuel out
			fmt.Printf("FUEL OUT AT %.2f SECS\n", ll.L)
			// 04.40: Calculate impact
			ll.S = (math.Sqrt(ll.V*ll.V+2*ll.A*ll.G) - ll.V) / ll.G
			ll.V = ll.V + ll.G*ll.S
			ll.L = ll.L + ll.S
			ll.landingSequence()
			return
		}

		// 03.10: Check time step
		if ll.T <= 0.001 {
			continue
		}

		ll.S = ll.T

		// 03.40: Adjust time step if not enough fuel
		if ll.N+ll.S*K > ll.M {
			ll.S = (ll.M - ll.N) / K
		}

		// 03.50: Call subroutine 9 (numerical integration)
		ll.subroutine9(K)

		// 03.50: I (I)7.1,7.1;I (V)3.8,3.8;I (J)8.1
		if ll.I <= 0 {
			// 07.10: I (S-.005)5.1 - if S > 0.005 then go to 5.1 (landing)
			if ll.S > 0.005 {
				ll.landingSequence()
				return
			}
			ll.S = 2 * ll.A / (ll.V + math.Sqrt(ll.V*ll.V+2*ll.A*(ll.G-ll.Z*K/ll.M)))
			ll.subroutine9(K)
			ll.subroutine6(K)
			continue
		}

		if ll.V < 0 {
			// 03.80: D 6;G 3.1 - Call subroutine 6 and continue
			ll.subroutine6(K)
			continue
		}

		if ll.J <= 0 {
			// 08.10: Special velocity handling
			ll.W = (1 - ll.M*ll.G/(ll.Z*K)) / 2
			ll.S = ll.M*ll.V/(ll.Z*K*(ll.W+math.Sqrt(ll.W*ll.W+ll.V/ll.Z))) + 0.05
			ll.subroutine9(K)

			// 08.30: I (I)7.1,7.1;D 6;I (-J)3.1,3.1;I (V)3.1,3.1,8.1
			if ll.I <= 0 {
				continue
			}
			ll.subroutine6(K)
			if ll.J < 0 {
				continue
			}
			if ll.V <= 0 {
				continue
			}
		}

		// Normal case: 03.80: D 6;G 3.1 - Call subroutine 6 and continue
		ll.subroutine6(K)
	}
}

// subroutine9 implements FOCAL subroutine 9 (lines 09.10-09.40)
func (ll *LunarLander) subroutine9(K float64) {
	ll.Q = ll.S * K / ll.M
	ll.J = ll.V + ll.G*ll.S + ll.Z*(-ll.Q-math.Pow(ll.Q, 2)/2-math.Pow(ll.Q, 3)/3-math.Pow(ll.Q, 4)/4-math.Pow(ll.Q, 5)/5)
	ll.I = ll.A - ll.G*ll.S*ll.S/2 - ll.V*ll.S + ll.Z*ll.S*(ll.Q/2+math.Pow(ll.Q, 2)/6+math.Pow(ll.Q, 3)/12+math.Pow(ll.Q, 4)/20+math.Pow(ll.Q, 5)/30)
}

// subroutine6 implements FOCAL subroutine 6 (state update)
// FOCAL line 06.10: S L=L+S;S T=T-S;S M=M-S*K;S A=I;S V=J
func (ll *LunarLander) subroutine6(K float64) {
	ll.L = ll.L + ll.S
	ll.T = ll.T - ll.S
	ll.M = ll.M - ll.S*K // Consume fuel
	ll.A = ll.I
	ll.V = ll.J
}

// landingSequence handles the landing evaluation
func (ll *LunarLander) landingSequence() {
	fmt.Printf("ON THE MOON AT %.2f SECS\n", ll.L)
	ll.W = 3600 * ll.V
	fmt.Printf("IMPACT VELOCITY OF %.3f M.P.H.\n", ll.W)
	fmt.Printf("FUEL LEFT: %.2f LBS\n", ll.M-ll.N)

	// Evaluate landing quality
	if ll.W <= 1 {
		fmt.Println("PERFECT LANDING !-(LUCKY)")
	} else if ll.W <= 10 {
		fmt.Println("GOOD LANDING-(COULD BE BETTER)")
	} else if ll.W <= 22 {
		fmt.Println("CONGRATULATIONS ON A POOR LANDING")
	} else if ll.W <= 40 {
		fmt.Println("CRAFT DAMAGE. GOOD LUCK")
	} else if ll.W <= 60 {
		fmt.Println("CRASH LANDING-YOU'VE 5 HRS OXYGEN")
	} else {
		fmt.Println("SORRY,BUT THERE WERE NO SURVIVORS-YOU BLEW IT!")
		fmt.Printf("IN FACT YOU BLASTED A NEW LUNAR CRATER %.2f FT. DEEP\n", ll.W*0.277777)
	}

	fmt.Println()
	fmt.Println("TRY AGAIN?")
}

// SimulateWithInputs runs simulation with predefined inputs (for testing)
func (ll *LunarLander) SimulateWithInputs(inputs []float64) (float64, float64, float64) {
	inputIndex := 0

	for {
		// Get fuel burn rate
		var K float64
		if inputIndex < len(inputs) {
			K = inputs[inputIndex]
			inputIndex++
		} else {
			K = 0
		}

		ll.T = 10 // Reset time step

		// Validate fuel burn rate (skip invalid inputs)
		if K > 200 || (K > 0 && K < 8) {
			continue
		}

		// 03.10: Check fuel availability
		if ll.M-ll.N <= 0.001 {
			// 04.40: Calculate impact when fuel runs out
			ll.S = (math.Sqrt(ll.V*ll.V+2*ll.A*ll.G) - ll.V) / ll.G
			ll.V = ll.V + ll.G*ll.S
			ll.L = ll.L + ll.S
			break
		}

		// Check time step
		if ll.T <= 0.001 {
			continue
		}

		ll.S = ll.T

		// 03.40: Adjust time step if not enough fuel
		if ll.N+ll.S*K > ll.M {
			ll.S = (ll.M - ll.N) / K
		}

		// 03.50: Call subroutine 9 (numerical integration)
		ll.subroutine9(K)

		// 03.50: I (I)7.1,7.1;I (V)3.8,3.8;I (J)8.1
		if ll.I <= 0 {
			// 07.10: I (S-.005)5.1 - if S > 0.005 then go to 5.1 (landing)
			if ll.S > 0.005 {
				break
			}
			ll.S = 2 * ll.A / (ll.V + math.Sqrt(ll.V*ll.V+2*ll.A*(ll.G-ll.Z*K/ll.M)))
			ll.subroutine9(K)
			ll.subroutine6(K)
			continue
		}

		if ll.V < 0 {
			// 03.80: D 6;G 3.1 - Call subroutine 6 and continue
			ll.subroutine6(K)
			continue
		}

		if ll.J <= 0 {
			// 08.10: Special velocity handling
			ll.W = (1 - ll.M*ll.G/(ll.Z*K)) / 2
			ll.S = ll.M*ll.V/(ll.Z*K*(ll.W+math.Sqrt(ll.W*ll.W+ll.V/ll.Z))) + 0.05
			ll.subroutine9(K)

			// 08.30: I (I)7.1,7.1;D 6;I (-J)3.1,3.1;I (V)3.1,3.1,8.1
			if ll.I <= 0 {
				continue
			}
			ll.subroutine6(K)
			if ll.J < 0 {
				continue
			}
			if ll.V <= 0 {
				continue
			}
		}

		// Normal case: 03.80: D 6;G 3.1 - Call subroutine 6 and continue
		ll.subroutine6(K)
	}

	impactVelocity := 3600 * ll.V
	fuelLeft := ll.M - ll.N
	return ll.L, impactVelocity, fuelLeft
}

func main() {
	lander := NewLunarLander()
	lander.RunSimulation()
}
