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

	halfStop  = 141 // = 100 * (2 ^ (1 / 2))
	thirdStop = 125 // = 100 * (2 ^ (1 / 3))

	longPress = 1 * time.Second
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

func toggleStateBit(s *stateBits, b stateBits) {
	*s = *s ^ b
}

type stateData struct {
	nextTick uint32

	flags stateBits
	pots  [4]uint8

	baseTime        uint32 // This is the base exposure time
	remaingingTime  uint32 // Time remaining during running exposure
	exposureRunning bool   // is an exposure currently running

	lastMode       mode // when returning from Focus
	lastSubMode    subMode
	currentMode    mode
	currentSubMode subMode
}

/*
   It would be nice to have a better structure abstracting the
	 state and transitions
*/

func (s *stateData) ButtonPress(b button.Button) bool {
	switch b {
	case button.Run:
		if s.exposureRunning {
			// pause exposure
		}
	case button.Cancel:
		if s.exposureRunning {
			return false
			// stop exposure, reset time
		}
		if state.currentMode == modeFocus {
			state.currentMode = state.lastMode
			state.currentSubMode = state.lastSubMode
			clearStateBit(&state.flags, statebitFocusColour)
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
		} else {
			state.currentMode = state.lastMode
			state.currentSubMode = state.lastSubMode
			clearStateBit(&state.flags, statebitFocusColour)
		}
		return true
	default:
		return false
	}
	return false
}

func (s *stateData) ButtonLongPress(b button.Button) bool {
	switch b {
	case button.Focus:
		if s.currentMode == modeFocus {
			toggleStateBit(&state.flags, statebitFocusColour)
			return true
		}
	default:
		return false
	}
	return false
}

func (s *stateData) ButtonHoldRepeat(b button.Button) bool {
	return false
}

func (s *stateData) UpdateDisplay() {
	lcd.SetCursor(0, 0)
	switch state.currentMode {
	case modeFocus:
		lcd.Print(stringTable[1])
		if !checkStateBit(&state.flags, statebitFocusColour) {
			setLEDPanel([4]uint8{0, 255, 0, 0})
		} else {
			setLEDPanel([4]uint8{0, 0, 0, 255})
		}
	case modeBW:
		lcd.Print(stringTable[0])
		setLEDPanel([4]uint8{0, 0, 0, 0})
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

	stringTable = [][]byte{
		[]byte("Hello: "),
		[]byte("---  Focus  ---"),
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

	state = stateData{}
)

const (
	conPotUpdated = 1 << iota
	cyanPotUpdated
	magentaPotUpdated
	yellowPotUpdated
)

type potUpdate struct {
	vals    [4]uint16
	updated uint8
}

func potChanged(o, n uint16) bool {
	// This gives us 128 valid pot positions
	minDiff := uint16(128)

	if o > n {
		return (minDiff < (o - n))
	}

	if o < n {
		return (minDiff < (n - o))
	}

	return false
}

func setLEDPanel(c [4]uint8) {
	for i := 0; i < ledCount; i++ {
		ledDriver.WriteByte(c[0]) // Green
		ledDriver.WriteByte(c[1]) // Red
		ledDriver.WriteByte(c[2]) // Blue
		ledDriver.WriteByte(c[3]) // White
	}
}

// pinToButton converts the hardware pin number to
// an internal number that is easy to work with
func pinToButton(p machine.Pin) button.Button {
	switch p {
	case machine.D7:
		return button.TimePlus
	case machine.D8:
		return button.TimeMinus
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

		now := time.Now()

		// queueUp events from button and pot changes
		butManager.Process(now.UnixNano())
		potManager.Process(now.UnixNano())

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

		if updated {
			state.UpdateDisplay()
		}

		// this can be a more subtle calculation
		state.nextTick = uint32(20 * time.Millisecond)
		time.Until(now.Add(time.Duration(state.nextTick)))
	}
}
