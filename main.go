package main

import (
	"fmt"
	"math/big"
	"os"
	"reflect"
	"strconv"

	"github.com/zellyn/openers/cmd"

	"github.com/alecthomas/kong"
)

var cli struct {
	Debug int `kong:"short='v',type='counter',help='Enable debug mode.'"`

	Secplus  cmd.SecplusCmd  `cmd:"" help:"Work with Security+2.0 devices."`
	Megacode cmd.MegaCodeCmd `cmd:"" help:"Work with MegaCode devices."`
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
		kong.NamedMapper("anybaseuint32", hexUint32Mapper{}),
		kong.NamedMapper("anybaseuint72", hexUint72Mapper{}),
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

type hexUint32Mapper struct{}

func (h hexUint32Mapper) Decode(ctx *kong.DecodeContext, target reflect.Value) error {
	t, err := ctx.Scan.PopValue("int")
	if err != nil {
		return err
	}
	var sv string
	switch v := t.Value.(type) {
	case string:
		sv = v

	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		sv = fmt.Sprintf("%v", v)

	default:
		return fmt.Errorf("expected an int but got %q (%T)", t, t.Value)
	}
	n, err := strconv.ParseUint(sv, 0, 32)
	if err != nil {
		return fmt.Errorf("expected a valid %d bit uint but got %q", 32, sv)
	}
	target.SetUint(n)
	return nil
}

type hexUint72Mapper struct{}

func (h hexUint72Mapper) Decode(ctx *kong.DecodeContext, target reflect.Value) error {
	t, err := ctx.Scan.PopValue("int")
	if err != nil {
		return err
	}
	var sv string
	switch v := t.Value.(type) {
	case string:
		sv = v
	default:
		return fmt.Errorf("expected a string value but got %q (%T)", t, t.Value)
	}
	if !target.CanInterface() {
		return fmt.Errorf("internal error: type 'anybaseuint72' only works with fields of type *big.Int; got %s", target.Type())
	}
	i := target.Interface()
	v, ok := i.(*big.Int)
	if !ok {
		return fmt.Errorf("internal error: type 'anybaseuint72' only works with fields of type *big.Int; got %s", target.Type())
	}

	err = v.UnmarshalText([]byte(sv))
	if err != nil {
		return err
	}
	if v.BitLen() > 72 {
		return fmt.Errorf("%s is > 72 bits", sv)
	}
	return nil
}
