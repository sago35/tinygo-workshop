package main

import (
	"device/arm"
	"machine"
)

var timerCh = make(chan struct{}, 1)

func main() {
	//led := machine.LED
	led := machine.LCD_BACKLIGHT

	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// timer fires 10 times per second
	arm.SetupSystemTimer(machine.CPUFrequency() / 10)

	for {
		led.Low()
		<-timerCh
		led.High()
		<-timerCh
	}
}

//export SysTick_Handler
func timer_isr() {
	select {
	case timerCh <- struct{}{}:
	default:
		// The consumer is running behind.
	}
}
