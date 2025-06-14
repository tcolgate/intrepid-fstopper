package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/hd44780i2c"
	"tinygo.org/x/drivers/ws2812"
)

const (
	tickChanLen = 0
	ledCount    = 56

	halfStop  = 141 // = 100 * (2 ^ (1 / 2))
	thirdStop = 125 // = 100 * (2 ^ (1 / 3))
)

type mode int

const (
	modeBW = iota
	modeFocus
)

type subMode int

const (
	modeBWPrint = iota
	modeBWTestStrip
)

var (
	tickChan      = make(chan struct{}, tickChanLen)
	potUpdateChan = make(chan potUpdate, 4)

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

	activeMode mode
	lastMode   mode
)

func ticker() {
	for {
		time.Sleep(10 * time.Millisecond)
		select {
		case tickChan <- (struct{}{}):
		default:
		}
	}
}

type potUpdateStatus int

const (
	conPotUpdated     = 1
	cyanPotUpdated    = 2
	magentaPotUpdated = 4
	yellowPotUpdated  = 8
)

type potUpdate struct {
	vals    [4]uint16
	updated uint8
}

type buttonState struct {
}

func potChanged(o, n uint16) bool {
	minDiff := uint16(10)

	if o > n {
		return (minDiff < (o - n))
	}

	if o < n {
		return (minDiff < (n - o))
	}

	return false
}

func potWatcher(tick <-chan struct{}) {
	conV := contrast.Get()
	cyanV := cyan.Get()
	magentaV := magenta.Get()
	yellowV := yellow.Get()

	for _ = range tick {
		var updated uint8

		if newConV := contrast.Get(); potChanged(newConV, conV) {
			updated |= conPotUpdated
			conV = newConV
		}

		if newCyanV := cyan.Get(); potChanged(newCyanV, cyanV) {
			updated |= cyanPotUpdated
			cyanV = newCyanV
		}

		if newMagentaV := magenta.Get(); potChanged(newMagentaV, magentaV) {
			updated |= magentaPotUpdated
			magentaV = newMagentaV
		}

		if newYellowV := yellow.Get(); potChanged(newYellowV, yellowV) {
			updated |= yellowPotUpdated
			yellowV = newYellowV
		}

		if updated == 0 {
			continue
		}

		update := potUpdate{
			updated: updated,
		}
		update.vals[0] = conV
		update.vals[1] = cyanV
		update.vals[2] = magentaV
		update.vals[3] = yellowV

		select {
		case potUpdateChan <- update:
		default:
		}
	}
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

	butInt := func(p machine.Pin) {
		p.Get()
		println("Button ", p, " ", p.Get())
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

	// setup

	// down here is using stuff

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

	go ticker()
	go potWatcher(tickChan)

	lcd.ClearDisplay()
	lcd.SetCursor(0, 0)
	lcd.Print(stringTable[0])

	for {
		select {
		case upd := <-potUpdateChan:
			var out numBuf
			if (upd.updated & conPotUpdated) > 0 {
				numOut(&out, num(upd.vals[0]))
				lcd.SetCursor(0, 1)
				lcd.Print(out[:])
			}
		}
	}
}
