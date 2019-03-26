package bminer

type StateFunc func(*Bminer) StateFunc

// Bminer represents a dynamic balanced miner
type Bminer struct {
	startState StateFunc
}

// New creates a Bminer
func New(start StateFunc) *Bminer {
	return &Bminer{
		startState: start,
	}
}

// Start executing the Bminer asynchronously using a goroutine.
func (m *Bminer) Start() {
	go m.run()
}

// Private methods

func (m *Bminer) reset() {
	// flush and cleanup
}

// Drive the Bminer state machine.
func (m *Bminer) run() {
	state := m.startState
	for state != nil {
		state = state(m)
	}
	m.reset()
}
