package main

import (
	"bytes"
	"intrepidfstopper/button"
	"intrepidfstopper/num"
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

	baseTime           uint32  // This is the base exposure time
	exposureFactor     int8    //
	exposureFactorUnit expUnit // 0 = stops 1 = 1/2 stops 2 = 1/3 stops 3 = 1/10 stops

	remainingTime   int64 // Time remaining during running exposure
	currentExposure uint8
	exposureRunning bool // is an exposure currently running
	exposurePaused  bool // is an exposure currently running

	lastSubMode    subMode
	currentMode    mode
	currentSubMode subMode
	lastLED        [4]uint8
	currentLED     [4]uint8

	activeTouchPoints     []touchPoint
	activeTouchPointIndex uint8
	lastDisplay           [2][]byte
	nextDisplay           [2][]byte

	exposureMode Mode
	focusMode    Mode
	bwMode       Mode

	lastMode   Mode
	activeMode Mode
	nextMode   Mode
}

/*
   It would be nice to have a better structure abstracting the
	 state and transitions
*/

func (s *stateData) ButtonHoldRepeat(b button.Button) bool {
	switch b {
	case button.Plus:
		return s.activeMode.LongPlus(s.activeTouchPointIndex)
	case button.Minus:
		return s.activeMode.LongMinus(s.activeTouchPointIndex)
	default:
		return false
	}
}

func (s *stateData) ButtonPress(b button.Button) bool {
	switch b {
	case button.Plus:
		s.activeMode.Plus(s.activeTouchPointIndex)
	case button.Minus:
		s.activeMode.Minus(s.activeTouchPointIndex)
	case button.Run:
		s.activeMode.Minus(s.activeTouchPointIndex)
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
	case button.Cancel:
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
	case button.Focus:
		if s.exposureRunning {
			return false
		}
		if state.currentMode != modeFocus {
			state.lastMode = state.currentMode
			state.lastSubMode = state.currentSubMode
			state.currentMode = modeFocus
			clearStateBit(&state.flags, statebitFocusColour)
			state.currentLED = ledRed
		} else {
			state.currentMode = state.lastMode
			state.currentSubMode = state.lastSubMode
			clearStateBit(&state.flags, statebitFocusColour)
			state.currentLED = ledOff
		}
		return true
	}
	return false
}

func (s *stateData) ButtonLongPress(b button.Button) bool {
	switch b {
	case button.Focus:
		if s.currentMode == modeFocus {
			if toggleStateBit(&state.flags, statebitFocusColour) {
				s.currentLED = ledWhite
			} else {
				s.currentLED = ledRed
			}
		} else {
			state.lastMode = state.currentMode
			state.lastSubMode = state.currentSubMode
			state.currentMode = modeFocus
			setStateBit(&state.flags, statebitFocusColour)
			state.currentLED = ledWhite
		}
		return true
	}
	return false
}

func (s *stateData) UpdateDisplay() {

	nb := num.NumBuf{}
	hasTouchPoints := false
	var tp touchPoint

	s.activeMode.UpdateDisplay(&s.lastDisplay)

	switch {
	case state.exposureRunning:
		copy(s.nextDisplay[0], stringTable[2][0])
		copy(s.nextDisplay[1], stringTable[2][1])

		if s.exposureRunning {
			num.Out(&nb, num.Num(s.remainingTime/int64((10*time.Millisecond))))
			copy(s.nextDisplay[1][12:16], nb[0:4])
		}
	case state.currentMode == modeBW:
		hasTouchPoints = true
		tpi := s.pots[0]
		s.activeTouchPoint = uint8(tpi) // WRONG, shouldn't be updating state in here
		tp = s.activeTouchPoints[tpi]

		copy(s.nextDisplay[0], stringTable[0][0])
		copy(s.nextDisplay[1], stringTable[0][1])

		num.Out(&nb, num.Num(s.baseTime))
		copy(s.nextDisplay[0][2:6], nb[0:4])

		if s.exposureFactor < 0 {
			s.nextDisplay[1][1] = signMinus
		} else {
			s.nextDisplay[1][1] = signPlus
		}
		absExpFact := s.exposureFactor
		if absExpFact < 0 {
			absExpFact = absExpFact * -1
		}
		num.OutLeft(&nb, num.Num(absExpFact))
		copy(s.nextDisplay[1][2:6], nb[0:4])

		copy(s.nextDisplay[0][10:15], expUnitNames[s.exposureFactorUnit][0:4])

		res := num.Num(11_23)
		num.Out(&nb, num.Num(res))
		resLen := num.Len(res)
		s.nextDisplay[1][13-resLen] = []byte("(")[0]
		s.nextDisplay[1][14-resLen] = []byte("=")[0]
		copy(s.nextDisplay[1][11:15], nb[0:4])
		s.nextDisplay[1][15] = []byte(")")[0]
	}

	for i := uint8(0); i < 2; i++ {
		if bytes.Compare(s.lastDisplay[i], s.nextDisplay[i]) != 0 {
			lcd.SetCursor(0, i)
			lcd.Print(s.nextDisplay[i])
			copy(s.lastDisplay[i][0:16], s.nextDisplay[i][0:16])
		}
	}

	if len(s.activeTouchPoints) > 0 {
		lcd.SetCursor(tp[1], tp[0])
	}

	lcd.CursorOn(hasTouchPoints)
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

	bwM       = &exposureMode{}
	exposureM = &exposureMode{}
	focusM    = &focusMode{}

	state = stateData{
		baseTime:           7_00,
		exposureFactorUnit: 1, // default to 1/2 stops

		lastDisplay: [2][]byte{
			make([]byte, 16),
			make([]byte, 16),
		},
		nextDisplay: [2][]byte{
			make([]byte, 16),
			make([]byte, 16),
		},

		exposureMode: exposureM,
		bwMode:       bwM,
		focusMode:    focusM,

		activeMode: exposureM,
	}
)

