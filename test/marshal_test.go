package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	fsm "github.com/lemenendez/fsm"
)

func TestBasic3StatesJson(t *testing.T) {

	const TRIAL = "TRIAL"
	const BASIC = "BASIC"
	const PREMIUM = "PREMIUM"
	const UPGRATE = "UPGRATE"
	const DOWNGRATE = "DOWNGRATE"

	f := fsm.NewFSM("Customer Plan")

	f.AddState(TRIAL)
	f.AddState(BASIC)
	f.AddState(PREMIUM)

	f.AddTrans(TRIAL, BASIC, UPGRATE)
	f.AddTrans(TRIAL, PREMIUM, UPGRATE)
	f.AddTrans(BASIC, PREMIUM, UPGRATE)
	f.AddTrans(PREMIUM, BASIC, DOWNGRATE)

	f.Init(TRIAL)

	b, err := json.Marshal(f)
	if err != nil {
		fmt.Println(err)
		return
	}
	_ = ioutil.WriteFile("ThreeStateFST.json", b, 0644)
	//t.Log(string(b))

	var f2 fsm.FSM
	err = json.Unmarshal(b, &f2)

	myFunc := func(pre string, cur string, action string) {
		t.Logf("Previous State:%v, New State:%v, Action:%v", pre, cur, action)
	}

	if err != nil {
		err = f2.Exec(UPGRATE, BASIC, myFunc)
		if err != nil {
			t.Logf(err.Error())
			t.Errorf("It should should works")
		}
	}
}

func TestBasic3StatesJsonUnmarshal1(t *testing.T) {
	b := []byte(`{"Id":"0","UUID":"642f24f0-132d-4413-80b9-0fc43408c302","Name":"Customer Plan Status","Current":"TRIAL","States":[{"Id":0,"UUID":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","Name":"TRIAL"},{"Id":1,"UUID":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Name":"BASIC"},{"Id":2,"UUID":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Name":"PREMIUM"}],"Transitions":[{"From":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","To":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Action":"UPGRATE"},{"From":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","To":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Action":"UPGRATE"},{"From":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","To":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Action":"UPGRATE"},{"From":"c7560801-7ade-4a73-866b-fe8e0128d5b1","To":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Action":"DOWNGRATE"}]}`)
	var f fsm.FSM
	err := json.Unmarshal(b, &f)
	if err == nil {
		t.Errorf("Id must be integer not string")
	}

}

func TestBasic3StatesJsonUnmarshal2(t *testing.T) {
	b := []byte(`{"Id":0,"UUID":"642f24f0-132d-4413-80b9-0fc43408cZZZ","Name":"Customer Plan Status","Current":"TRIAL","States":[{"Id":0,"UUID":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","Name":"TRIAL"},{"Id":1,"UUID":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Name":"BASIC"},{"Id":2,"UUID":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Name":"PREMIUM"}],"Transitions":[{"From":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","To":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Action":"UPGRATE"},{"From":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","To":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Action":"UPGRATE"},{"From":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","To":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Action":"UPGRATE"},{"From":"c7560801-7ade-4a73-866b-fe8e0128d5b1","To":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Action":"DOWNGRATE"}]}`)
	var f fsm.FSM
	err := json.Unmarshal(b, &f)
	if err == nil {
		t.Errorf("State in transition does not exists")
	}
	t.Log(err)

}

func TestBasic3StatesJsonUnmarshal3(t *testing.T) {
	b := []byte(`{"Id":0,"UUID":"642f24f0-132d-4413-80b9-0fc43408c302","Name":"Customer Plan Status","Current":"TRIAL","States":[{"Id":0,"UUID":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","Name":"TRIAL"},{"Id":1,"UUID":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Name":"BASIC"},{"Id":2,"UUID":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Name":"PREMIUM"}],"Transitions":[{"From":"c73d9ad0-89fb-4117-baab-5c3751ae6315","To":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Action":"UPGRATE"},{"From":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","To":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Action":"UPGRATE"},{"From":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","To":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Action":"UPGRATE"},{"From":"c7560801-7ade-4a73-866b-fe8e0128d5b1","To":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Action":"DOWNGRATE"}]}`)
	var f fsm.FSM
	err := json.Unmarshal(b, &f)
	if err == nil {
		t.Errorf("State in transition does not exists")
	}
	t.Log(err)
}

