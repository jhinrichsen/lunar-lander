# Prompt for Evaluating AI IDE/LLM/Agents on the 1969 FOCAL Lunar Lander Simulation

You are tasked with evaluating AI models (IDEs, LLMs, agents) based on their capability to handle, translate, and optimize the 1969 FOCAL Lunar Lander simulation.

##  Objective

This evaluation has three distinct targets:

1. Understanding the Physics (Target 1)

Clearly demonstrate the AI’s ability to comprehend the underlying physics involved in the Lunar Lander scenario (gravity, fuel consumption, velocity, thrust).

Definition of Done (DoD): A clearly structured README.md file including a detailed and accurate "Physics Section" that explains the equations, variables, assumptions, and physical laws used in the original simulation.

2. Transpilation from FOCAL to Go 1.24 (Target 2)

Accurately transpile the original 1969 FOCAL source code of the Lunar Lander to modern Go 1.24, fully utilizing the latest features (generics, slices, etc.).

Take special care regarding floating-point arithmetic differences between FOCAL and Go, ensuring high fidelity replication of the original logic and outcomes.

Provide a detailed explanation of any differences or adjustments required due to floating-point arithmetic disparities.

Definition of Done (DoD): A runnable and fully tested Go file named main.go, capable of producing simulation outputs that match the behavior and numeric results of the original Lunar Lander implementation as closely as possible, within the constraints of Go’s standard float64 or float128 arithmetic.

Optional Recommendation: Extra recognition will be given for implementations using a custom numeric type closely mimicking the original PDP-8/FOCAL floating-point behavior for enhanced historical accuracy.

3. Self-Optimizing AI Solutions (Target 3)

Develop a self-optimizing implementation that intelligently adjusts its own inputs to achieve:

a) The softest possible landing.

b) The most fuel-efficient landing.

Implement the optimization strategies as two distinct and clearly defined test functions, ensuring robustness and repeatability of results.

Clearly outline the AI's optimization strategy, whether it leverages reinforcement learning, genetic algorithms, Bayesian optimization, or other self-improvement techniques.

Definition of Done (DoD): Two test functions (TestSoftLanding and TestFuelEfficientLanding) that clearly demonstrate the optimization approaches and provide reproducible evaluation metrics for validation.

## Prerequisites

Ensure the following prerequisites for the implementation:

Go 1.24

Optional: Provide a Makefile if non-Go dependencies are required

Source environment: Fedora 42 DNF5

