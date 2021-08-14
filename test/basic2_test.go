package test

import (
	"testing"

	fsm "github.com/lemenendez/fsm"
)

func TestTransDups(t *testing.T) {

	const ACTIVE = "ACTIVE"
	const INACTIVE = "INACTIVE"
	const ACTIVATE = "ACTIVATE"

	f := fsm.NewFSM("BASIC")

	f.AddState(ACTIVE)
	f.AddState(INACTIVE)

	f.AddTrans(ACTIVE, INACTIVE, ACTIVATE)
	err := f.AddTrans(ACTIVE, INACTIVE, ACTIVATE)

	if err == nil {
		t.Errorf("should error dup transition")
	}

	err = f.AddTrans("NON_EXISTING1", "NON_EXISTING2", "ACTION")
	if err == nil {
		t.Errorf("SHOULD ERRORED")
	}
}

func TestBasic(t *testing.T) {

	const ACTIVE = "ACTIVE"
	const INACTIVE = "INACTIVE"

	fsm := fsm.NewFSM("BASIC ACTIVE/INACTIVE")

	fsm.AddState(ACTIVE)
	fsm.AddState(INACTIVE)

	err := fsm.Init(ACTIVE)

	if err != nil {
		t.Errorf("ACTIVE STATE SHOULD EXISTS")
	}

}

func TestBasicTransition(t *testing.T) {

	const ACTIVE = "ACTIVE"
	const INACTIVE = "INACTIVE"
	const ACTIVATE = "ACTIVATE"
	const DEACTIVATE = "DEACTIVATE"

	fsm := fsm.NewFSM("BASIC")

	fsm.AddState(ACTIVE)
	fsm.AddState(INACTIVE)

	fsm.AddTrans(ACTIVE, INACTIVE, ACTIVATE)
	fsm.AddTrans(INACTIVE, ACTIVE, DEACTIVATE)

	var err error
	err = fsm.Init(string(ACTIVE))
	if err != nil {
		t.Errorf("status should exists")
	}

	err = fsm.Init("DELETED")
	if err == nil {
		t.Errorf("should does not exist")
	}
}

func TestBasic3States(t *testing.T) {
	fsm := fsm.NewFSM("Three States")

	const TRIAL = "TRIAL"
	const BASIC = "BASIC"
	const PREMIUM = "PREMIUM"
	const UPGRATE = "UPGRATE"
	const DOWNGRATE = "DOWNGRATE"
	const NON_EXISTING = "NONEXISTING"

	fsm.AddState(TRIAL)
	fsm.AddState(BASIC)
	fsm.AddState(PREMIUM)

	fsm.AddTrans(TRIAL, BASIC, UPGRATE)
	fsm.AddTrans(TRIAL, PREMIUM, UPGRATE)
	fsm.AddTrans(BASIC, PREMIUM, UPGRATE)
	fsm.AddTrans(PREMIUM, BASIC, DOWNGRATE)

	trans := fsm.GetTrans()

	t.Log(trans)

	var err error

	err = fsm.AddTrans(NON_EXISTING, PREMIUM, "DUMMY_ACTION")

	if err == nil {
		t.Errorf("SHOULD not allow adding a transtion FROM a non loaded/non existing state")
	}

	err = fsm.AddTrans(BASIC, NON_EXISTING, "DUMMY_ACTION")
	//t.Log(err)
	if err == nil {
		t.Errorf("SHOULD not allow adding a transtion TO a non loaded/non existing state")
	}

	myFunc := func(pre string, cur string, action string) {
		t.Logf("Previous State:%v, New State:%v, Action:%v", pre, cur, action)
	}

	err = fsm.Init(TRIAL)
	if err != nil {
		t.Errorf("TRIAL STATE SHOULD EXISTS")
	}

	err = fsm.Exec(UPGRATE, BASIC, myFunc)

	if err != nil {
		t.Errorf(err.Error())
		t.Errorf("SHOULD ALLOW TRIAL->BASIC")
	}

	err = fsm.Exec(DOWNGRATE, TRIAL, myFunc)

	if err == nil {
		t.Errorf("SHOULD NOT ALLOW DOWNGRADE BASIC->TRIAL")
	}

	err = fsm.Exec(UPGRATE, PREMIUM, myFunc)

	if err != nil {
		t.Errorf("SHOULD ALLOW BASIC->PREMIUM")
	}

	err = fsm.Init("NOT_LOADED")
	if err == nil {
		t.Errorf("SHOLD NOT ALLOW NON EXISTING / NOT LOADED STATUS")
	}
}
