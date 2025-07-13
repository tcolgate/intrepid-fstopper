package main

type Mode struct {
	TouchPoints func() []touchPoint
	SwitchTo    func(*Mode)
	SwitchAway  func() *Mode

	Tick          func(passed int64) (updateDisplay bool, exit bool)
	UpdateDisplay func(*[2][]byte) *touchPoint

	PressPlus      func(touchPointaIndex uint8) bool
	PressLongPlus  func(touchPointaIndex uint8) bool
	PressMinus     func(touchPointaIndex uint8) bool
	PressLongMinus func(touchPointaIndex uint8) bool

	PressRun       func() bool
	PressFocus     func() bool
	PressLongFocus func() bool
	PressCancel    func(touchPointIndex uint8) (updateDisplay bool, exit bool)
}
