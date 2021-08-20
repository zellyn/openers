package main

import (
	"fmt"
	"os"

	"github.com/zellyn/openers/cmd"

	"github.com/alecthomas/kong"
)

var cli struct {
	Debug int `kong:"short='v',type='counter',help='Enable debug mode.'"`

	Secplus cmd.SecplusCmd `cmd:"" help:"Work with Security+2.0 devices."`
}

func run() error {
	ctx := kong.Parse(&cli,
		kong.Name("openers"),
		kong.Description("A commandline tool for working with gate and garage openers."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}),
	)

	globals := &cmd.Globals{
		Debug: cli.Debug,
	}
	// Call the Run() method of the selected parsed command.
	return ctx.Run(globals)
}

func main() {

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
