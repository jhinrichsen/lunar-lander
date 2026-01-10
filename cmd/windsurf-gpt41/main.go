package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	initialAltitude = 120.0   // miles
	initialVelocity = 1.0     // mph (relative to moon's surface)
	emptyMass       = 16500.0 // lbs (capsule weight without fuel)
	fuelMass        = 16000.0 // lbs (initial fuel)
	gravity         = 0.001   // miles/sec² (lunar gravity)
	exhaustVelocity = 1.8     // miles/sec (effective exhaust velocity)
)

type Lander struct {
	time        float64               // seconds
	altitude    float64               // miles
	velocity    float64               // miles/second
	mass        float64               // lbs (total mass including fuel)
	fuel        float64               // lbs
	fuelRate    float64               // lbs/sec
	getFuelRate func(*Lander) float64 // Function to get fuel rate (for testing)
}

func NewLander() *Lander {
	l := &Lander{
		altitude: initialAltitude,
		velocity: initialVelocity / 3600, // Convert mph to miles/sec
		mass:     emptyMass + fuelMass,
		fuel:     fuelMass,
	}
	l.getFuelRate = defaultGetFuelRate // Default to the real implementation
	return l
}

func (l *Lander) Update(dt float64) {
	// Save initial mass for calculations
	initialMass := l.mass

	// Calculate fuel consumption
	fuelUsed := l.fuelRate * dt
	if fuelUsed > l.fuel {
		fuelUsed = l.fuel
		l.fuel = 0
		l.fuelRate = 0
		l.mass = emptyMass // Set to empty mass when out of fuel
	} else {
		l.fuel -= fuelUsed
		l.mass = emptyMass + l.fuel
	}

	// If there's no fuel, just apply gravity
	if l.fuel <= 0 && l.fuelRate == 0 {
		l.velocity += gravity * dt
		l.altitude -= l.velocity * dt
		l.time += dt
		if l.altitude < 0 {
			l.altitude = 0
		}
		return
	}

	// Calculate acceleration based on thrust and gravity
	// Thrust = fuelRate * exhaustVelocity
	// Acceleration = Thrust / mass - gravity
	thrust := l.fuelRate * exhaustVelocity
	accel := (thrust / initialMass) - gravity

	// Update velocity and position using kinematic equations
	l.velocity += accel * dt
	l.altitude += l.velocity*dt - 0.5*gravity*dt*dt

	// Update time
	l.time += dt

	// Prevent negative altitude
	if l.altitude < 0 {
		l.altitude = 0
	}
}

func (l *Lander) PrintStatus() {
	miles := math.Floor(l.altitude)
	feet := (l.altitude - miles) * 5280
	mph := l.velocity * 3600 // Convert miles/sec to mph

	fmt.Printf("%10.2f  %12.0f %6.0f  %12.2f  %10.2f  ",
		l.time, miles, feet, mph, l.fuel)
}

// GetFuelRate is a method that gets the fuel rate from the user
func (l *Lander) GetFuelRate() float64 {
	return l.getFuelRate(l)
}

// defaultGetFuelRate is the default implementation of getFuelRate that reads from stdin
func defaultGetFuelRate(l *Lander) float64 {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("      K=")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			return l.fuelRate // Keep current rate if no input
		}

		rate, err := strconv.ParseFloat(input, 64)
		if err != nil || rate < 0 {
			fmt.Println("INVALID INPUT")
			continue
		}

		if rate == 0 || (rate >= 8 && rate <= 200) {
			return rate
		}

		fmt.Println("NOT POSSIBLE" + strings.Repeat(".", 50))
	}
}

func (l *Lander) Land() {
	fmt.Println("CONTROL CALLING LUNAR MODULE. MANUAL CONTROL IS NECESSARY")
	fmt.Println("YOU MAY RESET FUEL RATE K EACH 10 SECS TO 0 OR ANY VALUE")
	fmt.Println("BETWEEN 8 & 200 LBS/SEC. YOU'VE 16000 LBS FUEL. ESTIMATED")
	fmt.Println("FREE FALL IMPACT TIME=120 SECS. CAPSULE WEIGHT=32500 LBS")
	fmt.Println("FIRST RADAR CHECK COMING UP")
	fmt.Println()
	fmt.Println("COMMENCE LANDING PROCEDURE")
	fmt.Println("TIME,SECS   ALTITUDE,MILES+FEET   VELOCITY,MPH   FUEL,LBS   FUEL RATE")

	// Initial status
	l.PrintStatus()
	l.fuelRate = l.GetFuelRate()

	// Main simulation loop
	for l.altitude > 0 {
		// Save state before update
		prevTime := l.time
		prevFuel := l.fuel

		// Update for 10 seconds or until impact
		dt := 10.0
		if l.altitude < 1.0 {
			dt = 0.1 // Smaller steps near impact for better accuracy
		}

		l.Update(dt)

		// Print status at least every 10 seconds or when fuel is out
		if math.Floor(l.time/10) > math.Floor(prevTime/10) ||
			(prevFuel > 0 && l.fuel <= 0) ||
			l.altitude <= 0 {

			l.PrintStatus()
			if l.fuel > 0 {
				l.fuelRate = l.getFuelRate(l)
			}
		}

		// Check for impact
		if l.altitude <= 0 {
			impactVelocity := math.Abs(l.velocity) * 3600 // mph
			fmt.Printf("\nON THE MOON AT %6.2f SECS\n", l.time)
			fmt.Printf("IMPACT VELOCITY OF %8.2f M.P.H.\n", impactVelocity)
			fmt.Printf("FUEL LEFT: %8.2f LBS\n", l.fuel)

			switch {
			case impactVelocity < 1.0:
				fmt.Println("PERFECT LANDING! - (LUCKY)")
			case impactVelocity < 10.0:
				fmt.Println("GOOD LANDING - (COULD BE BETTER)")
			case impactVelocity < 22.0:
				fmt.Println("CONGRATULATIONS ON A POOR LANDING")
			case impactVelocity < 40.0:
				fmt.Println("CRAFT DAMAGE. GOOD LUCK")
			case impactVelocity < 60.0:
				fmt.Println("CRASH LANDING - YOU'VE 5 HRS OXYGEN")
			default:
				craterDepth := impactVelocity * 0.277777 // Convert to ft/s * time factor
				fmt.Println("SORRY, BUT THERE WERE NO SURVIVORS - YOU BLEW IT!")
				fmt.Printf("IN FACT YOU BLASTED A NEW LUNAR CRATER %6.2f FT. DEEP\n", craterDepth)
			}

			// Ask to play again
			fmt.Print("TRY AGAIN? (YES/NO): ")
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(strings.ToUpper(answer))
			if answer == "YES" || answer == "Y" {
				// Reset lander
				l = NewLander()
				fmt.Println("\nFIRST RADAR CHECK COMING UP")
				fmt.Println()
				fmt.Println("COMMENCE LANDING PROCEDURE")
				fmt.Println("TIME,SECS   ALTITUDE,MILES+FEET   VELOCITY,MPH   FUEL,LBS   FUEL RATE")
				l.PrintStatus()
				l.fuelRate = l.GetFuelRate()
			} else {
				fmt.Println("CONTROL OUT")
				return
			}
		}
	}
}

func main() {
	lander := NewLander()
	lander.Land()
}
