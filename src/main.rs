
// Allow all uppercase variable names
#![allow(non_snake_case)]
#![allow(unused_variables)]

use std::io::{self,BufRead};

fn intro() {
    println!("{:^60}", "LUNAR");
    println!("{:^60}", "CREATIVE COMPUTING MORRISTOWN, NEW JERSEY");
    println!("\n\n");
    println!("THIS IS A COMPUTER SIMULATION OF AN APOLLO LUNAR");
    println!("LANDING CAPSULE.");
    println!("\n");
    println!("THE ON-BOARD COMPUTER HAS FAILED (IT WAS MADE BY");
    println!("XEROX) SO YOU HAVE TO LAND THE CAPSULE MANUALLY.");
    println!("");
    println!("SET BURN RATE OF RETRO ROCKETS TO ANY VALUE BETWEEN");
    println!("0 (FREE FALL) AND 200 (MAXIMUM BURN) POUNDS PER SECOND.");
    println!("SET NEW BURN RATE EVERY 10 SECONDS.");
    println!("");
    println!("CAPSULE WEIGHT 32,500 LBS; FUEL WEIGHT 16,500 LBS.");
    println!("\n\n");
    println!("GOOD LUCK");
}

fn main() {
    intro();
    let L = 0;
    println!("{:>10} {:>10} {:>10} {:>10} {:>10}", "SEC", "MI + FT", "MPH", "LB FUEL", "BURN RATE\n");

    let A = 120.0;
    let V = 1;
    let M = 33000.0;
    let N = 16500.0;
    let G = 1E-03;
    let Z = 1.8;

    println!("{:>10} {:>10} {:>10} {:>10} {:>10}", L, A, (5280 * ( A - ((A as i64) as f64)) as i64), 3600 * V, M - N);

    //let K = io::stdin().read_line().ok().expect("Failed to read line");
    let sin = io::stdin();
    let s = sin.lock().lines().next().unwrap().unwrap();
    let input: Option<f64> = s.trim().parse::<f64>().ok();
    let K = match input {
        Some(f) => f,
        None    => {
            println!("please input a number");
            return;
        }
    };
    let T = 10.0;
    if M - N < 1E-03 {
        // 240
    } else {
        if T < 1E-03 {
            // 150
        } else {
            let mut S = T; // 180
            if M >= N + S * K {
                // 200
            } else {
                S = (M - N) / K; // 190
            }
            // 200
            f420();
        }
    }
}

fn f420() {
}

// EOF
