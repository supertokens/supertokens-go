package core

import (
	"flag"
	"sync"
)

type processState struct {
	history []int
}

var processStateInstantiated *processState
var processStateLock sync.Mutex

// ResetProcessState to be used for testing only
func ResetProcessState() {
	processStateInstantiated = nil
}

// GetProcessStateInstance used to get processState struct
func GetProcessStateInstance() *processState {
	if processStateInstantiated == nil {
		processStateLock.Lock()
		defer processStateLock.Unlock()
		if processStateInstantiated == nil {
			processStateInstantiated = &processState{
				history: []int{},
			}
		}
	}
	return processStateInstantiated
}

func (p *processState) GetLastEventByName(state int) *int {
	if flag.Lookup("test.v") == nil {
		return nil
	}
	processStateLock.Lock()
	defer processStateLock.Unlock()

	for i := len(p.history) - 1; i >= 0; i-- {
		if p.history[i] == state {
			return &p.history[i]
		}
	}
	return nil
}

func (p *processState) AddState(state int) {
	if flag.Lookup("test.v") == nil {
		return
	}
	processStateLock.Lock()
	defer processStateLock.Unlock()
	p.history = append(p.history, state)
}

// Process states
const (
	CallingServiceInVerify = iota
)
