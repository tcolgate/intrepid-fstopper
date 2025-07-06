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
)

var (
	ledOff   = [4]uint8{0, 0, 0, 0}
	ledWhite = [4]uint8{0, 0, 0, 255}
	ledRed   = [4]uint8{0, 255, 0, 0}

	signMinus = []byte(`-`)[0]
	signPlus  = []byte(`-`)[0]
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
	pots  [4]uint8

	baseTime           uint32 // This is the base exposure time
	exposureFactor     int8   //
	exposureFactorUnit uint8  // 0 = stops 1 = 1/2 stops 2 = 1/3 stops 3 = 1/10 stops

	remainingTime   int64 // Time remaining during running exposure
	currentExposure uint8
	exposureRunning bool // is an exposure currently running
	exposurePaused  bool // is an exposure currently running

	lastMode       mode // when returning from Focus
	lastSubMode    subMode
	currentMode    mode
	currentSubMode subMode
	lastLED        [4]uint8
	currentLED     [4]uint8

	lastDisplay [2][]byte
	nextDisplay [2][]byte
}

/*
   It would be nice to have a better structure abstracting the
	 state and transitions
*/

func (s *stateData) ButtonHoldRepeat(b button.Button) bool {
	switch b {
	case button.Plus:
		if s.exposureRunning || state.currentMode == modeFocus {
			return false
		}
		if s.baseTime != 25500 {
			s.baseTime += 10
			return true
		}
	case button.Minus:
		if s.exposureRunning || state.currentMode == modeFocus {
			return false
		}
		if s.baseTime != 0 {
			s.baseTime -= 10
			return true
		}
	}
	return false
}

func (s *stateData) ButtonPress(b button.Button) bool {
	switch b {
	case button.Plus:
		if s.exposureRunning || state.currentMode == modeFocus {
			return false
		}
		if s.baseTime != 25500 {
			s.baseTime += 10
			return true
		}
	case button.Minus:
		if s.exposureRunning || state.currentMode == modeFocus {
			return false
		}
		if s.baseTime != 0 {
			s.baseTime -= 10
			return true
		}
	case button.Run:
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
		s.remainingTime = int64(s.baseTime) * 10 * int64(time.Millisecond)
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
		if s.exposureRunning {
			return false
		}
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
	switch state.currentMode {
	case modeFocus:
		copy(s.nextDisplay[0], stringTable[1][0])
		copy(s.nextDisplay[1], stringTable[1][1])
	case modeBW:
		nb := num.NumBuf{}

		copy(s.nextDisplay[0], stringTable[0][0])
		copy(s.nextDisplay[1], stringTable[0][1])

		num.Out(&nb, num.Num(s.baseTime))
		copy(s.nextDisplay[1][0:4], nb[0:4])

		if s.exposureFactor < 0 {
			s.nextDisplay[1][4] = signMinus
		} else {
			s.nextDisplay[1][4] = signPlus
		}
		absExpFact := s.exposureFactor
		if absExpFact < 0 {
			absExpFact = absExpFact * -1
		}
		num.OutLeft(&nb, num.Num(absExpFact))
		copy(s.nextDisplay[1][5:9], nb[0:4])

		if s.exposureRunning {
			num.Out(&nb, num.Num(s.remainingTime/int64((10*time.Millisecond))))
			copy(s.nextDisplay[1][13:17], nb[0:4])
		}
	}
	for i := uint8(0); i < 2; i++ {
		if bytes.Compare(s.lastDisplay[i], s.nextDisplay[i]) != 0 {
			lcd.SetCursor(0, i)
			lcd.Print(s.nextDisplay[i])
			copy(s.lastDisplay[i][0:16], s.nextDisplay[i][0:16])
		}
	}

}

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

	contrast = machine.ADC{machine.ADC0}
	cyan     = machine.ADC{machine.ADC1}
	magenta  = machine.ADC{machine.ADC2}
	yellow   = machine.ADC{machine.ADC3}

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
			[]byte("Print           "),
			[]byte("                "),
		},
		{
			[]byte("---  Focus   ---"),
			[]byte("                "),
		},
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
		lastDisplay: [2][]byte{
			make([]byte, 16),
			make([]byte, 16),
		},
		nextDisplay: [2][]byte{
			make([]byte, 16),
			make([]byte, 16),
		},
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

	contrast.Configure(machine.ADCConfig{})
	cyan.Configure(machine.ADCConfig{})
	magenta.Configure(machine.ADCConfig{})
	yellow.Configure(machine.ADCConfig{})

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

	for {
		updated := false
		if state.prevTick == 0 {
			state.prevTick = time.Now().UnixNano()
			// force a display update on the first iteration
			updated = true
		}
		now := time.Now()
		nowNS := now.UnixNano()

		// queueUp events from button and pot changes
		butManager.Process(nowNS)
		potManager.Process(nowNS)

	processEvents:
		// apply events to state
		for {
			select {
			case _ = <-potUpdateChan:
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
		}

		if updated {
			state.SetLEDPanel()
			state.UpdateDisplay()
		}

		// this can be a more subtle calculation
		state.prevTick = nowNS
		state.nextTick = int64(10 * time.Millisecond)
		time.Until(now.Add(time.Duration(state.nextTick)))
	}
}
