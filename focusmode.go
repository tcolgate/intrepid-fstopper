package main

type focusMode struct {
	prevMode Mode
}

func (e focusMode) SwitchTo(prev Mode) {
	e.prevMode = prev
}

func (e focusMode) SwitchAway() Mode {
	return e.prevMode
}

func (e focusMode) TouchPoints() []touchPoint {
	return nil
}

func (e focusMode) Tick(passed int64) (bool, bool) {
	return false, false
}

func (e focusMode) Run() bool {
	return false
}

func (e focusMode) Focus() bool {
	// this should cancel focus mode
	return true
}

func (e focusMode) LongFocus() bool {
	// this should toggle red/white led
	return false
}

func (e focusMode) Cancel(touchPoint uint8) (bool, bool) {
	return false, true
}

func (e focusMode) Plus(touchPoint uint8) bool {
	return false
}

func (e focusMode) LongPlus(touchPoint uint8) bool {
	return false
}

func (e focusMode) Minus(touchPoint uint8) bool {
	return false
}

func (e focusMode) LongMinus(touchPoint uint8) bool {
	return false
}

func (e focusMode) UpdateDisplay(nextDisplay *[2][]byte) *touchPoint {
	copy(nextDisplay[0], stringTable[1][0])
	copy(nextDisplay[1], stringTable[1][1])

	return nil
}

func (e focusMode) NextMode() Mode {
	return nil
}
