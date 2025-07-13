package main

import "intrepidfstopper/num"

type printMode struct {
	prevMode *Mode
	nextMode *Mode

	state *stateData

	baseTime           uint64
	paused             bool
	running            bool
	remainingTime      int64
	exposureFactor     int32
	exposureFactorUnit expUnit
}

func newBWMode(s *stateData) *Mode {
	m := &printMode{
		state: s,

		baseTime:           7_00,
		exposureFactorUnit: 1, // default to 1/2 stops
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
	println("printmode switchaway, e", e)
	return e.nextMode
}

func (e *printMode) TouchPoints() []touchPoint {
	return touchPoints[0]
}

func (e *printMode) PressRun() (bool, bool) {
	e.nextMode = e.state.exposureMode
	return true, true
}

func (e *printMode) PressFocus() (bool, bool) {
	println("printmode focus pressed, e: ", e.state)
	e.state.focusColour = ledRed
	e.nextMode = e.state.focusMode
	return true, true
}

func (e *printMode) PressLongFocus() (bool, bool) {
	println("printmode focus long pressed")
	e.state.focusColour = ledWhite
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
		if e.baseTime != 25500 {
			e.baseTime += 10
		}
	case 1:
		if e.exposureFactor != 126 {
			e.exposureFactor += 1
		}
	case 2:
		e.exposureFactorUnit++
		if e.exposureFactorUnit > 4 {
			e.exposureFactorUnit = 0
		}
	default:
		return false, false
	}
	return true, false
}

func (e *printMode) PressLongPlus(touchPointIndex uint8) (bool, bool) {
	switch touchPointIndex {
	case 0:
		if e.baseTime != 25500 {
			e.baseTime += 10
		}
	case 1:
		if e.exposureFactor != 126 {
			e.exposureFactor += 1
		}
	case 2:
		e.exposureFactorUnit++
		if e.exposureFactorUnit > 4 {
			e.exposureFactorUnit = 0
		}
	default:
		return false, false
	}
	return true, false
}

func (e *printMode) PressMinus(touchPointIndex uint8) (bool, bool) {
	switch touchPointIndex {
	case 0:
		if e.baseTime != 0 {
			e.baseTime -= 10
		}
	case 1:
		if e.exposureFactor != -126 {
			e.exposureFactor -= 1
		}
	case 2:
		e.exposureFactorUnit--
		if e.exposureFactorUnit == 0 {
			e.exposureFactorUnit = 4
		}
	default:
		return false, false
	}
	return true, false
}

func (e *printMode) PressLongMinus(touchPointIndex uint8) (bool, bool) {
	switch touchPointIndex {
	case 0:
		if e.baseTime != 0 {
			e.baseTime -= 10
		}
	case 1:
		if e.exposureFactor != -126 {
			e.exposureFactor -= 1
		}
	case 2:
		e.exposureFactorUnit--
		if e.exposureFactorUnit < 0 {
			e.exposureFactorUnit = 4
		}
	default:
		return false, false
	}
	return true, false
}

func (e *printMode) UpdateDisplay(nextDisplay *[2][]byte) {
	nb := &num.NumBuf{}
	copy(nextDisplay[0], stringTable[0][0])
	copy(nextDisplay[1], stringTable[0][1])

	num.Out(nb, num.Num(e.baseTime))
	copy(nextDisplay[0][2:6], nb[0:4])

	if e.exposureFactor < 0 {
		nextDisplay[1][1] = signMinus
	} else {
		nextDisplay[1][1] = signPlus
	}
	absExpFact := e.exposureFactor
	if absExpFact < 0 {
		absExpFact = absExpFact * -1
	}
	num.OutLeft(nb, num.Num(absExpFact))
	copy(nextDisplay[1][2:6], nb[0:4])

	copy(nextDisplay[0][10:15], expUnitNames[e.exposureFactorUnit][0:4])

	res := num.Num(11_23)
	num.Out(nb, num.Num(res))
	resLen := num.Len(res)
	nextDisplay[1][13-resLen] = []byte("(")[0]
	nextDisplay[1][14-resLen] = []byte("=")[0]
	copy(nextDisplay[1][11:15], nb[0:4])
	nextDisplay[1][15] = []byte(")")[0]
}
