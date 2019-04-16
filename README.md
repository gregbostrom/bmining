# bmining

Balanced mining (bmining) in terms of what coin to mine with regard to their network hashrates. The goal is maintain an even hashrate distribution across all chains being mined.

Reference:
    "Responsible mining: probabilistic hashrate distribution", 
        Mitchell P. Krawiec-Thayer, other authors ...

## Motivation

Auto-switching mining pools typically follow a financially-greedy algorthim that distributes hashrate based on various token's exchange rates and mining costs. In the case where it is desired to deploy millions of dormant computers/GPUs for mining, using the greedy algorithm will add variablity to the network hashrate and thus reduce the predictability of mining income for other users. The bmining program implements a probabilistic coin selection based on the total network hashrate.


## Notes

Logic:

1. Get the data, the network hashrates, of all the coins of interest.
2. Compute the probabilities to achieve balananced mining.
3. Using a random number select a coin based on the probabilities.

## Directories

**hashrate** - Fetch network hashrate from a reliable source.

## Sample Runs

```
# Help
$ bmining -h
 usage: bmining [OPTION] [list of coins]

   -d      dump network hash rates
   -h      help
   -s [n]  run simulation for [n] trials
   -v      verbose

# Fetch and dump network hashrates
$ bmining -d
X-CASH         (XCASH)   1.5 MH/s    24h: $0.62    CryptoNight HeavyX
Monero         (XMR)   307.8 MH/s    24h: $0.43    CryptoNight R (v4)
Conceal        (CCX)     3.7 MH/s    24h: $0.43    CryptoNight Conceal
Lethean        (LTHN)    2.0 MH/s    24h: $0.42    CryptoNight R (v4)
GRAFT          (GRFT)    8.2 MH/s    24h: $0.32    CryptoNight v2
Masari         (MSR)     5.5 MH/s    24h: $0.29    CryptoNight Fast v2
Ryo Currency   (RYO)     2.2 MH/s    24h: $0.29    CryptoNight Heavy
BitTube        (TUBE)   12.1 MH/s    24h: $0.25    CryptoNight Saber
Haven Protocol (XHV)    20.4 MH/s    24h: $0.18    CryptoNight Haven
Stellite       (XTL)    20.5 MH/s    24h: $0.12    CryptoNight v2
TurtleCoin     (TRTL)  268.0 MH/s    24h: $0.12    CryptoNight-Lite v1
Aeon           (AEON)   25.4 MH/s    24h: $0.08    CryptoNight-Lite v1
Loki           (LOKI)  210.3 MH/s    24h: $0.02    CryptoNight Heavy
BLOC.money     (BLOC)    4.4 MH/s    24h: $0.01    CryptoNight Haven
Dero           (DERO)   80.3 MH/s    24h: $0.01    CryptoNight
Karbo          (KRB)    32.3 MH/s    24h: $0.01    CryptoNight
Electroneum    (ETN)       3 GH/s    24h: $0.01    CryptoNight
Sumokoin       (SUMO)   93.7 MH/s    24h: $0.01    CryptoNight
Bytecoin       (BCN)   519.1 MH/s    24h: $0.00    CryptoNight


# Run a simulation of 1000 random trials to verify the algorithm is working.
# Select coins of interest.  In this example, X-CASH, Monero, Conceal,
# Lethean, and GRAFT are selected.
# First the probabilities, the percents, of each coin is listed.
# Next, the result of the 1000 random trials is shown.

$ bmining -s 1000 xcash xmr ccx lthn grft
Simulation 1000 [XCASH XMR CCX LTHN GRFT]
 XCASH 0.4% XMR 95.5% CCX 1.1% LTHN 0.6% GRFT 2.4%
 XCASH: 5  XMR: 948  CCX: 11  LTHN: 8  GRFT: 28  


# And finally an example on how it is envisioned the program will be used.
# A coin is selected and passed on to the mining command line.

$ COIN=`bmining xcash xmr ccx lthn grft`
$ echo $COIN
XMR

# The $COIN evironmental variable could then be used when starting the 
# the mining program.
```