// fsm implements a simple Finite State Machine.
// A Finite State Machine has multiple use cases, the state of a User in a system,
// the state of a Tenant for a Saas company, the state of the execution of a long-running task,
// or the state of a money transfer.
package fsm

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var ErrStateNotFound error = errors.New("state not found")
var ErrStateAlExists = errors.New("state already exists")
var ErrTransNotAllowed = errors.New("transition not allowed")
var ErrTransAlExists = errors.New("transition already exists")
var ErrInvalidName = errors.New("invalid name")
var ErrExecNotAllowed = errors.New("execution not allowed")
var ErrNotReady = errors.New("not ready")

var exp = regexp.MustCompile(`^[A-Z]+(_?[A-Z])*$`)

const ready = "READY"
const notReady = "NOT_READY"
const check = "CHECK"

// A transition from State A to State B
// The combination of From, To, and Actions must be unique
type transition struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Action string `json:"action"`
}

func (t transition) String() string {
	return fmt.Sprintf("%v (%v) -> (%v)", t.From, t.To, t.Action)
}

// It represents a Finite State Machine.
type FSM struct {
	Name    string
	states  map[string]bool
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

// createIntState creates a FSM for internal use
func createIntState() *FSM {
	var err error
	f := &FSM{
		Name:       "FSM State",
		states:     make(map[string]bool),
		adj:        make([]transition, 0),
		isInternal: true,
	}
	if err = f.AddState(ready); err != nil {
		return nil
	}
	if err = f.AddState(notReady); err != nil {
		return nil
	}
	if err = f.Init(notReady); err != nil {
		return nil
	}
	if err = f.AddTrans(notReady, ready, check); err != nil {
		return nil
	}
	if err = f.AddTrans(ready, notReady, check); err != nil {
		return nil
	}
	return f
}

// NewFSM creates a pointer to a brand new FSM
func NewFSM(name string) *FSM {
	f := &FSM{
		Name:   name,
		states: make(map[string]bool),
		adj:    make([]transition, 0),
		state:  createIntState(),
	}
	// set int state to not ready
	f.state.current = notReady
	return f
}

// AddState adds a new state into the fsm
// It validates state name prior to add it
// It validates state is unique
func (f *FSM) AddState(name string) error {
	if !isValidName(name) {
		return ErrInvalidName
	}
	if ok := f.states[name]; ok {
		return ErrStateAlExists
	}

	f.states[name] = true
	return nil
}

// AdTrans adds a new transition between two states
// It validates transition name prior to add it
// It validates transition is unique
func (f *FSM) AddTrans(src string, des string, name string) error {
	if !isValidName(name) {
		return ErrInvalidName
	}
	if ok := f.states[src]; !ok {
		return ErrStateNotFound
	}
	if ok := f.states[des]; !ok {
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
	if !f.isInternal {
		f.state.current = ready
	}
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
// It validates state exists prior to set it to current
func (f *FSM) Init(state string) error {
	if ok := f.states[state]; !ok {
		return ErrStateNotFound
	}
	f.current = state
	return nil
}

func (f *FSM) Exec(action string, des string, callback func(previous string, new string, action string)) error {
	if !f.isInternal {
		if f.state.current != ready {
			return ErrNotReady
		}
	}
	if ok := f.states[des]; !ok {
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
	const MAX = 64
	if len(name) > MAX {
		return false
	}
	return exp.MatchString(name)
}

func New(name string, trans [][3]string) (*FSM, error) {
	f := &FSM{
		Name:   name,
		states: make(map[string]bool),
		adj:    make([]transition, 0),
		state:  createIntState(),
	}

	for _, t := range trans {
		if err := f.AddState(t[0]); err != nil && err != ErrStateAlExists {
			return nil, err
		}
		if err := f.AddState(t[1]); err != nil && err != ErrStateAlExists {
			return nil, err
		}
		err := f.AddTrans(t[0], t[1], t[2])
		if err != nil {
			return nil, err
		}
	}
	return f, nil
}
