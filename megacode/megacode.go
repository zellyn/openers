package megacode

import "fmt"

// Encode encodes a MegaCode identifier to a bitstream of simple 0s and 1s.
// ID must be < 2**24, and the high bit must be set.
func Encode(ID uint32) ([]byte, error) {
	return nil, fmt.Errorf("not implemented yet")
}
