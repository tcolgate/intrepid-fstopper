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

import (
	"intrepidfstopper/button"
	"machine"
	"time"

	"tinygo.org/x/drivers/hd44780i2c"
	"tinygo.org/x/drivers/ws2812"
)

const (
	ledCount = 56

	halfStop    = 141 // = 100 * (2 ^ (1 / 2))
	negHalfStop = 71  // = 100 * 1/(2 ^ (1 / 2))

	thirdStop    = 125 // = 100 * 1/(2 ^ (1 / 3))
	negThirdStop = 79

	tenthStop    = 107 // = 100 * (2 ^ (1 / 10))
	negTenthStop = 93  // = 100 * 1/(2 ^ (1 / 10))

	longPress = 1 * time.Second

	tick = uint32(10 * time.Millisecond)

	maxExposures = 9
)

var (
	ledOff   = [4]uint8{0, 0, 0, 0}
	ledWhite = [4]uint8{0, 0, 0, 255}
	ledRed   = [4]uint8{0, 255, 0, 0}

	signMinus = []byte(`-`)[0]
	signPlus  = []byte(`+`)[0]

	charPlusMinus = []byte{
		0b_01000,
		0b_11100,
		0b_01001,
		0b_00010,
		0b_00100,
		0b_01000,
		0b_10000,
		0b_00111,
	}
	charPlusMinusAddr = uint8(0)
)

type ledMode uint8

const (
	modeBW ledMode = iota
	modeRGB
)

type subMode uint8

const (
	modePrint subMode = iota
	modeTestStrip
)

type expUnit uint8

const (
	expUnitAbsolute expUnit = iota
	expUnitHalfStop
	expUnitThirdStop
	expUnitTenthStop
	expUnitPercent
	expUnitFreeHand
	expUnitLast
)

var (
	expUnitNames = [7][4]byte{
		[4]byte([]byte(`s   `)),
		[4]byte([]byte("/2\xDF ")),
		[4]byte([]byte("/3\xDF ")),
		[4]byte([]byte("/10\xDF")),
		[4]byte([]byte(`%   `)),
		[4]byte([]byte(`Free`)),
		[4]byte([]byte(`Off `)),
	}
)

type stateData struct {
	nextTick int64
	prevTick int64

	pots [4]uint16

	activeTouchPoints     []touchPoint
	activeTouchPointIndex uint8

	activeDisplay bool
	display1      [2][16]byte
	display2      [2][16]byte

	exposureMode *Mode
	focusMode    *Mode
	printMode    *Mode
	tsMode       *Mode
	activeMode   *Mode

	// Focus light settings
	// printMode/teststripMode write focusColor on long press
	// focusMode reads and writes
	lastColour  [4]uint8
	focusColour bool

	// printMode/testStripMode writes, exposureMode reads
	exposureSet exposureSet
}

func (s *stateData) ButtonHoldRepeat(b button.Button) (bool, bool) {
	switch b {
	case button.Plus:
		if state.activeMode.PressLongPlus != nil {
			return s.activeMode.PressLongPlus(s.activeTouchPoints[s.activeTouchPointIndex].action)
		}
	case button.Minus:
		if state.activeMode.PressLongMinus != nil {
			return s.activeMode.PressLongMinus(s.activeTouchPoints[s.activeTouchPointIndex].action)
		}
	}
	return false, false
}

func (s *stateData) ButtonPress(b button.Button) (bool, bool) {
	switch b {
	case button.Plus:
		if state.activeMode.PressPlus != nil {
			return s.activeMode.PressPlus(s.activeTouchPoints[s.activeTouchPointIndex].action)
		}
	case button.Minus:
		if state.activeMode.PressMinus != nil {
			return s.activeMode.PressMinus(s.activeTouchPoints[s.activeTouchPointIndex].action)
		}
	case button.Run:
		if state.activeMode.PressRun != nil {
			return s.activeMode.PressRun()
		}

	case button.Cancel:
		if state.activeMode.PressCancel != nil {
			tpa := tpExposure
			if s.activeTouchPoints != nil {
				tpa = s.activeTouchPoints[s.activeTouchPointIndex].action
			}
			return s.activeMode.PressCancel(tpa)
		}

	case button.Focus:
		if state.activeMode.PressFocus != nil {
			return s.activeMode.PressFocus()
		}

	case button.Mode:
		if state.activeMode.PressMode != nil {
			return s.activeMode.PressMode()
		}
	}
	return false, false
}

