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

type focusMode struct {
	prevMode *Mode

	// focusColour is set on state, as it needs to be
	// adjustable by the mode that saw the focsu press
	// events to select colour
	state *stateData
}

func newFocusMode(s *stateData) *Mode {
	m := &focusMode{
		state: s,
	}
	return &Mode{
		SwitchTo:      m.SwitchTo,
		SwitchAway:    m.SwitchAway,
		UpdateDisplay: m.UpdateDisplay,

		PressFocus:     m.PressFocus,
		PressLongFocus: m.PressLongFocus,
		PressCancel:    m.PressCancel,
	}
}

func (e *focusMode) SwitchTo(prev *Mode) {
	if !e.state.focusColour {
		e.state.SetLEDPanel(ledRed)
	} else {
		e.state.SetLEDPanel(ledWhite)
	}
	e.prevMode = prev
}

func (e *focusMode) SwitchAway() *Mode {
	e.state.SetLEDPanel(ledOff)

	return e.prevMode
}

func (e *focusMode) PressFocus() (bool, bool) {
	return true, true
}

func (e *focusMode) PressLongFocus() (bool, bool) {
	e.state.focusColour = !e.state.focusColour
	if !e.state.focusColour {
		e.state.SetLEDPanel(ledRed)
	} else {
		e.state.SetLEDPanel(ledWhite)
	}
	return false, false
}

func (e *focusMode) PressCancel(touchPoint uint8) (bool, bool) {
	return true, true
}

func (e *focusMode) UpdateDisplay(_ uint8, nextDisplay *[2][16]byte) {
	nextDisplay[0] = stringTable[3]
	nextDisplay[1] = stringTable[0]
}
