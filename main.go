package main

import (
	"fmt"
	"math"
)

func main() {
	var L int
	var A, V, M, N, G, Z, T, K, S, W, I, J float64

	// Initial setup
	fmt.Println("\t\tLUNAR")
	fmt.Println("\t\tCREATIVE COMPUTING MORRISTOWN, NEW JERSEY")
	fmt.Println()
	fmt.Println("THIS IS A COMPUTER SIMULATION OF AN APOLLO LUNAR")
	fmt.Println("LANDING CAPSULE.")
	fmt.Println()
	fmt.Println("THE ON-BOARD COMPUTER HAS FAILED (IT WAS MADE BY")
	fmt.Println("XEROX) SO YOU HAVE TO LAND THE CAPSULE MANUALLY.")
	fmt.Println()
	fmt.Println("SET BURN RATE OF RETRO ROCKETS TO ANY VALUE BETWEEN")
	fmt.Println("0 (FREE FALL) AND 200 (MAXIMUM BURN) POUNDS PER SECOND.")
	fmt.Println("SET NEW BURN RATE EVERY 10 SECONDS.")
	fmt.Println()
	fmt.Println("CAPSULE WEIGHT 32,500 LBS; FUEL WEIGHT 16,500 LBS.")
	fmt.Println()
	fmt.Println("GOOD LUCK")
	fmt.Println()

	// Initialize variables
	L = 0
	A = 120
	V = 1
	M = 33000
	N = 16500
	G = 1e-03
	Z = 1.8

	// Print header
	fmt.Println("SEC MI + FT MPH LB FUEL BURN RATE")
	fmt.Println()

	// Main loop
	for {
		fmt.Printf("%d %d %d %d %.2f\n", L, int(A), int(5280*(A-math.Floor(A))), int(3600*V), int(M-N))
		fmt.Print("Enter new burn rate: ")
		_, err := fmt.Scanf("%f", &K)
		if err != nil {
			fmt.Println("Invalid input.")
			return
		}

		if M-N < 1e-03 {
			break
		}

		if T < 1e-03 {
			continue
		}

		S = T
		if M >= N+S*K {
			break
		}
		S = (M - N) / K

		// Perform calculations based on conditions
		if I <= O {
			break
		}

		if V <= 0 {
			break
		}

		if J < 0 {
			break
		}

		// Calculate next step
		// Simulate the loop...
	}

	// Final calculations based on speed and conditions
	S = (-V + math.Sqrt(V*V+2*A*G)) / G
	V = V + G*S
	L = L + S

	W = 3600 * V
	fmt.Printf("ON MOON AT %d SECONDS - IMPACT VELOCITY %.2f MPH\n", L, W)

	if W <= 1.2 {
		fmt.Println("PERFECT LANDING!")
	} else if W <= 10 {
		fmt.Println("GOOD LANDING (COULD BE BETTER)")
	} else if W > 60 {
		fmt.Println("SORRY THERE WERE NO SURVIVORS. YOU BLEW IT!")
		fmt.Printf("IN FACT, YOU BLASTED A NEW LUNAR CRATER %.2f FEET DEEP!\n", W*0.227)
	} else {
		fmt.Println("CRAFT DAMAGE... YOU'RE STRANDED HERE UNTIL A RESCUE")
		fmt.Println("PARTY ARRIVES. HOPE YOU HAVE ENOUGH OXYGEN!")
	}

	// Retry option
	fmt.Println("\nTRY AGAIN??")
}
