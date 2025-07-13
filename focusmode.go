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
	e.state.SetLEDPanel(e.state.focusColour)
	e.prevMode = prev
}

func (e *focusMode) SwitchAway() *Mode {
	e.state.focusColour = ledOff
	e.state.SetLEDPanel(e.state.focusColour)

	return e.prevMode
}

func (e *focusMode) PressFocus() (bool, bool) {
	return true, true
}

func (e *focusMode) PressLongFocus() (bool, bool) {
	switch e.state.focusColour {
	case ledRed:
		e.state.focusColour = ledWhite
	default:
		e.state.focusColour = ledRed
	}
	e.state.SetLEDPanel(e.state.focusColour)
	return false, false
}

func (e *focusMode) PressCancel(touchPoint uint8) (bool, bool) {
	return true, true
}

func (e *focusMode) UpdateDisplay(nextDisplay *[2][]byte) {
	copy(nextDisplay[0], stringTable[1][0])
	copy(nextDisplay[1], stringTable[1][1])
}