func (s *stateData) ButtonLongPress(b button.Button) (bool, bool) {
	switch b {
	case button.Focus:
		if s.activeMode.PressLongFocus != nil {
			return s.activeMode.PressLongFocus()
		}
	case button.Cancel:
		if s.activeMode.PressLongCancel != nil {
			tpa := tpExposure
			if s.activeTouchPoints != nil {
				tpa = s.activeTouchPoints[s.activeTouchPointIndex].action
			}
			return s.activeMode.PressLongCancel(tpa)
		}
	case button.Mode:
		if s.activeMode.PressLongMode != nil {
			return s.activeMode.PressLongMode()
		}
	}
	return false, false
}

func (s *stateData) UpdateDisplay() {
	lastDisplay := &s.display1
	nextDisplay := &s.display2
	if s.activeDisplay {
		lastDisplay = &s.display2
		nextDisplay = &s.display1
	}
	s.activeDisplay = !s.activeDisplay

	if len(s.activeTouchPoints) > 0 {
		s.activeMode.UpdateDisplay(s.activeTouchPoints[s.activeTouchPointIndex].action, nextDisplay)
	} else {
		// problematic
		s.activeMode.UpdateDisplay(tpExposure, nextDisplay)
	}

	for i := uint8(0); i < 2; i++ {
		if lastDisplay[i] != nextDisplay[i] {
			lcd.SetCursor(0, i)
			lcd.Print(nextDisplay[i][:])
			lastDisplay[i] = nextDisplay[i]
		}
	}

	if len(s.activeTouchPoints) > 0 {
		tp := s.activeTouchPoints[s.activeTouchPointIndex]
		lcd.SetCursor(tp.x, tp.y)
		lcd.CursorOn(true)
	} else {
		lcd.CursorOn(false)
	}
}

type tpAction uint8
type touchPoint struct {
	y      uint8
	x      uint8
	action tpAction
}

const (
	tpBaseTime tpAction = iota
	tpExpVal
	tpExpUnit
	tpExposure
	tpRGBG
	tpRGBR
	tpRGBB
	tpRGBW
	tpTSStrips
	tpTSMode
)

