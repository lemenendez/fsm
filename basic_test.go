package fsm

import (
	"testing"
)

func TestCheckDupState(t *testing.T) {

	const ACTIVE = "ACTIVE"
	const INACTIVE = "INACTIVE"

	f := NewFSM("BASIC")

	err := f.AddState(ACTIVE)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = f.AddState(INACTIVE)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = f.AddState(ACTIVE)
	if err == nil {
		t.Errorf("Should errored state already exists")
	}

}

func TestAddTransitions(t *testing.T) {

	const ACTIVE = "ACTIVE"
	const INACTIVE = "INACTIVE"
	const ACTIVATE = "ACTIVATE"
	const DEACTIVATE = "DEACTIVATE"

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
		t.Errorf(err.Error())
	}
	err = f.AddTrans(INACTIVE, ACTIVE, DEACTIVATE)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestWrongNames(t *testing.T) {
	const INVALIDNAME = " in valid 1 "
	const INVALIDNAME2 = "_INVALID"
	const INVALIDNAME3 = "INVALID_"
	const INVALIDNAME4 = "invalid"
	const INVALIDNAME5 = "THISISAVERYLONGLONGLONGNAMEANDSHOULDBE_AVOIDED_AT_ALL_COST_ACBCDEFTGADERE"
	const ACTIVATE = "ACTIVATE"
	const INACTIVE = "INACTIVE"
	const ACTIVE = "ACTIVE"
	const DEACTIVATE = "DEACTIVATE"

	f := NewFSM("BASIC")

	f.AddState(ACTIVE)
	f.AddState(INACTIVE)
	err := f.AddState(INVALIDNAME)
	if err == nil {
		t.Errorf("Should errored name is invalid")
	}

	err = f.AddState(INVALIDNAME2)
	if err == nil {
		t.Errorf("Should errored name is invalid")
	}

	err = f.AddState(INVALIDNAME3)
	if err == nil {
		t.Errorf("Should errored name is invalid")
	}

	err = f.AddState(INVALIDNAME4)
	if err == nil {
		t.Errorf("Should errored name is invalid")
	}

	err = f.AddState(INVALIDNAME5)
	if err == nil {
		t.Errorf("Should errored name is invalid")
	}

	f.AddTrans(ACTIVE, INACTIVE, ACTIVATE)
	f.AddTrans(INACTIVE, ACTIVE, DEACTIVATE)
	err = f.AddTrans(INACTIVE, ACTIVE, INVALIDNAME)
	if err == nil {
		t.Errorf("Should errored name is invalid")
	}

	err = f.AddTrans(INACTIVE, ACTIVE, INVALIDNAME2)
	if err == nil {
		t.Errorf("Should errored name is invalid")
	}

	err = f.AddTrans(INACTIVE, ACTIVE, INVALIDNAME3)
	if err == nil {
		t.Errorf("Should errored name is invalid")
	}
}

func TestGetName(t *testing.T) {

	const ACTIVE = "ACTIVE"
	const INACTIVE = "INACTIVE"
	const ACTIVATE = "ACTIVATE"
	const DEACTIVATE = "DEACTIVATE"

	f := NewFSM("BASIC")

	f.AddState(ACTIVE)
	f.AddState(INACTIVE)

	f.AddTrans(ACTIVE, INACTIVE, ACTIVATE)
	f.AddTrans(INACTIVE, ACTIVE, DEACTIVATE)

	f.Init(INACTIVE)

	t.Log(f.GetState())
}

func TestExecWithEmptyStates(t *testing.T) {
	f := NewFSM("BASIC")

	myFunc := func(pre string, cur string, action string) {
		t.Logf("Previous State:%v, New State:%v, Action:%v", pre, cur, action)
	}

	err := f.Exec("DO", "SOME_MORE", myFunc)

	if err == nil {
		t.Errorf("It should failed:fsm does not have any states")
	}

}

func TestExecWithEmptyTrans(t *testing.T) {
	f := NewFSM("BASIC")

	f.AddState("ACTIVE")
	f.AddState("INACTIVE")
	f.AddTrans("ACTIVE", "INACTIVE", "ACTIVATE")
	f.Init("ACTIVE")

	myFunc := func(pre string, cur string, action string) {
		t.Logf("Previous State:%v, New State:%v, Action:%v", pre, cur, action)
	}

	err := f.Exec("DO", "SOME_MORE", myFunc)

	if err == nil {
		t.Errorf("It should failed: SOME_MORE does not exist")
	}

}
