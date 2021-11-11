package main

import (
	"fmt"
	"machine"
	"time"
)

func main() {
	machine.InitADC()

	led := machine.LCD_BACKLIGHT
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	sensor := machine.ADC{Pin: machine.WIO_LIGHT}
	sensor.Configure(machine.ADCConfig{})

	for {
		val := sensor.Get()
		fmt.Printf("%04X\r\n", val)
		if val < 0x8000 {
			led.Low()
		} else {
			led.High()
		}
		time.Sleep(time.Millisecond * 100)
	}
}
