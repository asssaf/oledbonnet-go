package cmd

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3"

	"github.com/asssaf/oledbonnet-go/cli/cmd/drawutil"
	"github.com/asssaf/oledbonnet-go/oledbonnet"
)

type ButtonCommand struct {
	fs *flag.FlagSet
}

func NewButtonCommand() *ButtonCommand {
	c := &ButtonCommand{
		fs: flag.NewFlagSet("button", flag.ExitOnError),
	}

	c.fs.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: oledbonnet button\n", c.fs.Name())
	}

	return c
}

func (c *ButtonCommand) Name() string {
	return c.fs.Name()
}

func (c *ButtonCommand) Init(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	flag.Usage = c.fs.Usage

	if c.fs.NArg() > 0 {
		return errors.New("Too many arguments")
	}

	return nil
}

func (c *ButtonCommand) Execute() error {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	bonnet, err := oledbonnet.NewI2C(&oledbonnet.DefaultOpts)
	if err != nil {
		log.Fatal(err)
	}

	display := bonnet.GetDisplay()

	buttons := []*button{}
	buttons = append(buttons, &button{name: "Button A", pin: bonnet.GetButtonA(), pos: image.Point{80, 50}})
	buttons = append(buttons, &button{name: "Button B", pin: bonnet.GetButtonB(), pos: image.Point{110, 30}})
	buttons = append(buttons, &button{name: "Button L", pin: bonnet.GetButtonL(), pos: image.Point{9, 31}})
	buttons = append(buttons, &button{name: "Button R", pin: bonnet.GetButtonR(), pos: image.Point{51, 31}})
	buttons = append(buttons, &button{name: "Button U", pin: bonnet.GetButtonU(), pos: image.Point{30, 11}})
	buttons = append(buttons, &button{name: "Button D", pin: bonnet.GetButtonD(), pos: image.Point{30, 51}})
	buttons = append(buttons, &button{name: "Button C", pin: bonnet.GetButtonC(), pos: image.Point{30, 31}})

	img := image.NewRGBA(image.Rectangle{Max: image.Point{X: 128, Y: 64}})
	drawBackground(img)
	for _, b := range buttons {
		defer b.pin.In(gpio.PullUp, gpio.NoEdge)
		b.pin.In(gpio.PullUp, gpio.BothEdges)
		level := b.pin.Read()
		if level == gpio.Low {
			b.pressed = true
		}
		drawButton(img, b.pos, b.pressed)
	}

	if err := display.Draw(img.Bounds(), img, image.ZP); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Press buttons... \n")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	stopCh := make(chan struct{})
	go func() {
		_, ok := <-sigs
		if ok {
			fmt.Printf("Terminating...\n")
			close(stopCh)
		}
	}()

	chs := make([]chan struct{}, len(buttons))
	cases := make([]reflect.SelectCase, len(chs))
	for i, b := range buttons {
		chs[i] = make(chan struct{})
		go pinHandler(b.pin, chs[i], stopCh)
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(chs[i])}
	}

	updateButton := func(b *button) {
		level := b.pin.Read()
		pressed := false
		if level == gpio.Low {
			pressed = true
		}

		action := "released"
		if pressed {
			action = "pressed"
		}
		fmt.Printf("%s %s\n", b.name, action)

		drawButton(img, b.pos, pressed)
		if err := display.Draw(img.Bounds(), img, image.ZP); err != nil {
			log.Fatal(err)
		}

		b.pressed = pressed
	}

	for {
		chosen, _, ok := reflect.Select(cases)
		if !ok {
			break
		}

		b := buttons[chosen]
		updateButton(b)
	}

	drawBackground(img)
	if err := display.Draw(img.Bounds(), img, image.ZP); err != nil {
		log.Fatal(err)
	}

	if err := display.Halt(); err != nil {
		log.Fatal(err)
	}
	return nil
}

type button struct {
	name    string
	pressed bool
	pin     gpio.PinIn
	pos     image.Point
}

func pinHandler(pin gpio.PinIn, ch chan<- struct{}, stopCh <-chan struct{}) {
	defer close(ch)
	for {
		select {
		case <-stopCh:
			return
		default:
		}

		res := pin.WaitForEdge(time.Second)
		if !res {
			continue
		}
		ch <- struct{}{}
	}
}

func drawBackground(img *image.RGBA) {
	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)
}

func drawButton(img *image.RGBA, p image.Point, state bool) {
	white := color.RGBA{255, 255, 255, 255}
	draw.DrawMask(img, img.Bounds(), &image.Uniform{white}, image.ZP, &drawutil.Circle{p, 8}, image.ZP, draw.Over)
	if !state {
		black := color.RGBA{0, 0, 0, 255}
		draw.DrawMask(img, img.Bounds(), &image.Uniform{black}, image.ZP, &drawutil.Circle{p, 6}, image.ZP, draw.Over)
	}
}
