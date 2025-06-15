package internal

import (
	"testing"
	"time"
)

func TestButMgr_Process(t *testing.T) {
	evs := make(chan ButEvent, 4)
	intEvs := make(chan ButIntEvent, 4)
	mgr := &ButMgr{
		Events:    evs,
		IntEvents: intEvs,
	}

	mgr.Process(time.Now())
}
