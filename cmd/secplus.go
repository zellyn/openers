package cmd

// SecplusCmd is the kong `secplus` subcommand.
type SecplusCmd struct {
	EncodeV2   EncodeV2Cmd   `kong:"cmd,name='encodev2',help='Encode Security+2.0 data and display the results.'"`
	TransmitV2 TransmitV2Cmd `kong:"cmd,name='transmitv2',help='Encode Security+2.0 data and transmit it using a GPIO pin.'"`
}
