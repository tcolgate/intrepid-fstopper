package main

import (
	"intrepidfstopper/num"
	"time"
)

type exposureMode struct {
	prevMode Mode

	paused        bool
	running       bool
	remainingTime int64
}

func (e exposureMode) SwitchTo(prev Mode) {
	e.prevMode = prev
}

func (e exposureMode) SwitchAway() Mode {
	return e.prevMode
}

func (e exposureMode) TouchPoints() []touchPoint {
	return nil
}

func (e exposureMode) Tick(passed int64) (bool, bool) {
	if e.paused {
		return false, false
	}

	e.remainingTime -= passed
	if e.remainingTime <= 0 {
		// exposure finished
		e.paused = false
		e.running = false
		e.remainingTime = 0
		state.currentLED = ledOff
		return true, true
	}

	return true, false
}

func (e exposureMode) Run() bool {
	e.paused = !e.paused

	if e.paused {
		state.currentLED = ledOff
	} else {
		state.currentLED = ledWhite
	}

	return true
}

func (e exposureMode) Focus() bool {
	return false
}

func (e exposureMode) LongFocus() bool {
	return false
}

func (e exposureMode) Cancel(touchPoint uint8) (bool, bool) {
	// cancel running exposure, reset
	e.paused = false
	e.running = false
	e.remainingTime = 0
	state.currentLED = ledOff

	return true, true
}

func (e exposureMode) Plus(touchPoint uint8) bool {
	return false
}

func (e exposureMode) LongPlus(touchPoint uint8) bool {
	return false
}

func (e exposureMode) Minus(touchPoint uint8) bool {
	return false
}

func (e exposureMode) LongMinus(touchPoint uint8) bool {
	return false
}

func (e exposureMode) UpdateDisplay(nextDisplay *[2][]byte) *touchPoint {
	nb := num.NumBuf{}
	copy(nextDisplay[0], stringTable[2][0])
	copy(nextDisplay[1], stringTable[2][1])

	if e.running {
		num.Out(&nb, num.Num(e.remainingTime/int64((10*time.Millisecond))))
		copy(nextDisplay[1][12:16], nb[0:4])
	}

	return nil
}
