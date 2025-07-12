package main

type Mode interface {
	TouchPoints() []touchPoint
	SwitchTo(Mode)
	SwitchAway() Mode

	Tick(passed int64) (updateDisplay bool, exit bool)
	UpdateDisplay(*[2][]byte) *touchPoint

	Plus(touchPointaIndex uint8) bool
	LongPlus(touchPointaIndex uint8) bool
	Minus(touchPointaIndex uint8) bool
	LongMinus(touchPointaIndex uint8) bool

	Run() bool
	Focus() bool
	LongFocus() bool
	Cancel(touchPointIndex uint8) (updateDisplay bool, exit bool)
}

type baseMode struct {
}
