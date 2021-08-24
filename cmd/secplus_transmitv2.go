package cmd

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/zellyn/openers/bits"
	"github.com/zellyn/openers/gpiod"
	"github.com/zellyn/openers/secplus"
)

// TransmitV2Cmd is the kong `encodev2` command.
type TransmitV2Cmd struct {
	Chip       string        `kong:"default='gpiochip0',help='Chip name (device in /dev/). Must be supported by github.com/warthog618/gpiod.'"`
	Pulsewidth time.Duration `kong:"default='250µs',help='Duration of a single on/off pulse (half a Manchester-coded bit).'"`
	Burstgap   time.Duration `kong:"default='90ms',help='Gap between first and second burst in a message.'"`
	Repeatgap  time.Duration `kong:"default='90ms',help='Gap between repeats of the whole message.'"`
	Repeats    int           `kong:"default='4',help='Number of times to send the whole message.'"`

	Pin     int      `kong:"required,placeholder='pin#',help='GPIO pin number.'"`
	Fixed   *big.Int `kong:"required,type='anybaseuint72',placeholder='72-bit-integer',help='Fixed part of opener code.'"`
	Rolling uint32   `kong:"required,type='anybaseuint32',placeholder='28-bit-integer',help='Rolling code.'"`
}

// Help displays extended help and examples.
func (t TransmitV2Cmd) Help() string {
	return `Examples:
	# TODO(zellyn)
	# shorter v2 code (encodes to 80 bits, in two 40-bit packets)
	openers secplus encodev2 --rolling=240124710 --fixed=70678577664
	# longer v2 code (encodes to 128 bits, in two 64-bit packets
	openers secplus encodev2 --rolling=240129675 --fixed=4616223061045564932096`
}

// Run the `encode` command.
func (t *TransmitV2Cmd) Run(globals *Globals) error {

	if err := gpiod.IsChip(t.Chip); err != nil {
		chips := gpiod.Chips()
		if len(chips) == 0 {
			return fmt.Errorf("%q is not an available chip: there are no chips available", t.Chip)
		}
		return fmt.Errorf("%q is not an available chip; please choose one of %s", t.Chip, strings.Join(gpiod.Chips(), ","))
	}

	fixedBytes := make([]byte, 9)
	t.Fixed.FillBytes(fixedBytes)
	fixedHigh := fixedBytes[0]
	fixedLow := binary.BigEndian.Uint64(fixedBytes[1:])

	bursts, err := secplus.EncodeV2ToBursts(fixedHigh, fixedLow, t.Rolling)
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

	for i := 0; i < t.Repeats; i++ {
		if i > 0 {
			time.Sleep(t.Repeatgap)
		}

		fmt.Printf("Sending repetition %d/%d\n", i+1, t.Repeats)

		for j := 0; j < 2; j++ {
			if j > 0 {
				target := time.Now().Add(t.Burstgap)
				for target.After(time.Now()) {
				}
			}

			runtime.GC()
			debug.SetGCPercent(-1)

			fmt.Printf("    sending: %s\n", bits.S(bursts[j]))
			target := time.Now()
			for _, bit := range bursts[j] {
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
	}

	return nil
}
