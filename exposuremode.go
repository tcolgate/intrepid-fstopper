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
}

func newExpMode(s *stateData) *Mode {
	m := &exposureMode{
		state: s,
	}
	return &Mode{
		TouchPoints:   m.TouchPoints,
		SwitchTo:      m.SwitchTo,
		SwitchAway:    m.SwitchAway,
		Tick:          m.Tick,
		UpdateDisplay: m.UpdateDisplay,
		PressRun:      m.PressRun,
		PressCancel:   m.PressCancel,
	}
}

func (e *exposureMode) SwitchTo(prev *Mode) {
	e.prevMode = prev
}

func (e *exposureMode) SwitchAway() *Mode {
	return e.prevMode
}

func (e *exposureMode) TouchPoints() []touchPoint {
	return nil
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
		//state.currentLED = ledOff
		return true, true
	}

	return true, false
}

/*
	case button.Run:
		if s.exposureRunning {
			s.exposurePaused = !s.exposurePaused
			if s.exposurePaused {
				state.currentLED = ledOff
			} else {
				state.currentLED = ledWhite
			}
			return true
		}

		// start exposure
		s.remainingTime = int64(s.baseTime) * tick
		s.exposureRunning = true
		s.exposurePaused = false
		state.currentLED = ledWhite
		return true
	case button.Cancel:
		if s.exposureRunning {
			s.exposurePaused = false
			s.exposureRunning = false
			s.remainingTime = 0
			state.currentLED = ledOff
			return true
			// stop exposure, reset time
		}
		if state.currentMode == modeFocus {
			state.currentMode = state.lastMode
			state.currentSubMode = state.lastSubMode
			clearStateBit(&state.flags, statebitFocusColour)
			state.currentLED = ledOff
			return true
		}
*/

func (e *exposureMode) PressRun() (bool, bool) {
	e.paused = !e.paused

	if e.paused {
		//		state.currentLED = ledOff
	} else {
		//		state.currentLED = ledWhite
	}

	return false, false
}

func (e *exposureMode) PressCancel(touchPoint uint8) (bool, bool) {
	// cancel running exposure, reset
	e.paused = false
	e.running = false
	e.remainingTime = 0
	//state.currentLED = ledOff

	return true, true
}

func (e *exposureMode) UpdateDisplay(nextDisplay *[2][]byte) {
	nb := num.NumBuf{}
	copy(nextDisplay[0], stringTable[2][0])
	copy(nextDisplay[1], stringTable[2][1])

	if e.running {
		num.Out(&nb, num.Num(e.remainingTime/int64((10*time.Millisecond))))
		copy(nextDisplay[1][12:16], nb[0:4])
	}
}
