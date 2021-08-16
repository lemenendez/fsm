// fsm implements a simple Finite State Machine.
// A Finite State Machine has multiple use cases, the state of a User in a system,
// the state of a Tenant for a Saas company, the state of the execution of a long-running task,
// or the state of a money transfer.
package fsm

import (
	"errors"
	"fmt"
	"strings"
)

var ErrStateNotFound = errors.New("state not found")
var ErrStateAlExists = errors.New("state already exists")
var ErrTransNotAllowed = errors.New("transition not allowed")
var ErrTransAlExists = errors.New("transition already exists")
var ErrInvalidName = errors.New("invalid name")
var ErrExecNotAllowed = errors.New("execution not allowed")

const ready = "READY"
const notReady = "NOT_READY"
const check = "CHECK"

// A transition from State A to State B
// The combination of From, To, and Actions must be unique
type transition struct {
	From   string `json:"From"`
	To     string `json:"To"`
	Action string `json:"Action"`
}

func (t transition) String() string {
	return fmt.Sprintf("%v (%v) -> (%v)", t.From, t.To, t.Action)
}

// It represents a Finite State Machine.
type FSM struct {
	Name    string
	states  []string
	adj     []transition
	current string
	// state is the internal fsm state
	state *FSM
	// isInt stands for isInternal, flag to determine if this state is an internal state, used in conjunction with state field
	isInternal bool
}

// GetState gets the current state
func (p *FSM) GetState() string {
	return p.current
}

// validate actions
func (f *FSM) validate() {
	if len(f.states) > 0 &&
		len(f.adj) > 0 &&
		f.GetState() != ready {
		f.state.Exec(check, ready, nil)
		return
	}
	f.state.Exec(check, notReady, nil)
}

// createIntState creates a FSM for internal use
func createIntState() *FSM {
	f := &FSM{
		Name:       "FSM State",
		states:     make([]string, 0),
		adj:        make([]transition, 0),
		isInternal: true,
	}
	f.AddState(ready)
	f.AddState(notReady)
	f.Init(notReady)
	f.AddTrans(notReady, ready, check)
	f.AddTrans(ready, notReady, check)
	return f
}

// NewFSM creates a pointer to a brand new FSM
func NewFSM(name string) *FSM {
	f := &FSM{
		Name:   name,
		states: make([]string, 0),
		adj:    make([]transition, 0),
		state:  createIntState(),
	}
	return f
}

// AddState adds a new state into the fsm
// It validates state name prior to add it
// It validates state is unique
func (f *FSM) AddState(name string) error {
	if !isValidName(name) {
		return ErrInvalidName
	}
	if f.existsState(name) {
		return ErrStateAlExists
	}
	f.states = append(f.states, name)
	return nil
}

func (p FSM) existsState(name string) bool {
	for _, value := range p.states {
		if value == name {
			return true
		}
	}
	return false
}

// AdTrans adds a new transition between two states
// It validates transition name prior to add it
// It validates transition is unique
func (f *FSM) AddTrans(src string, des string, name string) error {
	if !isValidName(name) {
		return ErrInvalidName
	}
	if !f.existsState(src) || !f.existsState(des) {
		return ErrStateNotFound
	}

	for _, trans := range f.adj {
		if trans.From == src &&
			trans.To == des &&
			trans.Action == name {
			return ErrTransAlExists
		}
	}
	f.adj = append(f.adj, transition{
		From:   src,
		To:     des,
		Action: name,
	})
	return nil
}

// GetTrans returns the transitions string representation of fsm
func (f *FSM) GetTrans() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Transitions %v:\n", f.Name))
	for _, ad := range f.adj {
		builder.WriteString(fmt.Sprintf("%v\n", ad))
	}
	return builder.String()
}

// Init set the current state to the given state
// It validates stte exists prior to set it to current
func (p *FSM) Init(state string) error {
	if !p.existsState(state) {
		return ErrStateNotFound
	}
	p.current = state
	return nil
}

func (f *FSM) Exec(action string, des string, callback func(previous string, new string, action string)) error {
	if !f.isInternal {
		f.validate()
		if f.state.current != ready {
			return errors.New(string(notReady))
		}
	}
	if !f.existsState(des) {
		return ErrStateNotFound
	}

	for _, adj := range f.adj {
		if adj.From == f.current &&
			adj.To == des &&
			adj.Action == action {
			previous := f.current
			f.current = des
			if callback != nil {
				callback(previous, f.current, action)
				return nil
			}
		}
	}
	return ErrExecNotAllowed
}

func isValidName(name string) bool {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ_"
	const MAX = 64

	count := len(name)

	if count == 0 || count > MAX {
		return false
	}
	for _, chr := range name {
		if !strings.Contains(chars, string(chr)) {
			return false
		}
	}
	last := len(name) - 1
	if name[0] == '_' || name[last] == '_' {
		return false
	}
	return true
}
