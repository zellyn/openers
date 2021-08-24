package cmd

// SecplusCmd is the kong `secplus` subcommand.
type MegaCodeCmd struct {
	Encode   EncodeCmd   `kong:"cmd,name='encode',help='Encode MegaCode data and display the results.'"`
	Transmit TransmitCmd `kong:"cmd,name='transmit',help='Encode MegaCode data and transmit it using a GPIO pin.'"`
}
