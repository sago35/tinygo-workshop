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

	initEyes()

	machine.WIO_5S_UP.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	machine.WIO_5S_LEFT.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	machine.WIO_5S_RIGHT.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	machine.WIO_5S_DOWN.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	redraw := true
	xofs := int16(0)
	yofs := int16(0)
	eye := eyeClose
	eyeCh := make(chan struct{}, 1)
	go func() {
		for {
			eyeCh <- struct{}{}
			time.Sleep(1500 * time.Millisecond)
			eyeCh <- struct{}{}
			time.Sleep(300 * time.Millisecond)
		}
	}()

	for {
		if !machine.WIO_5S_UP.Get() {
			if 0 < yofs {
				yofs--
				redraw = true
			}
		} else if !machine.WIO_5S_LEFT.Get() {
			if -4 < xofs {
				xofs--
				redraw = true
			}
		} else if !machine.WIO_5S_RIGHT.Get() {
			if xofs < 7 {
				xofs++
				redraw = true
			}
		} else if !machine.WIO_5S_DOWN.Get() {
			if yofs < 20 {
				yofs++
				redraw = true
			}
		}

		select {
		case <-eyeCh:
			if eye == eyeOpen {
				eye = eyeClose
			} else {
				eye = eyeOpen
			}
			redraw = true
		default:
		}

		if redraw {
			drawEye(display, 127+xofs, 91+yofs, eye)
			drawEye(display, 181+xofs, 91+yofs, eye)
			redraw = false
		}
		time.Sleep(50 * time.Millisecond)
	}
}

const (
	eyeOpen = iota
	eyeClose
	eyeClear
)

func drawEye(display *ili9341.Device, x, y int16, mode int) {
	switch mode {
	case eyeOpen:
		display.DrawRGBBitmap(x, y, eyeOpenRGB[:], 12, 12)
	case eyeClose:
		display.DrawRGBBitmap(x, y, eyeCloseRGB[:], 12, 12)
	case eyeClear:
		display.FillRectangle(x, y, 12, 12, white)
	}
}

var (
	eyeOpenRGB  [144]uint16
	eyeCloseRGB [144]uint16
)

func initEyes() {
	for i := 0; i < len(eye_open_rgbbitmap)/2; i++ {
		eyeOpenRGB[i] = uint16(eye_open_rgbbitmap[i*2])<<8 + uint16(eye_open_rgbbitmap[i*2+1])
	}

	for i := 0; i < len(eye_close_rgbbitmap)/2; i++ {
		eyeCloseRGB[i] = uint16(eye_close_rgbbitmap[i*2])<<8 + uint16(eye_close_rgbbitmap[i*2+1])
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
