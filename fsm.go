// fsm implements a simple Finite State Machine.
// A Finite State Machine has multiple use cases, the state of a User in a system,
// the state of a Tenant for a Saas company, the state of the execution of a long-running task,
// or the state of a money transfer.
package fsm

import (
	"encoding/json"
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

// A possible transition from an state A (src) to state B (des).
// For a given FSM the combination of: src,des and name must be unique.
// The code must enforce that rule all the time.
type transition struct {
	// source state
	src string
	// destination state)
	des string
	// transition's name
	name string
}

func (t transition) String() string {
	return fmt.Sprintf("%v (%v) -> (%v)", t.name, t.src, t.des)
}

// It represents a Finite State Machine.
type FSM struct {
	Id   int32
	Name string
	// set of states
	states []string
	// set of transitions
	adj     []transition
	current string
	// an FSM itself has its own internal state, used in conjunction with isInt field
	state *FSM
	// isInt stands for isInternal, flag to determine if this state is an internal state, used in conjunction with state field
	isInt bool
}

// It gets the name of the current state
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

// It creates a FSM for internal use
func createIntState() *FSM {
	f := &FSM{
		Id:     0,
		Name:   "FSM State",
		states: make([]string, 0),
		adj:    make([]transition, 0),
		isInt:  true,
	}
	f.AddState(ready)
	f.AddState(notReady)
	f.Init(notReady)
	f.AddTrans(notReady, ready, check)
	f.AddTrans(ready, notReady, check)
	return f
}

func NewFSM(name string) *FSM {
	f := &FSM{
		Id:     0,
		Name:   name,
		states: make([]string, 0),
		adj:    make([]transition, 0),
		state:  createIntState(),
	}
	return f
}

func (f *FSM) AddState(name string) error {
	if !valName(name) {
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

func (p *FSM) AddTrans(src string, des string, name string) error {
	if !valName(name) {
		return ErrInvalidName
	}
	if !p.existsState(src) || !p.existsState(des) {
		return ErrStateNotFound
	}

	for _, trans := range p.adj {
		if trans.src == src &&
			trans.des == des &&
			trans.name == name {
			return ErrTransAlExists
		}
	}
	t := transition{
		src:  src,
		des:  des,
		name: name,
	}
	p.adj = append(p.adj, t)
	return nil
}

func (f *FSM) GetTrans() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Transitions %v:\n", f.Name))
	for _, ad := range f.adj {
		builder.WriteString(fmt.Sprintf("%v\n", ad))
	}
	return builder.String()
}

func (p *FSM) Init(name string) error {
	if !p.existsState(name) {
		return ErrStateNotFound
	}
	p.current = name
	return nil
}

func (f *FSM) Exec(action string, des string, callback func(previous string, new string, action string)) error {
	if !f.isInt {
		f.validate()
		if f.state.current != ready {
			return errors.New(string(notReady))
		}
	}
	if !f.existsState(des) {
		return ErrStateNotFound
	}

	for i := 0; i < len(f.adj); i++ {
		if f.adj[i].src == f.current &&
			f.adj[i].des == des &&
			f.adj[i].name == action {
			// update the current state
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

func (f *FSM) MarshalJSON() ([]byte, error) {
	var states []string

	states = append(states, f.states...)

	type tran struct {
		From, To, Action string
	}

	trans := make([]tran, 0)
	for i := 0; i < len(f.adj); i++ {
		newtran := tran{
			From:   f.adj[i].src,
			To:     f.adj[i].des,
			Action: f.adj[i].name,
		}
		trans = append(trans, newtran)
	}
	return json.Marshal(&struct {
		Name        string   `json:"Name"`
		Current     string   `json:"Current"`
		States      []string `json:"States"`
		Transitions []tran   `json:"Transitions"`
	}{
		Name:        f.Name,
		Current:     f.GetState(),
		States:      states,
		Transitions: trans,
	})
}

func (f *FSM) UnmarshalJSON(data []byte) error {
	type tran struct {
		From, To, Action string
	}
	temp := struct {
		Name        string   `json:"Name"`
		Current     string   `json:"Current"`
		States      []string `json:"States"`
		Transitions []tran   `json:"Transitions"`
	}{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	f = &FSM{
		Name:  temp.Name,
		adj:   make([]transition, 0),
		state: createIntState(),
	}
	f.states = make([]string, 0)
	for _, val := range temp.States {
		if !f.existsState(val) {
			f.states = append(f.states, val)
		}
	}
	for _, trans := range temp.Transitions {
		if !valName(trans.Action) {
			return ErrInvalidName
		}
		if !f.existsState(trans.From) || !f.existsState(trans.To) {
			return ErrStateNotFound
		}

		f.adj = append(f.adj, transition{
			src:  trans.From,
			des:  trans.To,
			name: trans.Action,
		})
	}
	err := f.Init(temp.Current)
	if err != nil {
		return err
	}
	return nil
}

func valName(name string) bool {
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
