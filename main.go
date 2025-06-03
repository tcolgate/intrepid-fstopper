package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/hd44780i2c"
	"tinygo.org/x/drivers/ws2812"
)

const (
	tickChanLen = 4
	ledCount    = 56
)

var (
	tickChan      = make(chan struct{}, tickChanLen)
	potUpdateChan = make(chan potUpdate, tickChanLen)

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

func potWatcher(tick <-chan struct{}) {
	conV := contrast.Get()
	cyanV := cyan.Get()
	magentaV := magenta.Get()
	yellowV := yellow.Get()

	for _ = range tick {
		var updated uint8

		if newConV := contrast.Get(); newConV != conV {
			updated |= conPotUpdated
			conV = newConV
		}

		if newCyanV := cyan.Get(); newCyanV != cyanV {
			updated |= cyanPotUpdated
			cyanV = newCyanV
		}

		if newMagentaV := magenta.Get(); newMagentaV != magentaV {
			updated |= magentaPotUpdated
			magentaV = newMagentaV
		}

		if newYellowV := yellow.Get(); newYellowV != yellowV {
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

		potUpdateChan <- update
	}
}

func setLEDPanel(c color.RGBA) {
	for i := range leds {
		leds[i] = c
	}
}

func main() {
	time.Sleep(2 * time.Second)

	// setup

	machine.InitADC()

	contrast.Configure(machine.ADCConfig{})
	cyan.Configure(machine.ADCConfig{})
	magenta.Configure(machine.ADCConfig{})
	yellow.Configure(machine.ADCConfig{})

	err := i2c.Configure(machine.I2CConfig{})
	if err != nil {
		println("could not configure I2C:", err)
		return
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

	// down here is using stuff

	lcd.Print(stringTable[0])

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

	for {
		select {
		case upd := <-potUpdateChan:
			println("updated: ", upd.updated)
			println("contrast: ", upd.vals[0])
			println("cyan: ", upd.vals[1])
			println("magenta: ", upd.vals[2])
			println("yellow: ", upd.vals[3])
		}
	}
}
