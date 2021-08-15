package fsm

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestBasic3StatesJson(t *testing.T) {

	const TRIAL = "TRIAL"
	const BASIC = "BASIC"
	const PREMIUM = "PREMIUM"
	const UPGRATE = "UPGRATE"
	const DOWNGRATE = "DOWNGRATE"

	f := NewFSM("Customer Plan")

	err := f.AddState(TRIAL)
	if err != nil {
		t.Fatal(err)
	}
	err = f.AddState(BASIC)
	if err != nil {
		t.Fatal(err)
	}
	err = f.AddState(PREMIUM)
	if err != nil {
		t.Fatal(err)
	}

	err = f.AddTrans(TRIAL, BASIC, UPGRATE)
	if err != nil {
		t.Fatal(err)
	}
	err = f.AddTrans(TRIAL, PREMIUM, UPGRATE)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = f.AddTrans(BASIC, PREMIUM, UPGRATE)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = f.AddTrans(PREMIUM, BASIC, DOWNGRATE)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = f.Init(TRIAL)

	if err != nil {
		t.Errorf(err.Error())
	}

	b, err := json.Marshal(f)
	if err != nil {
		fmt.Println(err)
		return
	}

	var f2 FSM
	err = json.Unmarshal(b, &f2)

	myFunc := func(pre string, cur string, action string) {
		t.Logf("Previous State:%v, New State:%v, Action:%v", pre, cur, action)
	}

	if err != nil {
		err = f2.Exec(UPGRATE, BASIC, myFunc)
		if err != nil {
			t.Logf(err.Error())
			t.Errorf("It should work")
		}
	}
}

func TestBasic3StatesJsonUnmarshal1(t *testing.T) {
	b := []byte(`{"Id":"0","UUID":"642f24f0-132d-4413-80b9-0fc43408c302","Name":"Customer Plan Status","Current":"TRIAL","States":[{"Id":0,"UUID":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","Name":"TRIAL"},{"Id":1,"UUID":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Name":"BASIC"},{"Id":2,"UUID":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Name":"PREMIUM"}],"Transitions":[{"From":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","To":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Action":"UPGRATE"},{"From":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","To":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Action":"UPGRATE"},{"From":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","To":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Action":"UPGRATE"},{"From":"c7560801-7ade-4a73-866b-fe8e0128d5b1","To":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Action":"DOWNGRATE"}]}`)
	var f FSM
	err := json.Unmarshal(b, &f)
	if err == nil {
		t.Errorf("json is not valid fsm, it should errored")
	}

}

func TestBasic3StatesJsonUnmarshal2(t *testing.T) {
	b := []byte(`{"Name":"Customer Plan Status","Current":"TRIAL","States":["TRIAL","BASIC","PREMIUM"],"Transitions":[{"From":"TRIAL","To":"BASIC","Action":"UPGRATE"},{"From":"TRIAL","To":"PREMIUM","Action":"UPGRATE"},{"From":"BASIC","To":"PREMIUM","Action":"UPGRATE"},{"From":"PREMIUM","To":"BASIC","Action":"DOWNGRATE"}]}`)
	var f FSM
	err := json.Unmarshal(b, &f)
	if err != nil {
		t.Fatal(err)
	}
}
