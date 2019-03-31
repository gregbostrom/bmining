package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// Coin type
type Coin int

const (
	// Undef - undefined
	Undef Coin = 0

	// XMR - Monero
	XMR Coin = 1

	// LOKI - Loki
	LOKI Coin = 2

	// AEON - Aeon
	AEON Coin = 3
)

const gpuHashrate float64 = 1500

var miningXMR int
var miningLOKI int
var miningAEON int

var hXMR float64
var hLOKI float64
var hAEON float64
var hWorld float64

var cacheXMR float64

var xMhs string
var lMhs string
var aMhs string

func updateWorld() {
	hWorld = hXMR + hLOKI + hAEON
}

func incMiningXMR(v int) {
	var m sync.Mutex
	m.Lock()
	miningXMR += v
	hXMR += gpuHashrate
	updateWorld()
	m.Unlock()
}

func incMiningLOKI(v int) {
	var m sync.Mutex
	m.Lock()
	miningLOKI += v
	hLOKI += gpuHashrate
	updateWorld()
	m.Unlock()
}

func incMiningAEON(v int) {
	var m sync.Mutex
	m.Lock()
	miningAEON += v
	hAEON += gpuHashrate
	updateWorld()
	m.Unlock()
}

func loadBalance() {
	max := 100
	xmr := (int)((hXMR * 100) / hWorld)
	loki := (int)((hLOKI * 100) / hWorld)
	//aeon := (int)((hAEON * 100)/ hWorld)
	x := rand.Intn(max)
	if x < xmr {
		incMiningXMR(1)
	} else if x < xmr+loki {
		incMiningLOKI(1)
	} else {
		incMiningAEON(1)
	}
}

// Hashrate returns the global Hashrate value and string
func Hashrate(coin Coin) (float64, string) {

	var Mh float64 = 1000000

	// Hack for the MVP
	if coin == LOKI {
		return 32.82 * Mh, "32.8 Mh/s"
	}

	if coin == AEON {
		return 25.533 * Mh, "25.5 Mh/s"
	}

	if coin == XMR && cacheXMR != 0 {
		return cacheXMR, HumanMhs(cacheXMR)
	}

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
	cacheXMR = h
	return h, HumanMhs(h)
}

// Mhs converts hash/sec to Mhash/sec
func Mhs(h float64) float64 {
	return h / 1000000
}

// HumanMhs returns a human readable string for Hhs
func HumanMhs(h float64) string {

	return fmt.Sprintf("%.1f Mh/s", Mhs(h))
}

func coinInfo(coin Coin) string {
	switch coin {
	case XMR:
		return "https://moneroblocks.info/api/get_stats"
	}
	return ""
}

// Scale these times for the simulation
const scale time.Duration = 10000
const chargeReal time.Duration = time.Duration(30) * time.Minute
const chargeTime time.Duration = chargeReal / scale
const mineReal time.Duration = time.Duration(30) * time.Minute
const mineTime time.Duration = mineReal / scale
const workReal time.Duration = time.Duration(22*60) * time.Minute
const workTime time.Duration = workReal / scale

// StateFunc exported to allow multiple start points
type StateFunc func(*Bminer) StateFunc

// Bminer represents a dynamic balanced miner
type Bminer struct {
	id         int
	startState StateFunc
}

// WG (WaitGroup) to wait for all the miners.
var WG sync.WaitGroup

// New creates a Bminer
func New(start StateFunc, id int) *Bminer {
	return &Bminer{
		id:         id,
		startState: start,
	}
}

func (m *Bminer) Init(start StateFunc) {
	m.startState = start
}

// Start executing the Bminer asynchronously using a goroutine.
func (m *Bminer) Start() {
	m.trace("Start")
	go m.run()
}

// Charge the device
func (m *Bminer) Charge() StateFunc {
	m.trace("Charge")
	time.Sleep(chargeTime)
	return m.Charged()
}

// Charged and idle
func (m *Bminer) Charged() StateFunc {
	m.trace("Charged")
	return m.WorkOrMine()
}

// Mine state
func (m *Bminer) Mine() StateFunc {
	m.trace("Mining")
	loadBalance()
	time.Sleep(mineTime)
	return nil
}

// WorkOrMine are the choices to do next for this device
func (m *Bminer) WorkOrMine() StateFunc {
	if rand.Intn(2) == 0 {
		return m.Work()
	}
	return m.Mine()
}

// Work state
func (m *Bminer) Work() StateFunc {
	m.trace("drone work")
	time.Sleep(workTime)
	return m.Charge()
}

// Private functions

func (m *Bminer) reset() {
	// flush and cleanup
}

func (m *Bminer) trace(s string) {
	//fmt.Printf("id: %d, %s\n", m.id, s)
}

// Drive the Bminer state machine.
func (m *Bminer) run() {

	defer WG.Done()
	state := m.startState
	for state != nil {
		state = state(m)
	}
	m.reset()
}

func logger() {
	const sleepTime time.Duration = time.Duration(3) * time.Second
	for {
		time.Sleep(sleepTime)
		fmt.Printf("Monero hashrate: %s, Loki hashrate: %s, Aeon hashrate: %s\n",
			HumanMhs(hXMR), HumanMhs(hLOKI), HumanMhs(hAEON))
	}
}

func main() {
	fmt.Println("\nbmining MVP\n")
	hXMR, xMhs = Hashrate(XMR)
	hLOKI, lMhs = Hashrate(LOKI)
	hAEON, aMhs = Hashrate(AEON)
	updateWorld()
	fmt.Printf("World hashrate: %s, ", HumanMhs(hWorld))
	fmt.Printf("Monero hashrate: %s, Loki hashrate: %s, Aeon hashrate: %s\n",
		xMhs, lMhs, aMhs)

	rand.Seed(42)
	// rand.Seed(time.Now().UnixNano())

	for id := 1; id < 100000; id++ {
		m := New((*Bminer).Charge, id)
		WG.Add(1)
		m.Start()
	}
	const sleepTime time.Duration = time.Duration(5) * time.Second
	for i := 0; i < 13; i++ {
		time.Sleep(sleepTime)
		//fmt.Printf("Monero hashrate: %s, Loki hashrate: %s, Aeon hashrate: %s\n",
		//	HumanMhs(hXMR), HumanMhs(hLOKI), HumanMhs(hAEON))
		fmt.Printf("                            Monero miners: %d,        Loki miners: %d,        Aeon miners: %d\n",
			miningXMR, miningLOKI, miningAEON)
	}

	fmt.Printf("World hashrate: %s, ", HumanMhs(hWorld))
	fmt.Printf("Monero hashrate: %s, Loki hashrate: %s, Aeon hashrate: %s\n",
		HumanMhs(hXMR), HumanMhs(hLOKI), HumanMhs(hAEON))

	WG.Wait()
	fmt.Println("All done")

}

// World hashrate: 356.5 Mh/s, Monero hashrate: 298.1 Mh/s, Loki hashrate: 32.8 Mh/s, Aeon hashrate: 25.5 Mh/s
// 1234567890123456789012345678                     12345678                  12345678
//                             Monero miners: 41228,        Loki miners: 4456,        Aeon miners: 4047