Linux support required (no considerations for macOS or Windows; Go's cross-platform nature should inherently support this requirement)

Deliverables Summary

Target 1: README.md with Physics Section

Target 2: Fully runnable and tested Go file (main.go)

Target 3: Two clearly defined test functions (TestSoftLanding and TestFuelEfficientLanding) demonstrating robust and repeatable optimization outcomes.

Submit each target separately with detailed documentation for easy validation and reproducibility of your results.
##  FOCAL Source Code and Examples

Original FOCAL Source Code

```FOCAL
01.04 T "CONTROL CALLING LUNAR MODULE. MANUAL CONTROL IS NECESSARY"!
01.06 T "YOU MAY RESET FUEL RATE K EACH 10 SECS TO 0 OR ANY VALUE"!
01.08 T "BETWEEN 8 & 200 LBS/SEC. YOU'VE 16000 LBS FUEL. ESTIMATED"!
01.11 T "FREE FALL IMPACT TIME=120 SECS. CAPSULE WEIGHT=32500 LBS"!
01.20 T "FIRST RADAR CHECK COMING UP"!!!;E
01.30 T "COMMENCE LANDING PROCEDURE"!"TIME,SECS   ALTITUDE,"
01.40 T "MILES+FEET   VELOCITY,MPH   FUEL,LBS   FUEL RATE"!
01.50 S A=120;S V=1;S M=32500;S N=16500;S G=.001;S Z=1.8
02.10 T "    "%3,L,"           "FITR(A),"  "%4,5280*(A-FITR(A))
02.20 T %6.02,"       "3600*V,"    "%6.01,M-N,"      K=";A K;S T=10
02.70 T %7.02;I (200-K)2.72;I (8-K)3.1,3.1;I (K)2.72,3.1
02.72 T "NOT POSSIBLE";F X=1,51;T "."
02.73 T "K=";A K;G 2.7
03.10 I (M-N-.001)4.1;I (T-.001)2.1;S S=T
03.40 I (N+S*K-M)3.5,3.5;S S=(M-N)/K
03.50 D 9;I (I)7.1,7.1;I (V)3.8,3.8;I (J)8.1
03.80 D 6;G 3.1
04.10 T "FUEL OUT AT"L, " SECS"!
04.40 S S=(FSQT(V*V+2*A*G)-V)/G;S V=V+G*S;S L=L+S
05.10 T "ON THE MOON AT"L, " SECS"!;S W=3600*V
05.20 T "IMPACT VELOCITY OF"W, " M.P.H."!"FUEL LEFT:"M-N, " LBS"!
05.40 I (1-W)5.5,5.5;T "PERFECT LANDING !-(LUCKY)"!;G 5.9
05.50 I (10-W)5.6,5.6;T "GOOD LANDING-(COULD BE BETTER)";G 5.9
05.60 I (22-W)5.7,5.7;T "CONGRATULATIONS ON A POOR LANDING";G 5.9
05.70 I (40-W)5.81,5.81;T "CRAFT DAMAGE. GOOD LUCK";G 5.9
05.81 I (60-W)5.82,5.82;T "CRASH LANDING-YOU'VE 5 HRS OXYGEN";G 5.9
05.82 T "SORRY,BUT THERE WERE NO SURVIVORS-YOU BLEW IT!"!"IN "
05.83 T "FACT YOU BLASTED A NEW LUNAR CRATER",W*.277777," FT. DEEP"!
05.90 T !!!!"TRY AGAIN?"!
05.92 A "(ANS. YES OR NO)"P;I (P-0NO)5.94,5.98
05.94 I (P-0YES)5.92,1.2,5.92
05.98 T "CONTROL OUT"!!!;Q
06.10 S L=L+S;S T=T-S;S M=M-S*K;S A=I;S V=J
07.10 I (S-.005)5.1;S S=2*A/(V+FSQT(V*V+2*A*(G-Z*K/M)))
07.30 D 9;D 6;G 7.1
08.10 S W=(1-M*G/Z*K)/2;S S=M*V/(Z*K*(W+FSQT(W*W+V/Z)))+.05;D 9
08.30 I (I)7.1,7.1;D 6;I (-J)3.1,3.1;I (V)3.1,3.1,8.1
09.10 S Q=S*K/M;S J=V+G*S+Z*(-Q-Q^2/2-Q^3/3-Q^4/4-Q^5/5)
09.40 S I=A-G*S*S/2-V*S+Z*S*(Q/2+Q^2/6+Q^3/12+Q^4/20+Q^5/30)
```

## Example Outputs

It is safe to ignore the TRY AGAIN dialog, focus on values and proper formating.

### A bad crash landing typically looks like this sample output.

```
 G
CONTROL CALLING LUNAR MODULE. MANUAL CONTROL IS NECESSARY
YOU MAY RESET FUEL RATE K EACH 10 SECS TO 0 OR ANY VALUE
BETWEEN 8 & 200 LBS/SEC. YOU'VE 16000 LBS FUEL. ESTIMATED
FREE FALL IMPACT TIME-120 SECS. CAPSULE WEIGHT-32500 LBS
FIRST RADAR CHECK COMING UP


COMMENCE LANDING PROCEDURE
TIME,SECS   ALTITUDE,MILES+FEET   VELOCITY,MPH   FUEL,LBS   FUEL RATE
       0            120      0        3660.80     16000.0       K=:0
      10            109   5016        3636.88     16000.0       K=:0
      20             99   4224        3672 .80    16000.0       K=:0
      30             89   2904        3788 .88    16000.0       K=:0
      40             79   1056        3744.80     16000.0       K=:0
      50             68   3960         3788.88    16600.0       K=:0
      60             58   1056        3816.88     16800.0       K=:0
      70             47   2904        3852.80     16680.0       K=:170
      80             37   1474        3539.86     14308.0       K=:200
      90             27   S247        3148.80     12306.0       K=:200
     100             19   4537        2716.41     10300.0       K=:200
     110             12   5118        2243.83      8300.0       K=:200
     120              7   2284        1734.97      6380.0       K=:200
     130              3   1990        1176.06      4300.0       K=:200
     140              0   5040         556.96      2308.0       K=:190
     150              0   1581      -   97.44       400.0       K=:0
     160              0   2746      -   61.44       400.0       K=:0
     170              0   3383      -   25.44       400.0       K=:0
     180              0   3492          10.56       400.0       K=:0
     190              0   3073          46.56       400.0       K=:0
     200              0   2126          82.56       400.0       K=:0
     210              0    652         118.56       400.0       K=:20
ON THE MOON AT   214.03 SECS
IMPACT VELOCITY OF   102.180 M.P.H.
FUEL LEFT:   319.47 LBS
SORRY,BUT THERE WERE NO SURVIVORS-YOU BLEW IT!
IN FACT YOU BLASTED A NEW LUNAR CRATER    28.36 FT. DEEP
(ANS. YES OR NO):YES
FIRST RADAR CHECK COMING UP

```

### A good landing may look like this sample output.
Remember there is an multitude of safe landing alternatives.

```
COMMENCE LANDING PROCEDURE
TIME, SECS  ALTITUDE,MILES+FEET   VELOCITY,MPH   FUEL,LBS   FUEL RATE
        0           120      0        3600.00     16000.0      K=:0
       10           109   5016        3636.00     16006.0      K=:0
       20            99   4224        3672.00     16000.0      K=:0
       30            89   2904        3788.00     16008.0      K=:0
       40            79   1056        3744.00     16000.0      K=:0
       50            68   3960        3780.00     16000.0      K=:0
       60            58   1056        3816.00     16000.0      K=:0
       70            47   2904        3852.00     16000.0      K=:170
       80            37   1474        3539.00     14300.0      K=:200
       90            27   5247        3146.00     12300.0      K=:200
      100            19   4537        2710.00     10300.0      K=:200
      110            12   5118        2243.00      8300.0      K=:200
      120             7   2284        1734.97      6300.0      K=:200
      130             3   1990        1176.86      4300.0      K=:200
      140             0   5040         556.96      2300.0      K=:170
      150             0   1040      -   21.21       600.0      K=:0
      160             0   1087          14.79       680.0      K=:0
      170             0    606          50.79       600.0      K=:30
      180             0    436      -   27.90       300.0      K=:0
      190             0    581           8.10       300.0      K=:8
      200             0    425          13.17       220.0      K=:10
      210             0    253          10.30       120.0      K=:9
      220             0     96          11.11        30.0      K=:100
FUEL OUT AT   220.30 SECS
ON THE MOON AT   226.12 SECS
IMPACT VELOCITY OF    21.36 M.P.H.
FUEL LEFT:     0.00 LBS
CONGRATULATIONS ON A POOR LANDING



TRY AGAIN?
(ANS. YES OR NO):N0
(ANSe YES OR NO):NO
CONTROL OUT
```


