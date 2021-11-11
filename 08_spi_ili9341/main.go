package main

import (
	"image/color"
	"image/png"
	"log"
	"machine"
	"strings"
	"time"

	"tinygo.org/x/drivers/ili9341"
	"tinygo.org/x/tinydraw"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

var (
	white = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	red   = color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}
	green = color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}
	blue  = color.RGBA{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF}
	black = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
)

func main() {
	display := InitDisplay()

	// body
	display.FillRectangle(20, 20, 280, 200, white)

	// lcd
	display.FillRectangle(25, 25, 270, 160, black)

	// speaker
	for i := int16(0); i < 4; i++ {
		display.FillRectangle(40+i*15, 190, 5, 20, black)
	}

	// buttons
	for i := int16(0); i < 3; i++ {
		display.FillRectangle(40+i*60, 15, 40, 5, blue)
	}

	// 5-way key
	tinydraw.FilledCircle(display, 260, 180, 20, blue)

	// text
	tinyfont.WriteLine(display, &freemono.Regular9pt7b, 30, 40, "Booting Wio Terminal...", green)

	// tinygo logo
	{
		img, err := png.Decode(strings.NewReader(tinygo_logo_s_png))
		if err != nil {
			log.Fatal(err)
		}

		w := img.Bounds().Dx()
		h := img.Bounds().Dy()
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				display.SetPixel((320-int16(w))/2+int16(x), (240-int16(h))/2+int16(y), color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: 0xFF})
			}
		}
	}

	for {
		time.Sleep(time.Hour)
	}
}

func InitDisplay() *ili9341.Device {
	machine.SPI3.Configure(machine.SPIConfig{
		SCK:       machine.LCD_SCK_PIN,
		SDO:       machine.LCD_SDO_PIN,
		SDI:       machine.LCD_SDI_PIN,
		Frequency: 48000000,
	})

	// configure backlight
	backlight := machine.LCD_BACKLIGHT
	backlight.Configure(machine.PinConfig{machine.PinOutput})

	display := ili9341.NewSPI(
		machine.SPI3,
		machine.LCD_DC,
		machine.LCD_SS_PIN,
		machine.LCD_RESET,
	)

	// configure display
	display.Configure(ili9341.Config{})

	backlight.High()

	display.SetRotation(ili9341.Rotation270)
	display.FillScreen(color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF})

	return display
}
