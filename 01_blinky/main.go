package main

import (
	"fmt"
	"machine"
	"time"

	"tinygo.org/x/drivers/buzzer"
)

type note struct {
	tone     float64
	duration float64
}

func main() {
	led := machine.LED
	//led := machine.LCD_BACKLIGHT
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	cnt := 0

	bzrPin := machine.WIO_BUZZER
	bzrPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	bzr := buzzer.New(bzrPin)

	notes := []note{
		{buzzer.C3, buzzer.Quarter},
		{buzzer.D3, buzzer.Quarter},
		{buzzer.E3, buzzer.Quarter},
	}

	for {
		cnt++
		fmt.Printf("cnt %d\r\n", cnt)
		led.Low()
		time.Sleep(time.Millisecond * 500)

		led.High()
		time.Sleep(time.Millisecond * 500)

		for _, n := range notes {
			bzr.Tone(n.tone, n.duration)
			time.Sleep(10 * time.Millisecond)
		}
	}
}
