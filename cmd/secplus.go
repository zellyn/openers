package cmd

// Globals holds variables global to all commands.
type Globals struct {
	Debug int // Debugging level (number of -v's)
}

// SecplusCmd is the kong `secplus` subcommand.
type SecplusCmd struct {
	EncodeV2 EncodeV2Cmd `kong:"cmd,name='encodev2',help='Encode Security+2.0 data and display the results.'"`
}
