package cmd

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/zellyn/openers/gpiod"
	"github.com/zellyn/openers/megacode"
)

// TransmitCmd is the kong `encodev2` command.
type TransmitCmd struct {
	Chip       string        `kong:"default='gpiochip0',help='Chip name (device in /dev/). Must be supported by github.com/warthog618/gpiod.'"`
	Pulsewidth time.Duration `kong:"default='1ms',help='Duration of a single pulse (1/6 of a bit packet).'"`
	Repeats    int           `kong:"default='4',help='Number of times to send the whole message.'"`

	Pin        int    `kong:"required,placeholder='pin#',help='GPIO pin number.'"`
	Identifier uint32 `kong:"required,type='anybaseuint32',placeholder='24-bit-integer',help='Opener identifier.'"`
}

// Help displays extended help and examples.
func (t TransmitCmd) Help() string {
	return `Examples:
	# TODO(zellyn)
	# Encode an opener identifier.
	openers megacode encode --identifier=0x876543`
}

// Run the `encode` command.
func (t *TransmitCmd) Run(globals *Globals) error {

	if err := gpiod.IsChip(t.Chip); err != nil {
		chips := gpiod.Chips()
		if len(chips) == 0 {
			return fmt.Errorf("%q is not an available chip: there are no chips available", t.Chip)
		}
		return fmt.Errorf("%q is not an available chip; please choose one of %s", t.Chip, strings.Join(gpiod.Chips(), ","))
	}

	databits, err := megacode.Encode(t.Identifier)
	if err != nil {
		return err
	}
	chip, err := gpiod.NewChip(t.Chip)
	if err != nil {
		return err
	}
	defer chip.Close()

	line, err := chip.RequestLine(t.Pin, gpiod.AsOutput(0))
	if err != nil {
		return err
	}
	defer line.SetValue(0)
	defer line.Close()

	target := time.Now()
	for i := 0; i < t.Repeats; i++ {
		if i > 0 {
			target.Add(t.Pulsewidth * 6)
			for target.After(time.Now()) {
			}
		}

		fmt.Printf("Sending repetition %d/%d\n", i+1, t.Repeats)

		runtime.GC()
		debug.SetGCPercent(-1)

		for _, bit := range databits {
			target = target.Add(t.Pulsewidth)
			if err := line.SetValue(int(bit)); err != nil {
				return err
			}
			for target.After(time.Now()) {
			}
		}
		if err := line.SetValue(0); err != nil {
			return err
		}
	}

	return nil
}
