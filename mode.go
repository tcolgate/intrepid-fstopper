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

type Mode struct {
	TouchPoints func() []touchPoint

	SwitchTo   func(*Mode)  // enter a mode, passed the current mode (we may want to return to it)
	SwitchAway func() *Mode // exit a mode, returns the next mode we should enter

	Tick          func(passed int64) (updateDisplay bool, exit bool)
	UpdateDisplay func(page uint8, nextDisplay *[2][16]byte)

	PressPlus       func(touchPointAction tpAction) (updateDisplay bool, exit bool)
	PressLongPlus   func(touchPointAction tpAction) (updateDisplay bool, exit bool)
	PressMinus      func(touchPointAction tpAction) (updateDisplay bool, exit bool)
	PressLongMinus  func(touchPointAction tpAction) (updateDisplay bool, exit bool)
	PressCancel     func(touchPointAction tpAction) (updateDisplay bool, exit bool)
	PressLongCancel func(touchPointAction tpAction) (updateDisplay bool, exit bool)

	PressMode      func() (updateDisplay bool, exit bool)
	PressLongMode  func() (updateDisplay bool, exit bool)
	PressRun       func() (updateDisplay bool, exit bool)
	PressFocus     func() (updateDisplay bool, exit bool)
	PressLongFocus func() (updateDisplay bool, exit bool)
}
