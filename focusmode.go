package main

type focusMode struct {
	prevMode *Mode
	state    *stateData
}

func newFocusMode(s *stateData) *Mode {
	m := &focusMode{
		state: s,
	}
	return &Mode{
		TouchPoints:    m.TouchPoints,
		SwitchTo:       m.SwitchTo,
		SwitchAway:     m.SwitchAway,
		Tick:           m.Tick,
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

func (e *focusMode) SwitchTo(prev *Mode) {
	e.prevMode = prev
}

func (e *focusMode) SwitchAway() *Mode {
	return e.prevMode
}

func (e *focusMode) TouchPoints() []touchPoint {
	return nil
}

func (e *focusMode) Tick(passed int64) (bool, bool) {
	return false, false
}

func (e *focusMode) PressRun() bool {
	return false
}

func (e *focusMode) PressFocus() bool {
	// this should cancel focus mode
	return true
}

func (e *focusMode) PressLongFocus() bool {
	// this should toggle red/white led
	return false
}

func (e *focusMode) PressCancel(touchPoint uint8) (bool, bool) {
	return false, true
}

func (e *focusMode) PressPlus(touchPoint uint8) bool {
	return false
}

func (e *focusMode) PressLongPlus(touchPoint uint8) bool {
	return false
}

func (e *focusMode) PressMinus(touchPoint uint8) bool {
	return false
}

func (e *focusMode) PressLongMinus(touchPoint uint8) bool {
	return false
}

func (e *focusMode) UpdateDisplay(nextDisplay *[2][]byte) *touchPoint {
	copy(nextDisplay[0], stringTable[1][0])
	copy(nextDisplay[1], stringTable[1][1])

	return nil
}

func (e *focusMode) NextMode() *Mode {
	return nil
}
