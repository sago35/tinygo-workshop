package main

import (
	"machine"
)

func main() {
	usbcdc := machine.Serial

	for {
		if usbcdc.Buffered() > 0 {
			data, _ := usbcdc.ReadByte()

			// just echo the character
			usbcdc.WriteByte(data)
		}
	}
}
