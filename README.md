# bmining

Balanced mininig.  Balanced in terms of what coin to mine. The goal is maintain an even hashrate distribution across all chains being mined.

Reference:
    "Responsible mining: probabilistic hashrate distribution", 
        Mitchell P. Krawiec-Thayer, other authors ...

## Glossary

**bminer** - Balanced miner implemented in a distributed fashion.

## Notes

Logic:

1. Get the data, the network hashrates, of all the coins of interest.
2. Compute the probabilities to achieve balananced mining.
3. Using a random number select a coin based on the probabilities.

## Directories

**hashrate** - Fetch network hashrate from a reliable source.

**mvp.Mar29** - Minimal Viable Product (MVP) to demo after a one week sprint.

## Sample Runs

```
# Fetch and dump network hashrates
$ bmining -d
Aeon           (AEON)   25.4 MH/s
BLOC.money     (BLOC)    4.0 MH/s
BitTube        (TUBE)   15.2 MH/s
Bytecoin       (BCN)   513.4 MH/s
Conceal        (CCX)     2.8 MH/s
Dero           (DERO)   71.7 MH/s
Electroneum    (ETN)       3 GH/s
GRAFT          (GRFT)    8.7 MH/s
Haven Protocol (XHV)    16.2 MH/s
Karbo          (KRB)    24.2 MH/s
Lethean        (LTHN)    2.1 MH/s
Loki           (LOKI)  237.7 MH/s
Masari         (MSR)     4.4 MH/s
Monero         (XMR)   295.7 MH/s
Ryo Currency   (RYO)     1.6 MH/s
Stellite       (XTL)    10.7 MH/s
Sumokoin       (SUMO)   70.0 MH/s
TurtleCoin     (TRTL)  292.9 MH/s
UltraNote      (XUN)   464.0 KH/s
Webchain       (WEB)   737.1 KH/s
X-CASH         (XCASH)   1.7 MH/s

# Run a simulation of 1000 random trials to verify the algorithm is working.
# Select coins of interest.  In this example, Monero, Loki, and Aeon.
# First the probabilities, the percents, of each coin is listed.
# Next, the result of the 1000 random trials is shown.

$ bmining -s 1000 xmr loki aeon
Simulation 1000 [XMR LOKI AEON]
 XMR 52.2% LOKI 43.3% AEON 4.5%
 XMR: 511  LOKI: 444  AEON: 45

 # And finally an example on how it is envisioned the program will be used.
 # A coin is selected and passed on to the mining command line.

$ COIN=`bmining xmr loki aeon`
$ echo $COIN
XMR
```
