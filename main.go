package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/hd44780i2c"
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

	// For efficiency, it's best to get the device ID once and cache it
	// (e.g. on RP2040 XIP flash and interrupts disabled for period of
	// retrieving the hardware ID from ROM chip)
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
