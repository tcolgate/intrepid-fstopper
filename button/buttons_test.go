package button

import (
	"testing"
	"time"
)

func TestMgr_Process(t *testing.T) {
	evs := make(chan Event, 4)
	intEvs := make(chan IntEvent, 4)
	mgr := &Mgr{
		Events:    evs,
		IntEvents: intEvs,
	}

	mgr.Process(time.Now().UnixNano())
}
