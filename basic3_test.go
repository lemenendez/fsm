package fsm

import (
	"testing"
)

func TestSimpleEnableDisableActions(t *testing.T) {
	var err error

	const ENABLED = "ENABLED"
	const DISABLED = "DISABLED"
	const ENABLE = "ENABLE"
	const DISABLE = "DISABLE"

	f := NewFSM("Basic Disabled/Enabled")
	err = f.AddState(DISABLED)
	if err != nil {
		t.Fail()
	}
	err = f.AddState(ENABLED)
	if err != nil {
		t.Fail()
	}
	err = f.Init(DISABLED)
	if err != nil {
		t.Fail()
	}
	err = f.AddTrans(ENABLED, DISABLED, DISABLE)
	if err != nil {
		t.Fail()
	}
	err = f.AddTrans(DISABLED, ENABLED, ENABLE)
	if err != nil {
		t.Fail()
	}
	myFunc := func(pre string, cur string, action string) {
		t.Logf("Previous State:%v, New State:%v, Action:%v", pre, cur, action)
	}

	err = f.Exec(ENABLE, ENABLED, myFunc)

	if err != nil {
		t.Log(err)
		t.Errorf("should allow transicion")
	}
	t.Log(f.GetTrans())

}
