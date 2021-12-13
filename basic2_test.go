package fsm

import (
	"testing"
)

func TestTransDups(t *testing.T) {

	const ACTIVE = "ACTIVE"
	const INACTIVE = "INACTIVE"
	const ACTIVATE = "ACTIVATE"

	f := NewFSM("BASIC")

	err := f.AddState(ACTIVE)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = f.AddState(INACTIVE)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = f.AddTrans(ACTIVE, INACTIVE, ACTIVATE)
	if err != nil {
		t.Fail()
	}
	err = f.AddTrans(ACTIVE, INACTIVE, ACTIVATE)

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

	fsm := NewFSM("BASIC ACTIVE/INACTIVE")

	err := fsm.AddState(ACTIVE)
	if err != nil {
		t.Fail()
	}
	err = fsm.AddState(INACTIVE)
	if err != nil {
		t.Fail()
	}

	err = fsm.Init(ACTIVE)

	if err != nil {
		t.Errorf("ACTIVE STATE SHOULD EXISTS")
	}

}

func TestBasicTransition(t *testing.T) {
	var err error
	const ACTIVE = "ACTIVE"
	const INACTIVE = "INACTIVE"
	const ACTIVATE = "ACTIVATE"
	const DEACTIVATE = "DEACTIVATE"

	fsm := NewFSM("BASIC")

	err = fsm.AddState(ACTIVE)
	if err != nil {
		t.Fail()
	}
	err = fsm.AddState(INACTIVE)
	if err != nil {
		t.Fail()
	}
	err = fsm.AddTrans(ACTIVE, INACTIVE, ACTIVATE)
	if err != nil {
		t.Fail()
	}
	err = fsm.AddTrans(INACTIVE, ACTIVE, DEACTIVATE)
	if err != nil {
		t.Fail()
	}
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
	var err error
	f := NewFSM("Three States")

	const TRIAL = "TRIAL"
	const BASIC = "BASIC"
	const PREMIUM = "PREMIUM"
	const UPGRATE = "UPGRATE"
	const DOWNGRATE = "DOWNGRATE"
	const NON_EXISTING = "NONEXISTING"

	err = f.AddState(TRIAL)
	if err != nil {
		t.Fail()
	}
	err = f.AddState(BASIC)
	if err != nil {
		t.Fail()
	}
	err = f.AddState(PREMIUM)
	if err != nil {
		t.Fail()
	}

	err = f.AddTrans(TRIAL, BASIC, UPGRATE)
	if err != nil {
		t.Fail()
	}
	err = f.AddTrans(TRIAL, PREMIUM, UPGRATE)
	if err != nil {
		t.Fail()
	}
	err = f.AddTrans(BASIC, PREMIUM, UPGRATE)
	if err != nil {
		t.Fail()
	}
	err = f.AddTrans(PREMIUM, BASIC, DOWNGRATE)
	if err != nil {
		t.Fail()
	}

	trans := f.GetTrans()

	t.Log(trans)

	err = f.AddTrans(NON_EXISTING, PREMIUM, "DUMMY_ACTION")

	if err == nil {
		t.Errorf("SHOULD not allow adding a transtion FROM a non loaded/non existing state")
	}

	err = f.AddTrans(BASIC, NON_EXISTING, "DUMMY_ACTION")
	if err == nil {
		t.Errorf("SHOULD not allow adding a transtion TO a non loaded/non existing state")
	}

	myFunc := func(pre string, cur string, action string) {
		t.Logf("Previous State:%v, New State:%v, Action:%v", pre, cur, action)
	}

	err = f.Init(TRIAL)
	if err != nil {
		t.Errorf("TRIAL STATE SHOULD EXISTS")
	}

	err = f.Exec(UPGRATE, BASIC, myFunc)

	if err != nil {
		t.Errorf(err.Error())
		t.Errorf("SHOULD ALLOW TRIAL->BASIC")
	}

	err = f.Exec(DOWNGRATE, TRIAL, myFunc)

	if err == nil {
		t.Errorf("SHOULD NOT ALLOW DOWNGRADE BASIC->TRIAL")
	}

	err = f.Exec(UPGRATE, PREMIUM, myFunc)

	if err != nil {
		t.Errorf("SHOULD ALLOW BASIC->PREMIUM")
	}

	err = f.Init("NOT_LOADED")
	if err == nil {
		t.Errorf("SHOLD NOT ALLOW NON EXISTING / NOT LOADED STATUS")
	}
}
