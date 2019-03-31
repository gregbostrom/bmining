package bminer

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Scale these times for the simulation
const scale time.Duration = 100
const chargeReal time.Duration = time.Duration(30) * time.Minute
const chargeTime time.Duration = chargeReal / scale
const mineReal time.Duration = time.Duration(30) * time.Minute
const mineTime time.Duration = mineReal / scale
const workReal time.Duration = time.Duration(30) * time.Minute
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
	fmt.Printf("id: %d, %s\n", m.id, s)
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
