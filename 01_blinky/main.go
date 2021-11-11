package main

import (
	"fmt"
	"machine"
	"time"
)

func main() {
	led := machine.LED
	//led := machine.LCD_BACKLIGHT
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	cnt := 0
	for {
		cnt++
		fmt.Printf("cnt %d\r\n", cnt)
		led.Low()
		time.Sleep(time.Millisecond * 500)

		led.High()
		time.Sleep(time.Millisecond * 500)
	}
}