var (
	// hardware setup

	buttonPins = []machine.Pin{
		machine.D7,  // T+
		machine.D8,  // T-
		machine.D9,  // Run
		machine.D10, // Focus
		machine.D2,  // Cancel
		machine.D11, // Mode
		machine.D12, // Safelight
	}

	buttonPinsConfig = machine.PinConfig{Mode: machine.PinInputPullup}

	ledPin       = machine.PD4
	ledPinConfig = machine.PinConfig{Mode: machine.PinOutput}
	ledDriver    = ws2812.NewSK6812(ledPin)

	pots = [4]machine.ADC{
		machine.ADC{machine.ADC0},
		machine.ADC{machine.ADC1},
		machine.ADC{machine.ADC2},
		machine.ADC{machine.ADC3},
	}

	i2c       = machine.I2C0
	i2cConfig = machine.I2CConfig{
		Frequency: machine.TWI_FREQ_400KHZ,
	}

	lcdAddr   = uint8(0x27)
	lcd       = hd44780i2c.New(i2c, lcdAddr)
	lcdConfig = hd44780i2c.Config{
		Width:  16,
		Height: 2,
	}

	stringTable = [][16]byte{
		[16]byte([]byte("                ")),
		[16]byte([]byte("    s           ")),
		[16]byte([]byte("(=    s)   E: / ")),
		[16]byte([]byte("     Focus      ")),
		[16]byte([]byte("  Exposure  /   ")),
		[16]byte([]byte("    s \000         ")),
		[16]byte([]byte("   |            ")),
		[16]byte([]byte("Test Patch  /   ")),
		[16]byte([]byte("Brightness:     ")),
		[16]byte([]byte("R:       G:     ")),
		[16]byte([]byte("B:         E: / ")),
	}
	touchPoints = [][]touchPoint{
		[]touchPoint{
			{0, 3, tpBaseTime},
			{0, 7, tpExpVal},
			{0, 12, tpExpUnit},
			{1, 13, tpExposure},
			{0, 12, tpRGBW},
			{1, 13, tpExposure},
		}, // Print mode - ledBW
		[]touchPoint{
			{0, 3, tpBaseTime},
			{0, 9, tpExpVal},
			{0, 12, tpExpUnit},
			{1, 3, tpTSStrips},
			{1, 12, tpTSMode},
		}, // Test strip mode - brightness
		[]touchPoint{
			{0, 3, tpBaseTime},
			{0, 7, tpExpVal},
			{0, 12, tpExpUnit},
			{1, 13, tpExposure},
			{0, 3, tpRGBR},
			{0, 11, tpRGBG},
			{1, 3, tpRGBB},
			{1, 13, tpExposure},
		}, // Print mode - ledRGB
	}

	// Application state
	activeMode ledMode
	lastMode   ledMode

	potUpdateChan   = make(chan potUpdate, 8)
	butIntEventChan = make(chan button.IntEvent, 8)
	butEventChan    = make(chan button.Event, 8)

	potManager = &potMgr{}
	butManager = &button.Mgr{
		IntEvents: butIntEventChan,
		Events:    butEventChan,
	}

	state = stateData{
		display1: [2][16]byte{
			stringTable[0],
			stringTable[0],
		},
		display2: [2][16]byte{
			stringTable[0],
			stringTable[0],
		},
	}

	bwM       = newBWMode(&state)
	tsM       = newTestStripMode(&state)
	exposureM = newExpMode(&state)
	focusM    = newFocusMode(&state)
)

func (s *stateData) SetLEDPanel(col [4]uint8) {
	if s.lastColour == col {
		return
	}
	for i := 0; i < ledCount; i++ {
		ledDriver.WriteByte(col[0]) // Green
		ledDriver.WriteByte(col[1]) // Red
		ledDriver.WriteByte(col[2]) // Blue
		ledDriver.WriteByte(col[3]) // White
	}
	s.lastColour = col
}

// pinToButton converts the hardware pin number to
// an internal number that is easy to work with
func pinToButton(p machine.Pin) button.Button {
	switch p {
	case machine.D7:
		return button.Plus
	case machine.D8:
		return button.Minus
	case machine.D9:
		return button.Run
	case machine.D10:
		return button.Focus
	case machine.D2:
		return button.Cancel
	case machine.D11:
		return button.Mode
	case machine.D12:
		return button.Safelight

	default:
		// should never get here
		panic("pin is not valid")
	}
}

func configureDevices() {
	machine.InitADC()

	for i := range pots {
		pots[i].Configure(machine.ADCConfig{})
	}

	err := i2c.Configure(machine.I2CConfig{})
	if err != nil {
		panic(err)
	}

	i2c.Configure(i2cConfig)
	lcd.Configure(lcdConfig)

	lcd.CreateCharacter(charPlusMinusAddr, charPlusMinus)

	ledPin.Configure(ledPinConfig)

	butInt := func(p machine.Pin) {
		ev := button.IntEvent{
			Button: pinToButton(p),
			Status: p.Get(),
		}

		select {
		case butIntEventChan <- ev:
		default:
		}
	}

	for i := range buttonPins {
		buttonPins[i].Configure(buttonPinsConfig)
		buttonPins[i].SetInterrupt(machine.PinFalling|machine.PinRising, butInt)
	}
}

