package bits

// S converts from a byte slice of 0s and 1s to string of ASCII 0s and 1s. It
// does no validation, since it uses only the low bit.
func S(input []byte) string {
	result := make([]byte, len(input))
	for i, b := range input {
		result[i] = (b & 1) + '0'
	}
	return string(result)
}

// B converts from a string of ASCII 0s and 1s to a byte slice of 0s and 1s. It
// does no validation, since it only uses the low bit.
func B(input string) []byte {
	result := make([]byte, len(input))
	for i, b := range input {
		result[i] = byte((b & 1))
	}
	return result
}
