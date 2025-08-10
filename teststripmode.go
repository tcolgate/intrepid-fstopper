package main

import "intrepidfstopper/num"

type testStripMode struct {
	prevMode *Mode
	nextMode *Mode

	state *stateData
}

func newTestStripMode(s *stateData) *Mode {
	m := &testStripMode{
		state: s,
	}

	return &Mode{
		TouchPoints:     m.TouchPoints,
		SwitchTo:        m.SwitchTo,
		SwitchAway:      m.SwitchAway,
		UpdateDisplay:   m.UpdateDisplay,
		PressPlus:       m.PressPlus,
		PressMinus:      m.PressMinus,
		PressLongPlus:   m.PressLongPlus,
		PressLongMinus:  m.PressLongMinus,
		PressRun:        m.PressRun,
		PressFocus:      m.PressFocus,
		PressLongFocus:  m.PressLongFocus,
		PressCancel:     m.PressCancel,
		PressLongCancel: m.PressLongCancel,
		PressMode:       m.PressMode,
		PressLongMode:   m.PressLongMode,
	}
}

func (e *testStripMode) SwitchTo(prev *Mode) {
	e.prevMode = prev
	e.state.exposureSet.isTest = true
}

func (e *testStripMode) SwitchAway() *Mode {
	return e.nextMode
}

func (e *testStripMode) TouchPoints() []touchPoint {
	return touchPoints[1]
}

func (e *testStripMode) PressRun() (bool, bool) {
	// we should set up a set of exposures here
	// - duration
	// - type (e.g. regular vs freehand)
	// - things about the LED (e.g. brightness)

	for i := range e.state.exposureSet.exposures {
		switch e.state.exposureSet.exposures[i].expUnit {
		case expUnitOff:
		case expUnitFreeHand:
		default:
			for j := range e.state.exposureSet.exposures[i].colVals {
				e.state.exposureSet.exposures[i].colTime[j] = expUnitToS(
					e.state.exposureSet.baseTime,
					e.state.exposureSet.exposures[i].expUnit,
					e.state.exposureSet.exposures[i].colVals[j],
				)
			}
		}
	}

	e.nextMode = e.state.exposureMode
	return true, true
}

func (e *testStripMode) PressFocus() (bool, bool) {
	e.state.focusColour = false
	e.nextMode = e.state.focusMode
	return true, true
}

func (e *testStripMode) PressLongFocus() (bool, bool) {
	e.state.focusColour = true
	e.nextMode = e.state.focusMode
	return true, true
}

func (e *testStripMode) PressCancel(touchPointIndex uint8) (bool, bool) {
	return false, false
}

func (e *testStripMode) PressLongCancel(touchPointIndex uint8) (bool, bool) {
	return false, false
}

func (e *testStripMode) PressPlus(touchPointIndex uint8) (bool, bool) {
	return e.state.exposureSet.tpAdjustExposureSet(touchPointIndex, 0, 0, false, false), false
}

func (e *testStripMode) PressLongPlus(touchPointIndex uint8) (bool, bool) {
	return e.state.exposureSet.tpAdjustExposureSet(touchPointIndex, 0, 0, true, false), false
}

func (e *testStripMode) PressMinus(touchPointIndex uint8) (bool, bool) {
	return e.state.exposureSet.tpAdjustExposureSet(touchPointIndex, 0, 0, false, true), false
}

func (e *testStripMode) PressLongMinus(touchPointIndex uint8) (bool, bool) {
	return e.state.exposureSet.tpAdjustExposureSet(touchPointIndex, 0, 0, true, true), false
}

func (e *testStripMode) UpdateDisplay(nextDisplay *[2][16]byte) {
	nb := &num.NumBuf{}
	nextDisplay[0] = stringTable[5]
	nextDisplay[1] = stringTable[6]

	num.Out(nb, num.Num(e.state.exposureSet.baseTime))
	copy(nextDisplay[0][0:4], nb[0:4])

	// update method
	for i := uint8(0); i <= 2; i++ {
		c := byte(' ')
		if i <= e.state.exposureSet.testStrip.steps {
			c = byte('-')
		}

		nextDisplay[1][3-(i+1)] = c
		nextDisplay[1][3+(i+1)] = c
	}

	// update method
	copy(nextDisplay[1][12:16], testMethodStrs[e.state.exposureSet.testStrip.method])

	currExp := e.state.exposureSet.testStrip.exposure

	absExpFact := currExp.colVals[0]
	if absExpFact < 0 {
		absExpFact = absExpFact * -1
	}

	switch currExp.expUnit {
	case expUnitAbsolute:
		num.Out(nb, num.Num(absExpFact))
	default:
		num.IntOut(nb, num.Num(absExpFact))
	}

	copy(nextDisplay[0][8:12], nb[0:4])

	copy(nextDisplay[0][12:16], expUnitNames[currExp.expUnit][0:4])
}

func (e *testStripMode) PressMode() (bool, bool) {
	e.nextMode = e.state.printMode
	return true, true
}

func (e *testStripMode) PressLongMode() (bool, bool) {
	return false, false
}