func refreshTouchPoints() {
	state.activeTouchPoints = nil
	if state.activeMode.TouchPoints != nil {
		state.activeTouchPoints = state.activeMode.TouchPoints()
	}

	if len(state.activeTouchPoints) == 0 {
		potManager.SetDisabled(0, true)
	} else {
		potManager.SetDisabled(0, false)
		potManager.SetPotQuant(0, uint16(len(state.activeTouchPoints)))
	}
}

func main() {
	time.Sleep(1 * time.Second)
	configureDevices()

	// pot 0 is the touch point selector
	potManager.SetPotQuant(0, 3)

	// Not using these three pots currently
	potManager.SetDisabled(1, true)
	potManager.SetDisabled(2, true)
	potManager.SetDisabled(3, true)

	state.exposureSet.baseTime = 7_00
	for i := range state.exposureSet.exposures {
		state.exposureSet.exposures[i].grbw = [4]uint8{255, 255, 255, 255}
		if i == 0 {
			state.exposureSet.exposures[i].expUnit = expUnitAbsolute
			state.exposureSet.exposures[i].enabled = true
			continue
		}
	}

	state.exposureMode = exposureM
	state.focusMode = focusM
	state.printMode = bwM
	state.tsMode = tsM

	nextMode := state.printMode

	for {
		exitMode := false
		updateDisplay := false
		updateToucbPoints := false
		if state.prevTick == 0 {
			state.prevTick = time.Now().UnixNano()
			updateDisplay = true
		}

		now := time.Now()
		nowNS := now.UnixNano()
		if state.activeMode != nextMode {
			nextMode.SwitchTo(state.activeMode)
			state.activeMode = nextMode

			state.activeTouchPoints = nil
			if state.activeMode.TouchPoints != nil {
				state.activeTouchPoints = state.activeMode.TouchPoints()
			}

			if len(state.activeTouchPoints) == 0 {
				potManager.SetDisabled(0, true)
			} else {
				potManager.SetDisabled(0, false)
				potManager.SetPotQuant(0, uint16(len(state.activeTouchPoints)))
			}

			updateDisplay = true
		}

		// queueUp events from button and pot changes
		butManager.Process(nowNS)
		potManager.Process(nowNS)

	processEvents:
		// apply events to state
		for {
			select {
			case pu := <-potUpdateChan:
				if pu.updated > 0 {
					updateDisplay = true
					state.pots = pu.vals

					if len(state.activeTouchPoints) > 0 {
						state.activeTouchPointIndex = uint8(state.pots[0])
					}
				}
			case ev := <-butEventChan:
				var ud, em bool
				switch ev.EventType {
				case button.EventPress:
					ud, em = state.ButtonPress(ev.Button)
				case button.EventLongPress:
					ud, em = state.ButtonLongPress(ev.Button)
					if ud && ev.Button == button.Mode {
						updateToucbPoints = true
					}
				case button.EventHoldRepeat:
					ud, em = state.ButtonHoldRepeat(ev.Button)
				}
				updateDisplay = updateDisplay || ud
				exitMode = exitMode || em
			default:
				break processEvents
			}
		}

		passed := nowNS - state.prevTick
		if !(exitMode || updateDisplay) {
			if state.activeMode.Tick != nil {
				ud, em := state.activeMode.Tick(passed)
				updateDisplay = updateDisplay || ud
				exitMode = exitMode || em
			}
		}

		if updateToucbPoints {
			refreshTouchPoints()
		}

		if updateDisplay {
			state.UpdateDisplay()
		}

		if exitMode {
			nextMode = state.activeMode.SwitchAway()
		}

		// this can be a more subtle calculation
		state.prevTick = nowNS
		state.nextTick = int64(tick)

		// Rather than sleep here we may just be able to tight loop and
		// use a one of the hardware counters to track passage of ticks.
		time.Sleep(time.Duration(state.nextTick))
	}
}
