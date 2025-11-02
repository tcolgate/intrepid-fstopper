// Copyright 2025 Tristan Colgate-McFarlane
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	switch e.state.exposureSet.ledMode {
	case modeRGB:
		return touchPoints[2]
	default: // modeBW
		return touchPoints[0]
	}
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
			e.state.exposureSet.exposures[e.activeExposure].colVal,
		)
		e.state.exposureSet.exposures[e.activeExposure].colVal = 0

		return true, false
	case 1:
		e.state.exposureSet.exposures[e.activeExposure].colVal = 0
		return true, false
	case 2, 3:
		if e.activeExposure == 0 {
			// Not allowed to disable exposure 0
			return false, false
		}
		e.state.exposureSet.exposures[e.activeExposure].enabled = !e.state.exposureSet.exposures[e.activeExposure].enabled
		return true, false
	default:
		return false, false
	}
}

func (e *printMode) PressLongCancel(touchPointIndex uint8) (bool, bool) {
	e.state.exposureSet.baseTime = 7_00
	e.state.exposureSet.exposures[e.activeExposure].colVal = 0
	return true, false
}

func (e *printMode) adjustActiveExposure(inc bool) (bool, bool) {
	nextExp := e.activeExposure

	if inc {
		if nextExp == (maxExposures - 1) {
			nextExp = 0
		} else {
			nextExp += 1
		}
	} else {
		if nextExp == 0 {
			nextExp = (maxExposures - 1)
		} else {
			nextExp -= 1
		}
	}

	if e.activeExposure != nextExp {
		e.activeExposure = nextExp
		return true, false
	}

	return false, false
}

func (e *printMode) PressPlus(touchPointIndex uint8) (bool, bool) {
	switch touchPointIndex {
	case 3, 5:
		return e.adjustActiveExposure(true)
	default:
		return e.state.exposureSet.tpAdjustExposureSet(touchPointIndex, e.activeExposure, false, false), false
	}
}

func (e *printMode) PressLongPlus(touchPointIndex uint8) (bool, bool) {
	switch touchPointIndex {
	case 3, 5:
		return e.adjustActiveExposure(true)
	default:
		return e.state.exposureSet.tpAdjustExposureSet(touchPointIndex, e.activeExposure, true, false), false
	}
}

func (e *printMode) PressMinus(touchPointIndex uint8) (bool, bool) {
	switch touchPointIndex {
	case 3, 5:
		return e.adjustActiveExposure(false)
	default:
		return e.state.exposureSet.tpAdjustExposureSet(touchPointIndex, e.activeExposure, false, true), false
	}
}

func (e *printMode) PressLongMinus(touchPointIndex uint8) (bool, bool) {
	switch touchPointIndex {
	case 3, 5:
		return e.adjustActiveExposure(false)
	default:
		return e.state.exposureSet.tpAdjustExposureSet(touchPointIndex, e.activeExposure, true, true), false
	}
}

func (e *printMode) updateDisplayPage2(tp uint8, nextDisplay *[2][16]byte, nb *num.NumBuf) {
	// or the RGB line
	switch e.state.exposureSet.ledMode {
	case modeBW:
		nextDisplay[0] = stringTable[8]
		copy(nextDisplay[1][0:8], []byte(`         `))

		num.IntOutLeft(nb, num.Num(e.state.exposureSet.exposures[e.activeExposure].rgb[3]))
		copy(nextDisplay[0][12:16], nb[0:4])
	case modeRGB:
		nextDisplay[0] = stringTable[9]
		copy(nextDisplay[1][0:8], []byte(`B:       `))

		num.IntOutLeft(nb, num.Num(e.state.exposureSet.exposures[e.activeExposure].rgb[0]))
		copy(nextDisplay[0][3:7], nb[0:4])

		num.IntOutLeft(nb, num.Num(e.state.exposureSet.exposures[e.activeExposure].rgb[1]))
		copy(nextDisplay[0][11:15], nb[0:4])

		num.IntOutLeft(nb, num.Num(e.state.exposureSet.exposures[e.activeExposure].rgb[2]))
		copy(nextDisplay[1][3:7], nb[0:4])
	}
}

func (e *printMode) UpdateDisplay(tp uint8, nextDisplay *[2][16]byte) {
	nb := &num.NumBuf{}

	nextDisplay[0] = stringTable[1]
	nextDisplay[1] = stringTable[2]

	nextDisplay[1][13] = byte('1' + e.activeExposure)
	nextDisplay[1][15] = byte('0' + maxExposures)

	if tp >= 4 {
		e.updateDisplayPage2(tp, nextDisplay, nb)
		return
	}

	num.Out(nb, num.Num(e.state.exposureSet.baseTime))
	copy(nextDisplay[0][0:4], nb[0:4])

	currExp := e.state.exposureSet.exposures[e.activeExposure]
	switch {
	case !currExp.enabled || currExp.expUnit == expUnitFreeHand:
		nextDisplay[0][6] = byte(' ')
		copy(nextDisplay[1][0:8], []byte(`         `))
	default:
		if currExp.colVal < 0 {
			nextDisplay[0][6] = signMinus
		} else {
			nextDisplay[0][6] = signPlus
		}

		absExpFact := currExp.colVal
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
			e.state.exposureSet.exposures[e.activeExposure].colVal,
		)
		num.OutLeft(nb, num.Num(currTime))
		copy(nextDisplay[1][2:6], nb[0:4])
	}

	if currExp.enabled {
		copy(nextDisplay[0][12:16], expUnitNames[currExp.expUnit][0:4])
	} else {
		copy(nextDisplay[0][12:16], expUnitNames[expUnitLast][0:4])
	}
}

func (e *printMode) PressMode() (bool, bool) {
	e.nextMode = e.state.tsMode
	return true, true
}

func (e *printMode) PressLongMode() (bool, bool) {
	e.state.exposureSet.cycleLEDMode()

	return true, false
}
