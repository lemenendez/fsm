package fsm

import (
	"encoding/json"
)

func (t transition) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		From   string `json:"From"`
		To     string `json:"To"`
		Action string `json:"Action"`
	}{
		From:   t.From,
		To:     t.To,
		Action: t.Action,
	})
}

func (t *transition) UnmarshalJSON(data []byte) error {
	temp := struct {
		From   string `json:"From"`
		To     string `json:"To"`
		Action string `json:"Action"`
	}{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	t.From = temp.From
	t.To = temp.To
	t.Action = temp.Action

	return nil
}

func (f *FSM) MarshalJSON() ([]byte, error) {
	var states []string

	states = append(states, f.states...)

	trans := make([]transition, 0)
	for _, adj := range f.adj {
		newtran := transition{
			From:   adj.From,
			To:     adj.To,
			Action: adj.Action,
		}
		trans = append(trans, newtran)
	}
	return json.Marshal(&struct {
		Name        string       `json:"Name"`
		Current     string       `json:"Current"`
		States      []string     `json:"States"`
		Transitions []transition `json:"Transitions"`
	}{
		Name:        f.Name,
		Current:     f.GetState(),
		States:      states,
		Transitions: trans,
	})
}

func (f *FSM) UnmarshalJSON(data []byte) error {
	temp := struct {
		Name        string       `json:"Name"`
		Current     string       `json:"Current"`
		States      []string     `json:"States"`
		Transitions []transition `json:"Transitions"`
	}{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	f.Name = temp.Name
	f.state = createIntState()
	f.states = make([]string, 0)
	f.adj = make([]transition, 0)

	for _, val := range temp.States {
		if !f.existsState(val) {
			f.states = append(f.states, val)
		}
	}
	for _, trans := range temp.Transitions {
		if !isValidName(trans.Action) {
			return ErrInvalidName
		}
		if !f.existsState(trans.From) || !f.existsState(trans.To) {
			return ErrStateNotFound
		}

		f.adj = append(f.adj, transition{
			From:   trans.From,
			To:     trans.To,
			Action: trans.Action,
		})
	}
	err := f.Init(temp.Current)
	if err != nil {
		return err
	}
	return nil
}
