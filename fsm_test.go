package fsm_test

import (
	"github.com/lemenendez/fsm"
	"testing"
)

type LWFSM struct {
	name  string
	trans [][3]string
}

var FSMs = []struct {
	fsm      LWFSM
	expected error
}{
	{
		fsm: LWFSM{
			name: "BASIC",
			trans: [][3]string{
				{"ACTIVE", "INACTIVE", "ACTIVATE"},
				{"INACTIVE", "ACTIVE", "DEACTIVATE"},
			},
		},
		expected: nil,
	},
	{
		fsm: LWFSM{
			name: "BASIC",
			trans: [][3]string{
				{"NOT_READY", "READY", "CHECK"},
				{"READY", "NOT_READY", "CHECK"},
			},
		},
		expected: nil,
	},
}

func TestFull(t *testing.T) {
	for _, test := range FSMs {
		_, err := fsm.New(test.fsm.name, test.fsm.trans)
		if err != test.expected {
			t.Fail()
		}
	}
}
