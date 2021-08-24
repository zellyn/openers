package cmd

import (
	"fmt"

	"github.com/zellyn/openers/bits"
	"github.com/zellyn/openers/megacode"
)

// EncodeCmd is the kong `encodev2` command.
type EncodeCmd struct {
	Identifier uint32 `kong:"required,type='anybaseuint32',placeholder='24-bit-integer',help='Opener identifier.'"`
}

// Help displays extended help and examples.
func (e EncodeCmd) Help() string {
	return `Examples:
	# Encode an opener identifier.
	openers megacode encode --identifier=0x876543`
}

// Run the `encode` command.
func (e *EncodeCmd) Run(globals *Globals) error {
	databits, err := megacode.Encode(e.Identifier)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", bits.S(databits))
	return nil
}
