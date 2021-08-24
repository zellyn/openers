//+build !linux

//go:build !linux
package gpiod

import "fmt"

// IsChip checks if the named device is an accessible GPIO character device.
//
// Returns an error if not.
func IsChip(name string) error {
	return fmt.Errorf("this is just a fake; IsChip is not supported except on Linux")
}

// Chips returns the names of the available GPIO devices.
func Chips() []string {
	return nil
}

// NewChip opens a GPIO character device.
func NewChip(name string, options ...ChipOption) (*Chip, error) {
	return nil, fmt.Errorf("this is just a fake; NewChip is not supported except on Linux")
}

// AsOutput indicates that a line or lines be requested as an output.
//
// The initial active state for the line(s) can optionally be provided.
// If fewer values are provided than lines then the remaining lines default to
// inactive.
//
// This option overrides and clears any previous Input, RisingEdge, FallingEdge,
// BothEdges, or Debounce options.
func AsOutput(values ...int) OutputOption {
	return OutputOption{}
}

// OutputOption indicates the line direction should be set to an output.
type OutputOption []int

// ChipOption defines the interface required to provide a Chip option.
type ChipOption interface{}

// Chip represents a single GPIO chip that controls a set of lines.
type Chip struct{}

// Close releases the Chip.
//
// It does not release any lines which may be requested - they must be closed
// independently.
func (c *Chip) Close() error {
	return nil
}

// RequestLine requests control of a single line on the chip.
//
// If granted, control is maintained until either the Line or Chip are closed.
func (c *Chip) RequestLine(offset int, options ...LineReqOption) (*Line, error) {
	return nil, fmt.Errorf("this is just a fake; RequestLine is not supported except on Linux")
}

// Line represents a single requested line.
type Line struct{}

type LineReqOption interface {
}

// SetValue sets the current active state of the line.
//
// Only valid for output lines.
func (l *Line) SetValue(value int) error {
	return fmt.Errorf("this is just a fake; SetValue is not supported except on Linux")
}

// Close releases all resources held by the requested line.
//
// Note that this includes waiting for any running event handler to return.
// As a consequence the Close must not be called from the context of the event
// handler - the Close should be called from a different goroutine.
func (l *Line) Close() error {
	return nil
}
