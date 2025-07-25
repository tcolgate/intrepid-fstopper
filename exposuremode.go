package main

import (
	"intrepidfstopper/num"
	"time"
)

type exposureMode struct {
	prevMode *Mode
	state    *stateData

	paused        bool
	running       bool
	remainingTime int64
	activeExp     uint8
	totalExps     uint8
}

func newExpMode(s *stateData) *Mode {
	m := &exposureMode{
		state: s,
	}
	return &Mode{
		SwitchTo:      m.SwitchTo,
		SwitchAway:    m.SwitchAway,
		Tick:          m.Tick,
		UpdateDisplay: m.UpdateDisplay,
		PressRun:      m.PressRun,
		PressCancel:   m.PressCancel,
	}
}

func (e *exposureMode) SwitchTo(prev *Mode) {
	// need to get the exposure details in here
	// from the calling mode
	e.prevMode = prev

	e.remainingTime = int64(e.state.exposureSet.exposures[0].colTime[0]) * int64(tick)

	expCnt := uint8(0)
	for i := range e.state.exposureSet.exposures {
		if e.state.exposureSet.exposures[i].expUnit == expUnitOff {
			break
		}
		expCnt++
	}

	e.running = true
	e.state.SetLEDPanel(ledWhite)
}

func (e *exposureMode) SwitchAway() *Mode {
	e.state.SetLEDPanel(ledOff)

	return e.prevMode
}

func (e *exposureMode) Tick(passed int64) (bool, bool) {
	if e.paused {
		return false, false
	}

	e.remainingTime -= passed
	if e.remainingTime <= 0 {
		// exposure finished
		e.paused = false
		e.running = false
		e.remainingTime = 0
		e.state.SetLEDPanel(ledOff)
		return true, true
	}

	return true, false
}

func (e *exposureMode) PressRun() (bool, bool) {
	e.paused = !e.paused

	if e.paused {
		// or, maybe optionally ledRed?
		e.state.SetLEDPanel(ledOff)
	} else {
		e.state.SetLEDPanel(ledWhite)
	}

	return true, false
}

func (e *exposureMode) PressCancel(touchPoint uint8) (bool, bool) {
	// cancel running exposure, reset
	e.paused = false
	e.running = false
	e.remainingTime = 0
	e.state.SetLEDPanel(ledOff)

	return true, true
}

func (e *exposureMode) UpdateDisplay(nextDisplay *[2][16]byte) {
	nb := num.NumBuf{}
	if !e.state.exposureSet.isTest {
		nextDisplay[0] = stringTable[2][0]
		nextDisplay[1] = stringTable[2][1]

	} else {
		nextDisplay[0] = stringTable[4][0]
		nextDisplay[1] = stringTable[4][1]
	}

	nextDisplay[0][11] = byte('1' + e.activeExp)
	nextDisplay[0][13] = byte('1' + e.totalExps)

	num.Out(&nb, num.Num(e.remainingTime/int64((10*time.Millisecond))))
	copy(nextDisplay[1][12:16], nb[0:4])
	if e.paused {
		copy(nextDisplay[1][0:7], []byte("Paused ")[0:7])
	}
}