func (s *stateData) SetLEDPanel() {
	if s.lastLED == s.currentLED {
		return
	}
	for i := 0; i < ledCount; i++ {
		ledDriver.WriteByte(s.currentLED[0]) // Green
		ledDriver.WriteByte(s.currentLED[1]) // Red
		ledDriver.WriteByte(s.currentLED[2]) // Blue
		ledDriver.WriteByte(s.currentLED[3]) // White
	}
	s.lastLED = s.currentLED

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

func configureDevices() error {
	machine.InitADC()

	for i := range pots {
		pots[i].Configure(machine.ADCConfig{})
	}

	err := i2c.Configure(machine.I2CConfig{})
	if err != nil {
		return err
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

	return nil
}

func main() {
	time.Sleep(2 * time.Second)
	configureDevices()

	potManager.SetPotQuant(0, 3)

	potManager.SetPotDisabled(1, true)
	potManager.SetPotDisabled(2, true)
	potManager.SetPotDisabled(3, true)
	state.nextMode = state.exposureMode

	for {
		updated := false
		if state.prevTick == 0 {
			state.prevTick = time.Now().UnixNano()
			updated = true
		}

		now := time.Now()
		nowNS := now.UnixNano()

		if state.activeMode != state.nextMode {
			if state.activeMode != nil {
				state.nextMode = state.activeMode.SwitchAway()
			}
			state.nextMode.SwitchTo(state.activeMode)
			state.activeMode = state.nextMode
			state.activeTouchPoints = state.activeMode.TouchPoints()

			if len(state.activeTouchPoints) == 0 {
				potManager.SetDisabled(0, true)
			} else {
				potManager.SetDisabled(0, false)
				potManager.SetPotQuant(0, uint16(len(state.activeTouchPoints)))
			}
		}

		// queueUp events from button and pot changes
		butManager.Process(nowNS)
		potManager.Process(nowNS)

		var nextMode Mode
	processEvents:
		// apply events to state
		for {
			select {
			case pu := <-potUpdateChan:
				if pu.updated > 0 {
					updated = true
					state.pots = pu.vals
				}
			case ev := <-butEventChan:
				switch ev.EventType {
				case button.EventPress:
					updated = updated || state.ButtonPress(ev.Button)
				case button.EventLongPress:
					updated = updated || state.ButtonLongPress(ev.Button)
				case button.EventHoldRepeat:
					updated = updated || state.ButtonHoldRepeat(ev.Button)
				}
			default:
				break processEvents
			}
		}

		if state.exposureRunning && !state.exposurePaused {
			passed := nowNS - state.prevTick
			state.remainingTime -= passed
			if state.remainingTime <= 0 {
				state.exposurePaused = false
				state.exposureRunning = false
				state.remainingTime = 0
				state.currentLED = ledOff
			}
			updated = true
			updateDisplay, _ := state.activeMode.Tick(passed)
			updated = updated || updateDisplay
		}

		var exitMode bool
		if updated {
			state.SetLEDPanel()
			state.UpdateDisplay()
		}
		if exitMode {
			state.lastMode = state.activeMode
			nextMode = state.activeMode.SwitchAway()
			state.activeMode = nextMode
			nextMode.SwitchTo(state.lastMode)
		}

		// this can be a more subtle calculation
		state.prevTick = nowNS
		state.nextTick = tick
		time.Until(now.Add(time.Duration(state.nextTick)))
	}
}
