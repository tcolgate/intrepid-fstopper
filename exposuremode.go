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

	displayUpdated bool

	exposures [maxExposures]int64
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

func (e *exposureMode) nextTime() bool {
	if e.activeExp == e.totalExps {
		return false
	}

	e.activeExp += 1
	e.remainingTime = e.exposures[e.activeExp-1]

	return true
}

func (e *exposureMode) SwitchTo(prev *Mode) {
	// need to get the exposure details in here
	// from the calling mode
	e.prevMode = prev
	e.activeExp = 0

	e.totalExps = e.state.exposureSet.calcInto(&e.exposures)

	e.nextTime()

	if e.totalExps > 1 {
		e.paused = true
	}

	e.running = false
	e.displayUpdated = false
}

func (e *exposureMode) SwitchAway() *Mode {
	e.state.SetLEDPanel(ledOff)

	return e.prevMode
}

func (e *exposureMode) Tick(passed int64) (bool, bool) {
	if e.paused {
		return false, false
	}

	// For very short exposures the display update can impact the
	// timing
	if !e.displayUpdated {
		e.displayUpdated = true
		return true, false
	}

	if !e.running {
		e.running = true
		e.state.SetLEDPanel(ledWhite)
		return false, false
	}

	e.remainingTime -= passed
	if e.remainingTime <= 0 {
		e.remainingTime = 0

		// exposure finished
		e.paused = false
		e.running = false
		e.state.SetLEDPanel(ledOff)

		if !e.nextTime() {
			return true, true
		}

		e.paused = true

		return true, false
	}

	// TODO: it would be better to do the update of the
	// time here to reduce rather than leaving it to
	// a full call to UpdateDisplay, since we don't
	// need to re-render the entire display.
	return false, false
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
	nextDisplay[0] = stringTable[4]
	nextDisplay[1] = stringTable[0]
	nb := num.NumBuf{}

	nextDisplay[0][11] = byte('1' + e.activeExp - 1)
	nextDisplay[0][13] = byte('1' + e.totalExps - 1)

	num.Out(&nb, num.Num(e.remainingTime/int64((10*time.Millisecond))))
	copy(nextDisplay[1][12:16], nb[0:4])
	if e.paused {
		copy(nextDisplay[1][0:7], []byte("Paused ")[0:7])
	}
}
