package fsm

import (
	"testing"
)

type Plan struct {
	Id    int32
	Name  string
	State *FSM
	t     *testing.T
}

func NewPlan(plan string) (*Plan, error) {

	f := NewFSM("SAAS Account State V1.0")

	err := f.AddState("TRIAL")
	if err != nil {
		return nil, err
	}
	err = f.AddState("BASIC")
	if err != nil {
		return nil, err
	}
	err = f.AddState("PREMIUM")
	if err != nil {
		return nil, err
	}

	err = f.AddTrans("TRIAL", "BASIC", "UPGRATE")
	if err != nil {
		return nil, err
	}
	err = f.AddTrans("TRIAL", "PREMIUM", "UPGRATE")
	if err != nil {
		return nil, err
	}
	err = f.AddTrans("BASIC", "PREMIUM", "UPGRATE")
	if err != nil {
		return nil, err
	}
	err = f.AddTrans("PREMIUM", "BASIC", "DOWNGRATE")
	if err != nil {
		return nil, err
	}

	c := &Plan{
		Id:    0,
		Name:  "Standard",
		State: f,
	}
	err = f.Init(plan)
	if err == nil {
		return c, nil
	}
	return nil, err
}

func NewPlan2(plan string) (*Plan, error) {

	f := NewFSM("SAAS Account State V2.0")

	err := f.AddState("TRIAL")
	if err != nil {
		return nil, err
	}
	err = f.AddState("BASIC")
	if err != nil {
		return nil, err
	}
	err = f.AddState("PREMIUM")
	if err != nil {
		return nil, err
	}
	err = f.AddState("GOLD")
	if err != nil {
		return nil, err
	}

	err = f.AddTrans("TRIAL", "BASIC", "UPGRATE")
	if err != nil {
		return nil, err
	}
	err = f.AddTrans("TRIAL", "PREMIUM", "UPGRATE")
	if err != nil {
		return nil, err
	}
	err = f.AddTrans("TRIAL", "GOLD", "UPGRATE")
	if err != nil {
		return nil, err
	}

	err = f.AddTrans("BASIC", "PREMIUM", "UPGRATE")
	if err != nil {
		return nil, err
	}
	err = f.AddTrans("PREMIUM", "BASIC", "DOWNGRATE")
	if err != nil {
		return nil, err
	}
	err = f.AddTrans("GOLD", "PREMIUM", "DOWNGRATE")
	if err != nil {
		return nil, err
	}
	err = f.AddTrans("GOLD", "BASIC", "DOWNGRATE")
	if err != nil {
		return nil, err
	}

	c := &Plan{
		Id:    0,
		Name:  "Standard",
		State: f,
	}
	err = f.Init(plan)
	if err == nil {
		return c, nil
	}
	return nil, err
}

func (a *Plan) stateTransitionHandler(pre string, cur string, action string) {
	a.t.Logf("Previous State:%v, New State:%v, Action:%v", pre, cur, action)
}

func (a *Plan) Upgrade(name string) bool {
	err := a.State.Exec("UPGRATE", name, a.stateTransitionHandler)
	if err != nil {
		a.t.Logf(err.Error())
	}
	return err == nil
}

func (a *Plan) Downgrate(name string) bool {
	err := a.State.Exec("DOWNGRATE", name, a.stateTransitionHandler)
	if err != nil {
		a.t.Logf(err.Error())
	}
	return err == nil
}

func TestObject1(t *testing.T) {
	a, err := NewPlan("TRIAL")
	if a != nil {
		a.t = t
		if a.Upgrade("BASIC") {
			t.Logf("New State:%v", a.State.GetState())
		} else {
			t.Errorf("cannot updagrade")
		}
	} else {
		t.Errorf(err.Error())
	}
}

func TestObject2(t *testing.T) {
	a, err := NewPlan("BASIC")
	if a != nil {
		a.t = t
		if a.Upgrade("PREMIUM") {
			t.Logf("New State:%v", a.State.GetState())
		} else {
			t.Errorf("cannot updagrade")
		}
	} else {
		t.Errorf(err.Error())
	}
}

func TestObject3(t *testing.T) {
	a, err := NewPlan("PREMIUM")
	if a != nil {
		a.t = t
		if a.Downgrate("BASIC") {
			t.Logf("New State:%v", a.State.GetState())
		} else {
			t.Errorf("cannot updagrade")
		}
	} else {
		t.Errorf(err.Error())
	}
}
