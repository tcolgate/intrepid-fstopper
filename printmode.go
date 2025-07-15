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
		TouchPoints:    m.TouchPoints,
		SwitchTo:       m.SwitchTo,
		SwitchAway:     m.SwitchAway,
		UpdateDisplay:  m.UpdateDisplay,
		PressPlus:      m.PressPlus,
		PressLongPlus:  m.PressLongPlus,
		PressMinus:     m.PressMinus,
		PressLongMinus: m.PressLongMinus,
		PressRun:       m.PressRun,
		PressFocus:     m.PressFocus,
		PressLongFocus: m.PressLongFocus,
		PressCancel:    m.PressCancel,
	}
}

func (e *printMode) SwitchTo(prev *Mode) {
	e.prevMode = prev
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

func (e *printMode) PressCancel(touchPoint uint8) (bool, bool) {
	// should reset stuff and/or delete the current
	// exposure
	return false, false
}

func (e *printMode) PressPlus(touchPointIndex uint8) (bool, bool) {
	switch touchPointIndex {
	case 0:
		e.state.exposureSet.adjustBaseTime(10)
	case 1:
		e.state.exposureSet.adjustExposureTime(e.activeExposure, 0, 10)
	case 2:
		e.state.exposureSet.cycleExpUnit(e.activeExposure, true)
	default:
		return false, false
	}
	return true, false
}

func (e *printMode) PressLongPlus(touchPointIndex uint8) (bool, bool) {
	switch touchPointIndex {
	case 0:
		e.state.exposureSet.adjustBaseTime(100)
	case 1:
		e.state.exposureSet.adjustExposureTime(e.activeExposure, 0, 100)
	case 2:
		e.state.exposureSet.cycleExpUnit(e.activeExposure, true)
	default:
		return false, false
	}
	return true, false
}

func (e *printMode) PressMinus(touchPointIndex uint8) (bool, bool) {
	switch touchPointIndex {
	case 0:
		e.state.exposureSet.adjustBaseTime(-10)
	case 1:
		e.state.exposureSet.adjustExposureTime(e.activeExposure, 0, -10)
	case 2:
		e.state.exposureSet.cycleExpUnit(e.activeExposure, false)
	default:
		return false, false
	}
	return true, false
}

func (e *printMode) PressLongMinus(touchPointIndex uint8) (bool, bool) {
	switch touchPointIndex {
	case 0:
		e.state.exposureSet.adjustBaseTime(-100)
	case 1:
		e.state.exposureSet.adjustExposureTime(e.activeExposure, 0, -100)
	case 2:
		e.state.exposureSet.cycleExpUnit(e.activeExposure, false)
	default:
		return false, false
	}
	return true, false
}

func (e *printMode) UpdateDisplay(nextDisplay *[2][16]byte) {
	nb := &num.NumBuf{}
	nextDisplay[0] = stringTable[0][0]
	nextDisplay[1] = stringTable[0][1]

	num.Out(nb, num.Num(e.state.exposureSet.baseTime))
	copy(nextDisplay[0][0:4], nb[0:4])

	currExpIndex := 0
	nextDisplay[1][13] = byte('1' + currExpIndex)
	nextDisplay[1][15] = byte('0' + maxExposures)

	currExp := e.state.exposureSet.exposures[currExpIndex]
	switch currExp.expUnit {
	case expUnitOff, expUnitFreeHand:
		nextDisplay[0][6] = byte(' ')
		copy(nextDisplay[1][2:6], []byte(`    `))
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
	}

	copy(nextDisplay[0][12:16], expUnitNames[currExp.expUnit][0:4])
}
