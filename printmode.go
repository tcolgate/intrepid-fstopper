package main

import "intrepidfstopper/num"

type printMode struct {
	prevMode *Mode
	nextMode *Mode

	state *stateData

	// which exposure are we edditing
	activeExposure uint8
}

func newBWMode(s *stateData) *Mode {
	m := &printMode{
		state: s,
	}

	return &Mode{
		TouchPoints:     m.TouchPoints,
		SwitchTo:        m.SwitchTo,
		SwitchAway:      m.SwitchAway,
		UpdateDisplay:   m.UpdateDisplay,
		PressPlus:       m.PressPlus,
		PressLongPlus:   m.PressLongPlus,
		PressMinus:      m.PressMinus,
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

func (e *printMode) SwitchTo(prev *Mode) {
	e.prevMode = prev
	e.state.exposureSet.isTest = false
}

func (e *printMode) SwitchAway() *Mode {
	return e.nextMode
}

func (e *printMode) TouchPoints() []touchPoint {
	return touchPoints[0]
}

func (e *printMode) PressRun() (bool, bool) {
	// we should set up a set of exposures here
	// - duration
	// - type (e.g. regular vs freehand)
	// - things about the LED (e.g. brightness)

	e.nextMode = e.state.exposureMode
	return true, true
}

func (e *printMode) PressFocus() (bool, bool) {
	e.state.focusColour = false
	e.nextMode = e.state.focusMode
	return true, true
}

func (e *printMode) PressLongFocus() (bool, bool) {
	e.state.focusColour = true
	e.nextMode = e.state.focusMode
	return true, true
}

func (e *printMode) PressCancel(touchPointIndex uint8) (bool, bool) {
	switch touchPointIndex {
	case 0:
		// Quick press of cancel on the base time set resets the basetime
		// to whatever the calculated value for the current dsiplayed
		// baseTime and adjustment would be
		e.state.exposureSet.baseTime = expUnitToS(
			e.state.exposureSet.baseTime,
			e.state.exposureSet.exposures[e.activeExposure].expUnit,
			e.state.exposureSet.exposures[e.activeExposure].colVals[0],
		)
		e.state.exposureSet.exposures[e.activeExposure].colVals[0] = 0

		return true, false
	case 1:
		e.state.exposureSet.exposures[e.activeExposure].colVals[0] = 0
		return true, false
	default:
		return false, false
	}
}

func (e *printMode) PressLongCancel(touchPointIndex uint8) (bool, bool) {
	e.state.exposureSet.baseTime = 7_00
	e.state.exposureSet.exposures[e.activeExposure].colVals[0] = 0
	return true, false
}

func (e *printMode) PressPlus(touchPointIndex uint8) (bool, bool) {
	return e.state.exposureSet.tpAdjustExposureSet(touchPointIndex, e.activeExposure, 0, false, false), false
}

func (e *printMode) PressLongPlus(touchPointIndex uint8) (bool, bool) {
	return e.state.exposureSet.tpAdjustExposureSet(touchPointIndex, e.activeExposure, 0, true, false), false
}

func (e *printMode) PressMinus(touchPointIndex uint8) (bool, bool) {
	return e.state.exposureSet.tpAdjustExposureSet(touchPointIndex, e.activeExposure, 0, false, true), false
}

func (e *printMode) PressLongMinus(touchPointIndex uint8) (bool, bool) {
	return e.state.exposureSet.tpAdjustExposureSet(touchPointIndex, e.activeExposure, 0, true, true), false
}

func (e *printMode) UpdateDisplay(nextDisplay *[2][16]byte) {
	nb := &num.NumBuf{}
	nextDisplay[0] = stringTable[1]
	nextDisplay[1] = stringTable[2]

	num.Out(nb, num.Num(e.state.exposureSet.baseTime))
	copy(nextDisplay[0][0:4], nb[0:4])

	currExpIndex := 0
	nextDisplay[1][13] = byte('1' + currExpIndex)
	nextDisplay[1][15] = byte('0' + maxExposures)

	currExp := e.state.exposureSet.exposures[currExpIndex]
	switch currExp.expUnit {
	case expUnitOff, expUnitFreeHand:
		nextDisplay[0][6] = byte(' ')
		copy(nextDisplay[1][0:7], []byte(`        `))
	default:
		if currExp.colVals[0] < 0 {
			nextDisplay[0][6] = signMinus
		} else {
			nextDisplay[0][6] = signPlus
		}

		absExpFact := currExp.colVals[0]
		if absExpFact < 0 {
			absExpFact = absExpFact * -1
		}
		switch currExp.expUnit {
		case expUnitAbsolute:
			num.OutLeft(nb, num.Num(absExpFact))
		default:
			num.IntOutLeft(nb, num.Num(absExpFact))
		}
		copy(nextDisplay[0][7:11], nb[0:4])

		currTime := expUnitToS(
			e.state.exposureSet.baseTime,
			e.state.exposureSet.exposures[e.activeExposure].expUnit,
			e.state.exposureSet.exposures[e.activeExposure].colVals[0],
		)
		num.OutLeft(nb, num.Num(currTime))
		copy(nextDisplay[1][2:6], nb[0:4])
	}

	copy(nextDisplay[0][12:16], expUnitNames[currExp.expUnit][0:4])
}

func (e *printMode) PressMode() (bool, bool) {
	e.nextMode = e.state.tsMode
	return true, true
}

func (e *printMode) PressLongMode() (bool, bool) {
	// This should toggle between BW, Tri_color and RGB
	return false, false
}
