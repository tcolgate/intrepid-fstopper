package main

import (
	"bytes"
	"intrepidfstopper/button"
	"machine"
	"time"

	"tinygo.org/x/drivers/hd44780i2c"
	"tinygo.org/x/drivers/ws2812"
)

const (
	ledCount = 56

	halfStop  = 141 // = 100 * (2 ^ (1 / 2))
	thirdStop = 125 // = 100 * (2 ^ (1 / 3))

	longPress = 1 * time.Second

	tick = int64(10 * time.Millisecond)
)

var (
	ledOff   = [4]uint8{0, 0, 0, 0}
	ledWhite = [4]uint8{0, 0, 0, 255}
	ledRed   = [4]uint8{0, 255, 0, 0}

	signMinus = []byte(`-`)[0]
	signPlus  = []byte(`+`)[0]
)

type mode uint8

const (
	modeBW = iota
	modeRGB
	modeFocus
)

type subMode uint8

const (
	modePrint subMode = iota
	modeTestStrip
)

type expUnit uint8

const (
	expUnitHalfStop expUnit = iota
	expUnitThirdStop
	expUnitTenthStop
	expUnitPercent
	expUnitFreeHand
)

var (
	expUnitNames = [5][]byte{
		[]byte("\xDF/2 "),
		[]byte("\xDF/3 "),
		[]byte("\xDF/10"),
		[]byte(`%   `),
		[]byte(`Free`),
	}
)

type stateBits int

const (
	statebitFocusColour stateBits = 1 << iota
)

func setStateBit(s *stateBits, b stateBits) {
	*s = *s | b
}

func clearStateBit(s *stateBits, b stateBits) {
	*s = *s & (^b)
}

func checkStateBit(s *stateBits, b stateBits) bool {
	return (*s & b) > 0
}

func toggleStateBit(s *stateBits, b stateBits) bool {
	*s = *s ^ b
	return (*s & b) > 0
}

type stateData struct {
	nextTick int64
	prevTick int64

	flags stateBits
	pots  [4]uint16

	lastSubMode    subMode
	currentSubMode subMode

	activeTouchPoints     []touchPoint
	activeTouchPointIndex uint8

	activeDisplay bool
	display1      [2][]byte
	display2      [2][]byte

	exposureMode *Mode
	focusMode    *Mode
	bwMode       *Mode

	activeMode *Mode

	lastColour  [4]uint8
	focusColour [4]uint8
}

/*
   It would be nice to have a better structure abstracting the
	 state and transitions
*/

func (s *stateData) ButtonHoldRepeat(b button.Button) (bool, bool) {
	switch b {
	case button.Plus:
		if state.activeMode.PressLongPlus != nil {
			return s.activeMode.PressLongPlus(s.activeTouchPointIndex)
		}
	case button.Minus:
		if state.activeMode.PressLongMinus != nil {
			return s.activeMode.PressLongMinus(s.activeTouchPointIndex)
		}
	}
	return false, false
}

