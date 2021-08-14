package test

import (
	"testing"

	fsm "github.com/lemenendez/fsm"
)

type Plan struct {
	Id    int32
	Name  string
	State *fsm.FSM
	t     *testing.T
}

func NewPlan(plan string) (*Plan, error) {

	f := fsm.NewFSM("SAAS Account State V1.0")

	f.AddState("TRIAL")
	f.AddState("BASIC")
	f.AddState("PREMIUM")

	f.AddTrans("TRIAL", "BASIC", "UPGRATE")
	f.AddTrans("TRIAL", "PREMIUM", "UPGRATE")
	f.AddTrans("BASIC", "PREMIUM", "UPGRATE")
	f.AddTrans("PREMIUM", "BASIC", "DOWNGRATE")

	c := &Plan{
		Id:    0,
		Name:  "Standard",
		State: f,
	}
	err := f.Init(plan)
	if err == nil {
		return c, nil
	}
	return nil, err
}

func NewPlan2(plan string) (*Plan, error) {

	f := fsm.NewFSM("SAAS Account State V2.0")

	f.AddState("TRIAL")
	f.AddState("BASIC")
	f.AddState("PREMIUM")
	f.AddState("GOLD")

	f.AddTrans("TRIAL", "BASIC", "UPGRATE")
	f.AddTrans("TRIAL", "PREMIUM", "UPGRATE")
	f.AddTrans("TRIAL", "GOLD", "UPGRATE")

	f.AddTrans("BASIC", "PREMIUM", "UPGRATE")
	f.AddTrans("PREMIUM", "BASIC", "DOWNGRATE")
	f.AddTrans("GOLD", "PREMIUM", "DOWNGRATE")
	f.AddTrans("GOLD", "BASIC", "DOWNGRATE")

	c := &Plan{
		Id:    0,
		Name:  "Standard",
		State: f,
	}
	err := f.Init(plan)
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

/*
func TestMigration1(t *testing.T) {
	plan1, err := NewPlan("TRIAL")
	if err != nil {
		t.Errorf(err.Error())
	}
	plan1.t = t
	plan2, err := NewPlan2("PREMIUM")
	if err != nil {
		t.Errorf(err.Error())
	}

	plan2.t = t

	if plan1.Upgrade("PREMIUM") {
		t.Logf("New State:%v", plan1.State.GetState())
	} else {
		t.Errorf("cannot updagrade")
	}
}
*/
