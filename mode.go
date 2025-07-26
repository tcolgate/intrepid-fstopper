package main

type Mode struct {
	TouchPoints func() []touchPoint

	SwitchTo   func(*Mode)  // enter a mode, passed the current mode (we may want to return to it)
	SwitchAway func() *Mode // exit a mode, returns the next mode we should enter

	Tick          func(passed int32) (updateDisplay bool, exit bool)
	UpdateDisplay func(*[2][16]byte)

	PressPlus      func(touchPointIndex uint8) (updateDisplay bool, exit bool)
	PressLongPlus  func(touchPointIndex uint8) (updateDisplay bool, exit bool)
	PressMinus     func(touchPointIndex uint8) (updateDisplay bool, exit bool)
	PressLongMinus func(touchPointIndex uint8) (updateDisplay bool, exit bool)
	PressCancel    func(touchPointIndex uint8) (updateDisplay bool, exit bool)

	PressRun       func() (updateDisplay bool, exit bool)
	PressFocus     func() (updateDisplay bool, exit bool)
	PressLongFocus func() (updateDisplay bool, exit bool)
}
