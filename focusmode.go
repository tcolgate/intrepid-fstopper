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

func (e *focusMode) UpdateDisplay(nextDisplay *[2][16]byte) {
	nextDisplay[0] = stringTable[3]
	nextDisplay[1] = stringTable[0]
}