func TestBasic3StatesJsonUnmarshal4(t *testing.T) {
	b := []byte(`{"Id":0,"UUID":"642f24f0-132d-4413-80b9-0fc43408c302","Name":"Customer Plan Status","Current":"TRIAL","States":[{"Id":0,"UUID":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","Name":"TRIAL"},{"Id":1,"UUID":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Name":"BASIC"},{"Id":2,"UUID":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Name":"PREMIUM"}],"Transitions":[{"From":"8d6afd9e-c260-4ef8-b678-0baa796b60xx","To":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Action":"UPGRATE"},{"From":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","To":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Action":"UPGRATE"},{"From":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","To":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Action":"UPGRATE"},{"From":"c7560801-7ade-4a73-866b-fe8e0128d5b1","To":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Action":"DOWNGRATE"}]}`)
	var f fsm.FSM
	err := json.Unmarshal(b, &f)
	if err == nil {
		t.Errorf("State UUID is not valid length")
	}
	t.Log(err)
}

func TestBasic3StatesJsonUnmarshal5(t *testing.T) {
	b := []byte(`{"Id":0,"UUID":"642f24f0-132d-4413-80b9-0fc43408c302","Name":"Customer Plan Status","Current":"TRIAL","States":[{"Id":0,"UUID":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","Name":"TRIAL"},{"Id":1,"UUID":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Name":"BASIC"},{"Id":2,"UUID":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Name":"PREMIUM"}],"Transitions":[{"From":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","To":"4b4fc77c-6319-4ebe-863c-ab46e7611dxx","Action":"UPGRATE"},{"From":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","To":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Action":"UPGRATE"},{"From":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","To":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Action":"UPGRATE"},{"From":"c7560801-7ade-4a73-866b-fe8e0128d5b1","To":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Action":"DOWNGRATE"}]}`)
	var f fsm.FSM
	err := json.Unmarshal(b, &f)
	if err == nil {
		t.Errorf("State in transition does not exists or UUID is invalid")
	}
	t.Log(err)
}

func TestBasic3StatesJsonUnmarshal6(t *testing.T) {
	b := []byte(`{"Id":0,"UUID":"642f24f0-132d-4413-80b9-0fc43408c302","Name":"Customer Plan Status","Current":"TRIAL","States":[{"Id":0,"UUID":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","Name":"TRIAL"},{"Id":1,"UUID":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Name":"BASIC"},{"Id":2,"UUID":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Name":"PREMIUM"}],"Transitions":[{"From":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","To":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Action":"UPGRATE"},{"From":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","To":"98c4baca-b26c-452b-89e8-519fe9dae87e","Action":"UPGRATE"},{"From":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","To":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Action":"UPGRATE"},{"From":"c7560801-7ade-4a73-866b-fe8e0128d5b1","To":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Action":"DOWNGRATE"}]}`)
	var f fsm.FSM
	err := json.Unmarshal(b, &f)
	if err == nil {
		t.Errorf("State in transition does not exists or UUID is invalid")
	}
	t.Log(err)
}

func TestBasic3StatesJsonUnmarshal7(t *testing.T) {
	b := []byte(`{"Id":0,"UUID":"642f24f0-132d-4413-80b9-0fc43408c302","Name":"Customer Plan Status","Current":"TRIAL","States":[{"Id":0,"UUID":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","Name":"TRIAL"},{"Id":1,"UUID":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Name":"BASIC"},{"Id":2,"UUID":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Name":"PREMIUM"}],"Transitions":[{"From":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","To":"4b4fc77c-6319-4ebe-863c-ab46e7611dXX","Action":"UPGRATE"},{"From":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","To":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Action":"UPGRATE"},{"From":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","To":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Action":"UPGRATE"},{"From":"c7560801-7ade-4a73-866b-fe8e0128d5b1","To":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Action":"DOWNGRATE"}]}`)
	var f fsm.FSM
	err := json.Unmarshal(b, &f)
	if err == nil {
		t.Errorf("State in transition does not exists or UUID is invalid")
	}
	t.Log(err)
}

func TestBasic3StatesJsonUnmarshal8(t *testing.T) {
	b := []byte(`{"Id":0,"UUID":"642f24f0-132d-4413-80b9-0fc43408c302","Name":"Customer Plan Status","Current":"TRIAL","States":[{"Id":0,"UUID":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","Name":"TRIAL"},{"Id":1,"UUID":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Name":"BASIC"},{"Id":2,"UUID":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Name":"PREMIUM"}],"Transitions":[{"From":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","To":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Action":"UPGRATE"},{"From":"8d6afd9e-c260-4ef8-b678-0baa796b60cc","To":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Action":"PLAN UPGRADE UPGRATE"},{"From":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","To":"c7560801-7ade-4a73-866b-fe8e0128d5b1","Action":"UPGRATE"},{"From":"c7560801-7ade-4a73-866b-fe8e0128d5b1","To":"4b4fc77c-6319-4ebe-863c-ab46e7611d05","Action":"DOWNGRATE"}]}`)
	var f fsm.FSM
	err := json.Unmarshal(b, &f)
	if err == nil {
		t.Errorf("State in transition does not exists or UUID is invalid")
	}
	t.Log(err)
}