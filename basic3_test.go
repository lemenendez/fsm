package fsm

import (
	"testing"
)

func TestSimpleEnableDisableActions(t *testing.T) {
	const ENABLED = "ENABLED"
	const DISABLED = "DISABLED"
	const ENABLE = "ENABLE"
	const DISABLE = "DISABLE"

	f := NewFSM("Basic Disabled/Enabled")
	f.AddState(DISABLED)
	f.AddState(ENABLED)
	f.Init(DISABLED)

	f.AddTrans(ENABLED, DISABLED, DISABLE)
	f.AddTrans(DISABLED, ENABLED, ENABLE)
	myFunc := func(pre string, cur string, action string) {
		t.Logf("Previous State:%v, New State:%v, Action:%v", pre, cur, action)
	}

	err := f.Exec(ENABLE, ENABLED, myFunc)

	if err != nil {
		t.Log(err)
		t.Errorf("should allow transicion")
	}
	t.Log(f.GetTrans())

}
