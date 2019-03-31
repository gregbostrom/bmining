package hashrate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Coin type
type Coin int

const (
	// Undef - undefined
	Undef Coin = 0

	// XMR - Monero
	XMR Coin = 1
)

// Hashrate returns the global Hashrate value and string
func Hashrate(coin Coin) (float64, string) {

	url := coinInfo(coin)
	if url == "" {
		return 0, ""
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return 0, ""
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return 0, ""
	}

	var dat map[string]interface{}

	err = json.Unmarshal(body, &dat)
	if err != nil {
		fmt.Println(err)
		return 0, ""
	}

	h := dat["hashrate"].(float64)

	// Return Mh/s in value and string
	return h, HumanMhs(h)
}

// Mhs converts hash/sec to Mhash/sec
func Mhs(h float64) float64 {
	return h / 1000000
}

// HumanMhs returns a human readable string for Hhs
func HumanMhs(h float64) string {

	return fmt.Sprintf("%.1f Mh/s\n", Mhs(h))
}

func coinInfo(coin Coin) string {
	switch coin {
	case XMR:
		return "https://moneroblocks.info/api/get_stats"
	}
	return ""
}
