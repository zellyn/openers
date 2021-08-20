package cmd

import "fmt"

// EncodeV2Cmd is the kong `encodev2` command.
type EncodeV2Cmd struct {
	System string `kong:"default='auto',enum='auto,dos3',help='DOS system used for image.'"`

	Filename string `kong:"arg,required,help='Filename to use on disk.'"`
}

// Help displays extended help and examples.
func (d EncodeV2Cmd) Help() string {
	return `Examples:
	# TODO(zellyn)
	echo no example yet`
}

// Run the `encode` command.
func (d *EncodeV2Cmd) Run(globals *Globals) error {
	return fmt.Errorf("not implemented yet")
}
