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

   Select a coin from the list using probabilistic hashrate distribution
       based on the total network hashrate.

   -d      dump network hash rates
   -h      help
   -s [n]  run simulation for [n] trials
   -v      verbose

# Fetch and dump price, network hashrate, 24hr mining reward, and algorithm
$ bmining -d
X-CASH         (XCASH) $0.000040    1.7 MH/s   24h: $0.64   CryptoNight HeavyX
Conceal        (CCX)   $0.17        3.4 MH/s   24h: $0.49   CryptoNight Conceal
Monero         (XMR)   $68.20     314.2 MH/s   24h: $0.43   CryptoNight R (v4)
Lethean        (LTHN)  $0.0015      2.1 MH/s   24h: $0.42   CryptoNight R (v4)
Masari         (MSR)   $0.14        4.2 MH/s   24h: $0.34   CryptoNight Fast v2
Ryo Currency   (RYO)   $0.0643      2.1 MH/s   24h: $0.31   CryptoNight Heavy
Haven Protocol (XHV)   $0.45       13.5 MH/s   24h: $0.28   CryptoNight Haven
BitTube        (TUBE)  $0.0471     14.6 MH/s   24h: $0.23   CryptoNight Saber
Stellite       (XTL)   $0.0001     14.3 MH/s   24h: $0.12   CryptoNight v2
TurtleCoin     (TRTL)  $0.0002    304.1 MH/s   24h: $0.10   CryptoNight-Lite v1
Aeon           (AEON)  $0.38       23.9 MH/s   24h: $0.09   CryptoNight-Lite v1
Loki           (LOKI)  $0.25      236.5 MH/s   24h: $0.01   CryptoNight Heavy
BLOC.money     (BLOC)  $0.0082      4.6 MH/s   24h: $0.01   CryptoNight Haven
Sumokoin       (SUMO)  $0.0404     68.1 MH/s   24h: $0.01   CryptoNight
Dero           (DERO)  $1.04       77.5 MH/s   24h: $0.01   CryptoNight
Karbo          (KRB)   $0.0872     30.8 MH/s   24h: $0.01   CryptoNight
Electroneum    (ETN)   $0.0056        3 GH/s   24h: $0.01   CryptoNight
Bytecoin       (BCN)   $0.0009    398.6 MH/s   24h: $0.00   CryptoNight


# Run a simulation of 1000 random trials to verify the algorithm is working.
# In this example, X-CASH, Lethean, Conceal, and Monero are selected.
# First the probabilities, the percents, of each coin is listed.
# Next, the result of the 1000 random trials is shown.

$ bmining -s 1000  xcash lthn ccx xmr
Simulation 1000 [XCASH LTHN CCX XMR]
 XCASH 0.5% LTHN 0.6% CCX 1.1% XMR 97.7%
 XCASH: 10  LTHN: 5  CCX: 11  XMR: 974  


# And finally an example on how it is envisioned the program will be used.
# A coin is selected and passed on to the mining command line.

$ COIN=`bmining xcash lthn ccx xmr`
$ echo $COIN
XMR

# The $COIN evironmental variable would then be used when starting the 
# the mining program.
```