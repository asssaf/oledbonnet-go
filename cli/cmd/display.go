package cmd

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"periph.io/x/host/v3"

	"github.com/asssaf/oledbonnet-go/oledbonnet"
)

type DisplayCommand struct {
	fs     *flag.FlagSet
	action bool
}

func NewDisplayCommand() *DisplayCommand {
	c := &DisplayCommand{
		fs: flag.NewFlagSet("display", flag.ExitOnError),
	}

	c.fs.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: oledbonnet %s <on|off>\n", c.fs.Name())
	}

	return c
}

func (c *DisplayCommand) Name() string {
	return c.fs.Name()
}

func (c *DisplayCommand) Init(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	flag.Usage = c.fs.Usage

	if c.fs.NArg() < 1 {
		return errors.New("Missing required arguments")
	}

	if c.fs.NArg() > 1 {
		return errors.New("Too many arguments")
	}

	action := c.fs.Arg(0)

	switch action {
	case "on":
		c.action = true
	case "off":
		c.action = false
	default:
		return errors.New("unrecognized action, must be one of: [on, off]")
	}

	return nil
}

func (c *DisplayCommand) Execute() error {
	actionName := "off"
	if c.action == true {
		actionName = "on"
	}
	fmt.Printf("Turning display %s\n", actionName)

	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	bonnet, err := oledbonnet.NewI2C(&oledbonnet.DefaultOpts)
	if err != nil {
		log.Fatal(err)
	}

	display := bonnet.GetDisplay()

	if c.action {
		// any command turns the display on
		if err := display.StopScroll(); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := display.Halt(); err != nil {
			log.Fatal(err)
		}
	}

	return nil
}
