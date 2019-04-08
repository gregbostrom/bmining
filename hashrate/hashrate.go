package hashrate

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

const KHs float64 = 1000
const MHs float64 = 1000000
const GHs float64 = 1000000000
const PHs float64 = 1000000000000
const THs float64 = 1000000000000000

// CoinHash maps a coin's symbol to its current network hashrate.
var CoinHash = make(map[string]float64)

// CoinName maps a coin's symbol to the coin's name.
var CoinName = make(map[string]string)

// CoinSymb maps a coin's name to the coin's symbol.
var CoinSymb = make(map[string]string)

// PreInit will initialize for command line verification.
func PreInit() {
	initCoinName()
	initCoinSymb()
}

// InitCoinHash will initialize network hashrates of coins.
func InitCoinHash() {
	scrapeMineCryptoNight()
}

func initCoinName() {
	CoinName["AEON"] = "Aeon"
	CoinName["BCN"] = "Bytecoin"
	CoinName["BLOC"] = "BLOC.money"
	CoinName["DERO"] = "Dero"
	CoinName["ETN"] = "Electroneum"
	CoinName["IRD"] = "Iridium"
	CoinName["KRB"] = "Karbo"
	CoinName["LOKI"] = "Loki"
	CoinName["SUMO"] = "Sumokoin"
	CoinName["TRTL"] = "TurtleCoin"
	CoinName["TUBE"] = "BitTube"
	CoinName["XHV"] = "Haven Protocol"
	CoinName["XTL"] = "Stellite"
	CoinName["CCX"] = "Conceal"
	CoinName["GRFT"] = "GRAFT"
	CoinName["LTHN"] = "Lethean"
	CoinName["MSR"] = "Masari"
	CoinName["XMR"] = "Monero"
	CoinName["RYO"] = "Ryo Currency"
	CoinName["XUN"] = "UltraNote"
	CoinName["WEB"] = "Webchain"
	CoinName["XCASH"] = "X-CASH"
}

func initCoinSymb() {
	// assume InitCoinName() has been called
	for symb, name := range CoinName {
		CoinSymb[name] = symb
	}
}

func scrapeMineCryptoNight() {

	url := "https://minecryptonight.net/"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("http.Get failed", err)
		return
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll failed", err)
		return
	}
	scrape := string(bytes)

	/*
			*  Scrape the 'Network hash rate' and the coin symbol.
			*  The text is of the form:
			*
		    *  Algorithm: CryptoNight-Lite v1
		    *  Difficulty: 6,633,307,490
		    *  Network hash rate: 27.64 MH/s
			*  Last block reward: 6.413 AEON
			*
			*  In this case the coin is AEON and the netHashRate is 27.64 MH/s
	*/

	dif := "Difficulty: "
	nsr := "Network hash rate: "
	lbr := "Last block reward: "
	chunks := strings.SplitAfterN(scrape, dif, 2)

	for len(chunks) == 2 {
		// Find the "Network hash rate"
		chunks = strings.SplitAfterN(chunks[1], nsr, 2)
		hashUnits := strings.SplitN(chunks[1], " ", 2)
		// f will be the MH/s (or other hash units)
		f, err := strconv.ParseFloat(hashUnits[0], 64)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Verify the hash units
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
			fmt.Println("Unknown hashrate unit ", verify[0])
			return
		}

		// Now scrape the coin symbol
		chunks = strings.SplitAfterN(hashUnits[1], lbr, 2)
		chunks = strings.SplitAfterN(chunks[1], " ", 2)
		chunks = strings.SplitN(chunks[1], "<", 2)
		coin := chunks[0]
		// Map the coin to its network hashrate
		CoinHash[coin] = f
		// Continue for the next coin
		chunks = strings.SplitAfterN(chunks[1], dif, 2)
	}
}

func DumpCoinHash() []string {

	const humanPad int = 10
	const namePad int = 15
	const symbPad int = 23

	dump := make([]string, len(CoinHash))
	i := 0

	for k, v := range CoinHash {
		s := CoinName[k]
		for len(s) < namePad {
			s += " "
		}
		s = s + "(" + k + ") "
		for len(s) < symbPad {
			s += " "
		}
		// Pad to a length of 10
		h := HumanHs(v)
		for len(h) < humanPad {
			h = " " + h
		}
		s += h
		dump[i] = s
		i++
	}

	sort.Strings(dump)

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

// VerifyCoin returns true if the coin symbol is supported.
func VerifyCoin(symb string) bool {
	return CoinName[symb] != ""
}
