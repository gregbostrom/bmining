package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/gregbostrom/bmining/hashrate"
)

func dumpNetHashes(coins []string) {
	dumpCoins := hashrate.DumpCoinHash(coins)
	for i := 0; i < len(dumpCoins); i++ {
		fmt.Println(dumpCoins[i])
	}
}

// Select a coin from an array of coins and their associated
// device hashrates (h) and their network hashrates (H)
func selectCoin(coins []string, h []float64, H []float64, v bool) (string, error) {

	count := len(coins)

	if count != len(h) || count != len(H) {
		return "", errors.New("Number of coins, h, H must all be equal")
	}

	/*
	 * Reference:
	 *    Mitchell P. Krawiec-Thayer, other authors ...
	 *        Responsible mining: probabilistic hashrate distribution
	 */
	x := make([]float64, count)

	var i int

	for i = 0; i < count; i++ {
		x[i] = h[i] / H[i]
	}

	var Z float64

	for i = 0; i < count; i++ {
		Z += (1 / x[i])
	}

	P := make([]float64, count)

	for i = 0; i < count; i++ {
		P[i] = (1 / x[i]) * (1 / Z)
	}

	// Verbose
	if v == true {
		vs := ""
		for i = 0; i < count; i++ {
			percent := P[i] * 100.0
			vs += fmt.Sprintf(" %s %.1f%%", coins[i], percent)
		}
		fmt.Println(vs)
	}

	// Verify all the probabiliites add up to 1.
	var sum float64
	for i = 0; i < count; i++ {
		sum += P[i]
	}
	diff := sum - 1.0
	if diff > 0.01 || diff < -0.01 {
		s := fmt.Sprintf("Probability sum %.4f != 1", sum)
		return "", errors.New(s)
	}

	// pseudo-random number in [0.0,1.0)
	r := rand.Float64()

	sum = 0.0

	for i = 0; i < count; i++ {
		if (sum + P[i]) > r {
			break
		}
		sum += P[i]
	}

	return coins[i], nil
}

func simulation(count int, coins []string, h []float64, H []float64) {
	fmt.Println("Simulation", count, coins)

	tally := make(map[string]int)

	v := true
	for i := 0; i < count; i++ {
		coin, err := selectCoin(coins, h, H, v)
		if err != nil {
			fmt.Println(err)
			return
		}
		v = false
		tally[coin]++
	}

	// Dump out the results of the simulation
	s := " "
	for i := 0; i < len(coins); i++ {
		coin := coins[i]
		s += fmt.Sprintf("%s: %d  ", coin, tally[coin])
	}
	fmt.Println(s)
}

func main() {

	rand.Seed(time.Now().UnixNano())

	const defaultHR float64 = 500

	dmpf := flag.Bool("d", false, "dump coins supported and their network hashrate")
	hash := flag.Float64("hr", defaultHR, "hashrate of the device")
	help := flag.Bool("h", false, "help")
	simu := flag.Int("s", 0, "run simulation [n] times")
	//topc := flag.Int("t", 0, "select the the top [n] coins - USD 24h mining rewards")
	verb := flag.Bool("v", false, "Verbose")

	flag.Parse()

	diag := fmt.Sprintf("dev hashrate: %.0f Hs; Coins: ", *hash)

	coins := flag.Args()
	count := len(coins)

	dump := *dmpf // So we may modify dump for debugging
	// For debugging:
	// dump = true

	if *help == true || (count == 0 && dump == false) {
		fmt.Println(" usage: bmining [OPTION] [list of coins]\n")
		fmt.Println("   Select a coin from the list using probabilistic hashrate distribution")
		fmt.Println("       based on the total network hashrate.")
		fmt.Println("")
		fmt.Println("   -d      dump network hash rates")
		fmt.Println("   -h      help")
		fmt.Println("   -s [n]  run simulation for [n] trials")
		//fmt.Println("   -t [n]  select the top [n] coins")
		fmt.Println("   -v      verbose")
		fmt.Println("")
		return
	}

	hashrate.InitCoinHash(*verb)

	// h is our device hashrate for the coin
	// H is the network hashrate for the coin
	h := make([]float64, count)
	H := make([]float64, count)

	for i := 0; i < count; i++ {
		coins[i] = strings.ToUpper(coins[i])
		coin := hashrate.LookupCoin(coins[i])
		if coin == nil {
			fmt.Println("Unknown coin", coins[i])
			return
		}
		h[i] = *hash
		H[i] = coin.NetHashRate

		diag += coins[i]
		diag += " "
	}

	if dump == true {
		dumpNetHashes(coins)
		return
	}

	// Populate the two maps:
	//    h is the hashrate of our device for this coin
	//    H is the network hashrate for this coin
	for i := 0; i < count; i++ {
		coins[i] = strings.ToUpper(coins[i])
		coin := hashrate.LookupCoin(coins[i])
		if coin == nil {
			fmt.Println("Unknown coin", coins[i])
			return
		}
		h[i] = *hash
		H[i] = coin.NetHashRate

		diag += coins[i]
		diag += " "
	}

	if *verb == true {
		fmt.Println(diag)
	}

	if *simu != 0 {
		simulation(*simu, coins, h, H)
		return
	}

	coin, err := selectCoin(coins, h, H, false)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(coin)
}
