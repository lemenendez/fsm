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

// It represents a unique state
type State struct {
	Name string
}

// A possible transition from an estate A (src) to state B (des).
// For a given FSM the combination of: src,des and name must be unique.
// The code must enforce that rule all the time.
type transition struct {
	// pointer to the source state (A)
	src *State
	// pointer to destination state (B)
	des *State
	// transition's name
	name string
}

func (t transition) String() string {
	return fmt.Sprintf("%v (%v) -> (%v)", t.name, t.src.Name, t.des.Name)
}

// It represents a Finite State Machine.
type FSM struct {
	Id   int32
	Name string
	// set of states
	states map[string]*State
	// set of transitions
	adj []transition
	// pointer to the current state
	current *State
	// an FSM itself has its own internal state, used in conjunction with isInt field
	state *FSM
	// isInt stands for isInternal, flag to determine if this state is an internal state, used in conjunction with state field
	isInt bool
}

// It gets the name of the current state
func (p *FSM) GetState() string {
	return p.current.Name
}

// validate actions
func (f *FSM) validate() {
	if len(f.states) > 0 &&
		len(f.adj) > 0 &&
		f.current != nil &&
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
		states: make(map[string]*State),
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
		states: make(map[string]*State),
		adj:    make([]transition, 0),
		state:  createIntState(),
	}
	return f
}

func (p *FSM) AddState(name string) error {
	if !valName(name) {
		return ErrInvalidName
	}
	state := p.findStateByName(name)
	if state != nil {
		return ErrStateAlExists
	}
	nState := State{
		Name: name,
	}
	p.states[name] = &nState
	return nil
}

func (p FSM) findStateByName(name string) *State {
	for _, value := range p.states {
		if value.Name == name {
			return value
		}
	}
	return nil
}

func (p *FSM) AddTrans(src string, des string, name string) error {
	if !valName(name) {
		return ErrInvalidName
	}
	srcState := p.findStateByName(src)
	if srcState == nil {
		return ErrStateNotFound
	}
	desState := p.findStateByName(des)
	if desState == nil {
		return ErrStateNotFound
	}

	for i := 0; i < len(p.adj); i++ {
		if p.adj[i].src == srcState &&
			p.adj[i].des == desState &&
			p.adj[i].name == name {
			return ErrTransAlExists
		}
	}
	t := transition{
		src:  srcState,
		des:  desState,
		name: name,
	}
	p.adj = append(p.adj, t)
	return nil
}

func (p *FSM) GetTrans() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Transitions %v:\n", p.Name))
	for i := 0; i < len(p.adj); i++ {
		builder.WriteString(fmt.Sprintf("%v\n", p.adj[i]))
	}
	result := builder.String()
	return result
}

func (p *FSM) Init(name string) error {
	state := p.findStateByName(name)
	if state == nil {
		return ErrStateNotFound
	}
	p.current = state
	return nil
}

func (f *FSM) Exec(action string, des string, callback func(previous string, new string, action string)) error {
	if !f.isInt {
		f.validate()
		if f.state.current.Name != ready {
			return errors.New(string(notReady))
		}
	}
	desState := f.findStateByName(des)

	if desState == nil {
		return ErrStateNotFound
	}

	for i := 0; i < len(f.adj); i++ {
		if f.adj[i].src == f.current &&
			f.adj[i].des == desState &&
			f.adj[i].name == action {
			// update the current state
			previous := f.current.Name
			f.current = desState
			if callback != nil {
				callback(previous, f.current.Name, action)
				return nil
			}
		}
	}

	return ErrExecNotAllowed
}

func (f *FSM) MarshalJSON() ([]byte, error) {
	var states []*State
	for _, value := range f.states {
		states = append(states, value)
	}

	type tran struct {
		From, To, Action string
	}

	trans := make([]tran, 0)
	for i := 0; i < len(f.adj); i++ {
		newtran := tran{
			From:   (*f.adj[i].src).Name,
			To:     (*f.adj[i].des).Name,
			Action: f.adj[i].name,
		}
		trans = append(trans, newtran)
	}
	return json.Marshal(&struct {
		Name        string   `json:"Name"`
		Current     string   `json:"Current"`
		States      []*State `json:"States"`
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
		States      []*State `json:"States"`
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

	f.states = make(map[string]*State)

	for i := 0; i < len(temp.States); i++ {
		f.states[temp.States[i].Name] = temp.States[i]
	}

	for i := 0; i < len(temp.Transitions); i++ {

		if !valName(temp.Transitions[i].Action) {
			return ErrInvalidName
		}

		srcState := f.findStateByName(temp.Transitions[i].From)
		if srcState == nil {
			return ErrStateNotFound
		}
		desState := f.findStateByName(temp.Transitions[i].To)
		if desState == nil {
			return ErrStateNotFound
		}
		t := transition{
			src:  srcState,
			des:  desState,
			name: temp.Transitions[i].Action,
		}
		f.adj = append(f.adj, t)

	}

	f.Init(temp.Current)

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
