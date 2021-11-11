package main

import (
	"machine"
	"time"
)

func main() {
	led := machine.LCD_BACKLIGHT
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.High()

	usbcdc := machine.Serial
	usbcdc.Configure(machine.UARTConfig{})

	input := make([]byte, 64)
	i := 0
	for {
		if usbcdc.Buffered() > 0 {
			data, _ := usbcdc.ReadByte()

			switch data {
			case 13:
				// return key
				usbcdc.Write([]byte("\r\n"))

				switch string(input[:i]) {
				case "on":
					led.High()
				case "off":
					led.Low()
				case "toggle", "t":
					led.Toggle()
				}
				i = 0
			default:
				// just echo the character
				usbcdc.WriteByte(data)
				input[i] = data
				i++
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}
