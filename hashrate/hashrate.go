package hashrate

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const KHs float64 = 1000
const MHs float64 = 1000000
const GHs float64 = 1000000000
const PHs float64 = 1000000000000
const THs float64 = 1000000000000000

type Coin struct {
	Rank        int
	Symbol      string
	Name        string
	Algorithm   string
	XmrStakAlgo string
	NetHashRate float64
	USD24h      float64
}

var Coins []*Coin

// InitCoinHash will initialize network hashrates of coins.
func InitCoinHash() {
	scrapeMineCryptoNight()
}

func scrapeCoinSection(n int, s string) (*Coin, string) {

	var err error
	var f float64

	c := new(Coin)
	c.Rank = n
	// Pointing at the name
	chunks := strings.SplitN(s, "</span>", 2)
	c.Name = chunks[0]
	// Followed by the symbol
	chunks = strings.SplitN(chunks[1], "> (", 2)
	chunks = strings.SplitN(chunks[1], ")</span>", 2)
	c.Symbol = chunks[0]
	// Reward USD for 24hr
	chunks = strings.SplitN(chunks[1], "24h: $", 2)
	chunks = strings.SplitN(chunks[1], "<", 2)
	c.USD24h, err = strconv.ParseFloat(chunks[0], 64)
	// Algorithm
	chunks = strings.SplitN(chunks[1], "Algorithm: ", 2)
	chunks = strings.SplitN(chunks[1], "<", 2)
	c.Algorithm = chunks[0]
	// Network hash rate
	chunks = strings.SplitN(chunks[1], "Network hash rate: ", 2)
	hashUnits := strings.SplitN(chunks[1], " ", 2)
	// f will be the MH/s (or other hash units)
	f, err = strconv.ParseFloat(hashUnits[0], 64)
	if err != nil {
		fmt.Println(err)
		// Try to recover
		return nil, hashUnits[1]
	}
	// Verify and convert the hash units
	verify := strings.SplitN(hashUnits[1], "<", 2)
	if verify[0] == "KH/s" || verify[0] == "kH/s" {
		f *= KHs
	} else if verify[0] == "MH/s" {
		f *= MHs
	} else if verify[0] == "GH/s" {
		f *= GHs
	} else if verify[0] == "PH/s" {
		f *= PHs
	} else if verify[0] == "TH/s" {
		f *= THs
	} else {
		fmt.Println("Unknown hashrate unit ", verify[0], c.Name)
		// Try to recover
		return nil, verify[1]
	}
	c.NetHashRate = f

	return c, verify[1]
}

func scrapeMineCryptoNight() {

	url := "https://minecryptonight.net/"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("http.Get failed for url", url, err)
		return
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll failed", err)
		return
	}
	scrape := string(bytes)

	// Parse through at most 21 sections
	const maxCoins int = 21
	coins := make([]*Coin, maxCoins)
	// icoins indexs the actual coins parse successfully
	icoins := 0
	// n counts the sections
	n := 1
	for ; n < maxCoins+1; n++ {
		if scrape == "" {
			break
		}
		x := "<span>"
		if n < 10 {
			x += "[1-9]"
		} else if n < 20 {
			x += "1[0-9]"
		} else {
			x += "2[0-9]"
		}
		x += "\\. "
		chunks := regexp.MustCompile(x).Split(scrape, 2)
		if len(chunks) != 2 {
			break
		}
		coins[icoins], scrape = scrapeCoinSection(n, chunks[1])
		// Count it if we got one.
		if coins[icoins] != nil {
			icoins++
		}
	}

	if n == 1 {
		// Mishap
		fmt.Println("No coins found.")
		return
	}

	Coins = make([]*Coin, icoins)
	copy(Coins, coins[:icoins])
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func DumpCoinHash(coins []string) []string {

	const namePad int = 15
	const symbPad int = 23

	const humanPad int = 10
	const usd24Pad int = 30

	cnt := len(coins)
	all := (cnt == 0)

	if all == true {
		cnt = len(Coins)
	}

	dump := make([]string, cnt)
	i := 0

	for _, coin := range Coins {
		// Skip if not a coin of interest
		if !(all || contains(coins, coin.Symbol)) {
			continue
		}
		s := coin.Name
		for len(s) < namePad {
			s += " "
		}
		s = s + "(" + coin.Symbol + ") "
		for len(s) < symbPad {
			s += " "
		}
		// Pad to a length of 10
		h := HumanHs(coin.NetHashRate)
		for len(h) < humanPad {
			h = " " + h
		}
		s += h

		usd24 := fmt.Sprintf("    24h: $%.2f", coin.USD24h)
		s += usd24
		s += "    "
		s += coin.Algorithm
		dump[i] = s
		i++
	}

	//sort.Strings(dump)

	return dump
}

// HumanHs will display the Hs in human readable form.
func HumanHs(f float64) string {
	var s string

	if f < KHs {
		s = fmt.Sprintf("%f H/s", f)
	} else if f < MHs {
		s = fmt.Sprintf("%.1f KH/s", f/KHs)
	} else if f < GHs {
		s = fmt.Sprintf("%.1f MH/s", f/MHs)
	} else if f < PHs {
		s = fmt.Sprintf("%1.f GH/s", f/GHs)
	} else if f < THs {
		s = fmt.Sprintf("%1.f PH/s", f/PHs)
	} else {
		s = fmt.Sprintf("%.1f TH/s", f/THs)
	}

	return s
}

// LookupCoin returns pointer to the coin
func LookupCoin(symb string) *Coin {
	for _, coin := range Coins {
		if symb == coin.Symbol {
			return coin
		}
	}
	return nil
}
