package cmd

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/zellyn/openers/bits"
	"github.com/zellyn/openers/secplus"
)

// EncodeV2Cmd is the kong `encodev2` command.
type EncodeV2Cmd struct {
	Fixed   *big.Int `kong:"required,type='anybaseuint72',placeholder='72-bit-integer',help='Fixed part of opener code.'"`
	Rolling uint32   `kong:"required,type='anybaseuint32',placeholder='28-bit-integer',help='Rolling code.'"`
}

// Help displays extended help and examples.
func (e EncodeV2Cmd) Help() string {
	return `Examples:
	# shorter v2 code (encodes to 80 bits, in two 40-bit packets)
	openers secplus encodev2 --rolling=240124710 --fixed=70678577664
	# longer v2 code (encodes to 128 bits, in two 64-bit packets
	secplus encodev2 --rolling=240129675 --fixed=4616223061045564932096`
}

// Run the `encode` command.
func (e *EncodeV2Cmd) Run(globals *Globals) error {
	fixedBytes := make([]byte, 9)
	e.Fixed.FillBytes(fixedBytes)
	fixedHigh := fixedBytes[0]
	fixedLow := binary.BigEndian.Uint64(fixedBytes[1:])

	packets, err := secplus.EncodeV2(fixedHigh, fixedLow, e.Rolling)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n%s\n", bits.S(packets[0]), bits.S(packets[1]))
	return nil
}
