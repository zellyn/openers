package gpiod

import (
	realgpiod "github.com/warthog618/gpiod"
)

// IsChip checks if the named device is an accessible GPIO character device.
//
// Returns an error if not.
func IsChip(name string) error {
	return realgpiod.IsChip(name)
}

// Chips returns the names of the available GPIO devices.
func Chips() []string {
	return realgpiod.Chips()
}

// NewChip opens a GPIO character device.
func NewChip(name string, options ...realgpiod.ChipOption) (*realgpiod.Chip, error) {
	return realgpiod.NewChip(name, options...)
}

// AsOutput indicates that a line or lines be requested as an output.
//
// The initial active state for the line(s) can optionally be provided.
// If fewer values are provided than lines then the remaining lines default to
// inactive.
//
// This option overrides and clears any previous Input, RisingEdge, FallingEdge,
// BothEdges, or Debounce options.
func AsOutput(values ...int) realgpiod.OutputOption {
	return realgpiod.AsOutput(values...)
}
