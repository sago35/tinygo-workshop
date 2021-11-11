package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/buzzer"
)

type note struct {
	tone     float64
	duration float64
}

func main() {
	bzrPin := machine.WIO_BUZZER
	bzrPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	bzr := buzzer.New(bzrPin)

	notes := []note{
		{buzzer.C3, buzzer.Quarter},
		{buzzer.Rest, buzzer.Eighth},
		{buzzer.D3, buzzer.Eighth},
		{buzzer.E3, buzzer.Quarter},
		{buzzer.Rest, buzzer.Eighth},
		{buzzer.C3, buzzer.Eighth},
		{buzzer.E3, buzzer.Quarter},
		{buzzer.C3, buzzer.Quarter},
		{buzzer.E3, buzzer.Half},
	}

	for _, n := range notes {
		bzr.Tone(n.tone, n.duration)
		time.Sleep(10 * time.Millisecond)
	}

	for {
		time.Sleep(time.Hour)
	}
}
