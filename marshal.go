package fsm

import (
	"encoding/json"
)

func (t transition) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		From   string `json:"from"`
		To     string `json:"to"`
		Action string `json:"action"`
	}{
		From:   t.From,
		To:     t.To,
		Action: t.Action,
	})
}

func (t *transition) UnmarshalJSON(data []byte) error {
	temp := struct {
		From   string `json:"from"`
		To     string `json:"to"`
		Action string `json:"action"`
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

	for key := range f.states {
		states = append(states, key)
	}

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
		Name        string       `json:"name"`
		Current     string       `json:"current"`
		States      []string     `json:"states"`
		Transitions []transition `json:"transitions"`
	}{
		Name:        f.Name,
		Current:     f.GetState(),
		States:      states,
		Transitions: trans,
	})
}

func (f *FSM) UnmarshalJSON(data []byte) error {
	temp := struct {
		Name        string       `json:"name"`
		Current     string       `json:"current"`
		States      []string     `json:"states"`
		Transitions []transition `json:"transitions"`
	}{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	f.Name = temp.Name
	f.state = createIntState()
	f.states = make(map[string]bool)
	f.adj = make([]transition, 0)

	for _, val := range temp.States {
		if ok := f.states[val]; !ok {
			f.states[val] = true
		}
	}
	for _, trans := range temp.Transitions {
		if !isValidName(trans.Action) {
			return ErrInvalidName
		}
		if ok := f.states[trans.From]; !ok {
			return ErrStateNotFound
		}
		if ok := f.states[trans.To]; !ok {
			return ErrStateNotFound
		}

		f.adj = append(f.adj, transition{
			From:   trans.From,
			To:     trans.To,
			Action: trans.Action,
		})
	}
	if len(f.adj) > 0 {
		f.state.current = ready
	}
	err := f.Init(temp.Current)
	if err != nil {
		return err
	}
	return nil
}
