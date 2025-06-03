package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/hd44780i2c"
	"tinygo.org/x/drivers/ws2812"
)

func main() {
	time.Sleep(2 * time.Second)
	machine.InitADC()

	//	led := machine.LED
	//	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	contrast := machine.ADC{machine.ADC0}
	contrast.Configure(machine.ADCConfig{})

	cyan := machine.ADC{machine.ADC1}
	cyan.Configure(machine.ADCConfig{})

	magenta := machine.ADC{machine.ADC2}
	magenta.Configure(machine.ADCConfig{})

	yellow := machine.ADC{machine.ADC3}
	yellow.Configure(machine.ADCConfig{})

	// Test the display
	i2c := machine.I2C0
	err := i2c.Configure(machine.I2CConfig{})
	if err != nil {
		println("could not configure I2C:", err)
		return
	}

	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: machine.TWI_FREQ_400KHZ,
	})

	lcd := hd44780i2c.New(machine.I2C0, 0x27) // some modules have address 0x3F

	lcd.Configure(hd44780i2c.Config{
		Width:       16, // required
		Height:      2,  // required
		CursorOn:    true,
		CursorBlink: true,
	})
	lcd.Print([]byte("hello"))

	// test LED panel
	p := machine.PD4
	p.Configure(machine.PinConfig{Mode: machine.PinOutput})

	butInt := func(p machine.Pin) {
		println("Button pressed", p)
	}

	bPins := []machine.Pin{
		machine.D7,  // T+
		machine.D8,  // T-
		machine.D9,  // Run
		machine.D10, // Focus
		machine.D2,  // Cancel
		machine.D11, // Mode
		machine.D12, // Safelight
	}
	for i := range bPins {
		bPins[i].Configure(machine.PinConfig{Mode: machine.PinInputPullup})
		bPins[i].SetInterrupt(machine.PinToggle, butInt)
	}

	ws := ws2812.NewSK6812(p)
	count := 56
	leds := make([]color.RGBA, count)
	for i := range leds {
		switch i {
		case 0:
			leds[i] = color.RGBA{R: 0xff, G: 0x0, B: 0x0}
		case count - 1:
			leds[i] = color.RGBA{R: 0x0, G: 0x0, B: 0xff}
		default:
			leds[i] = color.RGBA{R: 0x0, G: 0x0, B: 0x00}
		}
	}
	ws.WriteColors(leds[:])

	id := machine.Device

	for {
		time.Sleep(1 * time.Second)
		println("Device ID:", id)
		println(time.Now().Unix())

		conV := contrast.Get()
		println("contrast: ", conV)

		cyanV := cyan.Get()
		println("cyan: ", cyanV)
		magentaV := magenta.Get()
		println("magenta: ", magentaV)
		yellowV := yellow.Get()
		println("yellow: ", yellowV)
	}
}
