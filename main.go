package main

import (
	"image/color"
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

type mode uint8

const (
	modeBW = iota
	modeFocus
)

type subMode uint8

const (
	modeBWPrint = iota
	modeBWTestStrip
)

var (
	potUpdateChan   = make(chan potUpdate, 4)
	butIntEventChan = make(chan butIntEvent, 4)
	butEventChan    = make(chan butEvent, 4)

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
	leds         [ledCount]color.RGBA
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
		[]byte("C: "),
		[]byte("M: "),
		[]byte("Y: "),
	}

	// Application state
	activeMode mode
	lastMode   mode

	potManager = &potMgr{}
	butManager = &butMgr{
		intEvents: butIntEventChan,
		events:    butEventChan,
	}
)

type potUpdateStatus uint8

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

func pinToButton(p machine.Pin) button {
	switch p {
	case machine.D7:
		return butTimePlus
	case machine.D8:
		return butTimeMinus
	case machine.D9:
		return butRun
	case machine.D10:
		return butFocus
	case machine.D2:
		return butCancel
	case machine.D11:
		return butMode
	case machine.D12:
		return butSafelight

	default:
		// should never get here
		return 0
	}
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

func setLEDPanel(c color.RGBA) {
	for i := range leds {
		leds[i] = c
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

	/*
		butInt := func(p machine.Pin) {
			ev := butIntEvent{
				button: pinToButton(p),
				status: p.Get(),
			}

			select {
			case butIntEventChan <- ev:
			default:
			}
		}
	*/
	for i := range buttonPins {
		buttonPins[i].Configure(buttonPinsConfig)
		buttonPins[i].SetInterrupt(machine.PinFalling|machine.PinRising, butManager.Int)
	}

	return nil
}

func main() {
	time.Sleep(2 * time.Second)
	println("starting")
	configureDevices()
	println("configured")

	// down here is using stuff

	/*
		for i := range leds {
			switch i {
			case 0:
				leds[i] = color.RGBA{R: 0xff, G: 0x0, B: 0x0}
			case ledCount - 1:
				leds[i] = color.RGBA{R: 0x0, G: 0x0, B: 0xff}
			default:
			}
		}

		setLEDPanel(color.RGBA{R: 0x0, G: 0x0, B: 0x00, A: 0xff})
	*/

	println("setting lcd")
	lcd.ClearDisplay()
	lcd.SetCursor(0, 0)
	lcd.Print(stringTable[0])
	println("string set")

	for {
		t := time.Now()

		butManager.process(t)
		potManager.process(t)

		select {
		case upd := <-potUpdateChan:
			var out num.NumBuf
			if (upd.updated & conPotUpdated) > 0 {
				num.Out(&out, num.Num(upd.vals[0]))
				lcd.SetCursor(0, 1)
				lcd.Print(out[:])
			}
		case upd := <-butEventChan:
			println("butEv ", upd.button, " ev ", upd.buttonEventType)
		default:
			time.Sleep(10 * time.Microsecond)
		}
	}
}