func (s *stateData) ButtonPress(b button.Button) (bool, bool) {
	switch b {
	case button.Plus:
		if state.activeMode.PressPlus != nil {
			return s.activeMode.PressPlus(s.activeTouchPointIndex)
		}
	case button.Minus:
		if state.activeMode.PressMinus != nil {
			return s.activeMode.PressMinus(s.activeTouchPointIndex)
		}
	case button.Run:
		if state.activeMode.PressRun != nil {
			return s.activeMode.PressRun()
		}

		/*
			if s.exposureRunning {
				s.exposurePaused = !s.exposurePaused
				if s.exposurePaused {
					state.currentLED = ledOff
				} else {
					state.currentLED = ledWhite
				}
				return true
			}

			// start exposure
			s.remainingTime = int64(s.baseTime) * tick
			s.exposureRunning = true
			s.exposurePaused = false
			state.currentLED = ledWhite
			return true
		*/
	case button.Cancel:
		if state.activeMode.PressCancel != nil {
			return s.activeMode.PressCancel(s.activeTouchPointIndex)
		}

		/*
			if s.exposureRunning {
				s.exposurePaused = false
				s.exposureRunning = false
				s.remainingTime = 0
				state.currentLED = ledOff
				return true
				// stop exposure, reset time
			}
			if state.currentMode == modeFocus {
				state.currentMode = state.lastMode
				state.currentSubMode = state.lastSubMode
				clearStateBit(&state.flags, statebitFocusColour)
				state.currentLED = ledOff
				return true
			}
		*/
	case button.Focus:
		if state.activeMode.PressFocus != nil {
			return s.activeMode.PressFocus()
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
	}
	return false, false
}

func (s *stateData) UpdateDisplay() {
	println("in update display")
	lastDisplay := &s.display1
	nextDisplay := &s.display2
	if s.activeDisplay {
		lastDisplay = &s.display2
		nextDisplay = &s.display1
	}

	s.activeMode.UpdateDisplay(nextDisplay)

	for i := uint8(0); i < 2; i++ {
		if bytes.Compare(lastDisplay[i], nextDisplay[i]) != 0 {
			lcd.SetCursor(0, i)
			lcd.Print(nextDisplay[i])
			copy(lastDisplay[i][0:16], nextDisplay[i][0:16])
		}
	}

	if len(s.activeTouchPoints) > 0 {
		tp := s.activeTouchPoints[s.activeTouchPointIndex]
		lcd.SetCursor(tp[1], tp[0])
		lcd.CursorOn(true)
	} else {
		lcd.CursorOn(false)
	}
}

type touchPoint [2]uint8

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

	stringTable = [][2][]byte{
		{
			[]byte("B: Adj  U:     >"),
			[]byte("*              )"),
		},
		{
			[]byte("---  Focus   ---"),
			[]byte("                "),
		},
		{
			[]byte("--- Exposure ---"),
			[]byte("                "),
		},
	}
	touchPoints = [][]touchPoint{
		[]touchPoint{{0, 0}, {1, 0}, {0, 8}},
		nil,
	}

	// Application state
	activeMode mode
	lastMode   mode

	potUpdateChan   = make(chan potUpdate, 8)
	butIntEventChan = make(chan button.IntEvent, 8)
	butEventChan    = make(chan button.Event, 8)

	potManager = &potMgr{}
	butManager = &button.Mgr{
		IntEvents: butIntEventChan,
		Events:    butEventChan,
	}

	state = stateData{
		display1: [2][]byte{
			make([]byte, 16),
			make([]byte, 16),
		},
		display2: [2][]byte{
			make([]byte, 16),
			make([]byte, 16),
		},
	}

	bwM       = newBWMode(&state)
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

func main() {
	time.Sleep(2 * time.Second)
	configureDevices()

	potManager.SetPotQuant(0, 3)

	potManager.SetDisabled(1, true)
	potManager.SetDisabled(2, true)
	potManager.SetDisabled(3, true)

	state.exposureMode = exposureM
	state.focusMode = focusM
	state.bwMode = bwM

	state.activeMode = state.bwMode
	nextMode := state.bwMode

	for {
		exitMode := false
		updateDisplay := false
		if state.prevTick == 0 {
			state.prevTick = time.Now().UnixNano()
			updateDisplay = true
		}

		now := time.Now()
		nowNS := now.UnixNano()
		if state.activeMode != nextMode {
			nextMode.SwitchTo(state.activeMode)
			state.activeMode = nextMode
			state.activeTouchPoints = state.activeMode.TouchPoints()

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
				}
			case ev := <-butEventChan:
				var ud, em bool
				switch ev.EventType {
				case button.EventPress:
					ud, em = state.ButtonPress(ev.Button)
				case button.EventLongPress:
					ud, em = state.ButtonLongPress(ev.Button)
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

		if updateDisplay {
			println("calling update display", nowNS)
			state.UpdateDisplay()
		}

		if exitMode {
			println("calling switchAway", nowNS)
			nextMode = state.activeMode.SwitchAway()
		}

		// this can be a more subtle calculation
		state.prevTick = nowNS
		state.nextTick = tick
		time.Until(now.Add(time.Duration(state.nextTick)))
	}
}
