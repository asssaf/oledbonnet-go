package oledbonnet

import (
	//"errors"
	//"fmt"

	"periph.io/x/conn/v3/gpio"
	//"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/ssd1306"
	"periph.io/x/host/v3/rpi"
)

var (
	LedADC1     = 0
	LedADC2     = 1
	LedADC3     = 2
	LedOutput1  = 3
	LedOutput2  = 4
	LedOutput3  = 5
	LedRelay1NO = 6
	LedRelay1NC = 7
	LedRelay2NO = 8
	LedRelay2NC = 9
	LedRelay3NO = 10
	LedRelay3NC = 11
	LedInput3   = 12
	LedInput2   = 13
	LedInput1   = 14
	LedWarn     = 15 // red
	LedComm     = 16 // blue
	LedPower    = 17 // green
)

type Opts struct {
}

var DefaultOpts = Opts{}

// Dev represents an OLED Bonnet
type Dev struct {
	opts    Opts
	buttonA gpio.PinIn
	buttonB gpio.PinIn
	buttonL gpio.PinIn
	buttonR gpio.PinIn
	buttonU gpio.PinIn
	buttonD gpio.PinIn
	buttonC gpio.PinIn
	display *ssd1306.Dev
}

// NewI2C returns a new driver.
func NewI2C(opts *Opts) (*Dev, error) {
	i2cPort, err := i2creg.Open("/dev/i2c-1")
	if err != nil {
		return nil, err
	}

	displayOpts := ssd1306.Opts{W: 128, H: 64, Rotated: false, Sequential: false, SwapTopBottom: false}

	display, err := ssd1306.NewI2C(i2cPort, &displayOpts)
	if err != nil {
		return nil, err
	}

	dev := &Dev{
		opts:    *opts,
		buttonA: rpi.P1_29, // GPIO 5
		buttonB: rpi.P1_31, // GPIO 6
		buttonL: rpi.P1_13, // GPIO 27
		buttonR: rpi.P1_16, // GPIO 23
		buttonU: rpi.P1_11, // GPIO 17
		buttonD: rpi.P1_15, // GPIO 22
		buttonC: rpi.P1_7,  // GPIO 4
		display: display,
	}

	return dev, nil
}

// GetButtonA returns gpio.PinIn corresponding to button A
func (d *Dev) GetButtonA() gpio.PinIn {
	return d.buttonA
}

// GetButtonB returns gpio.PinIn corresponding to button B
func (d *Dev) GetButtonB() gpio.PinIn {
	return d.buttonB
}

// GetButtonL returns gpio.PinIn corresponding to button L
func (d *Dev) GetButtonL() gpio.PinIn {
	return d.buttonL
}

// GetButtonR returns gpio.PinIn corresponding to button R
func (d *Dev) GetButtonR() gpio.PinIn {
	return d.buttonR
}

// GetButtonU returns gpio.PinIn corresponding to button U
func (d *Dev) GetButtonU() gpio.PinIn {
	return d.buttonU
}

// GetButtonD returns gpio.PinIn corresponding to button D
func (d *Dev) GetButtonD() gpio.PinIn {
	return d.buttonD
}

// GetButtonC returns gpio.PinIn corresponding to button C
func (d *Dev) GetButtonC() gpio.PinIn {
	return d.buttonC
}

func (d *Dev) GetDisplay() *ssd1306.Dev {
	return d.display
}

// Halt all internal devices.
func (d *Dev) Halt() error {
	if err := d.buttonA.Halt(); err != nil {
		return err
	}

	if err := d.buttonB.Halt(); err != nil {
		return err
	}

	if err := d.buttonL.Halt(); err != nil {
		return err
	}

	if err := d.buttonR.Halt(); err != nil {
		return err
	}

	if err := d.buttonU.Halt(); err != nil {
		return err
	}

	if err := d.buttonD.Halt(); err != nil {
		return err
	}

	if err := d.buttonC.Halt(); err != nil {
		return err
	}

	if d.display != nil {
		if err := d.display.Halt(); err != nil {
			return err
		}
	}

	return nil
}
