package secplus

import "fmt"

// orders is a map of order indicators to orders.
var orders = map[[4]byte][3]int{
	{0, 0, 0, 0}: {0, 2, 1},
	{0, 0, 0, 1}: {2, 0, 1},
	{0, 0, 1, 0}: {0, 1, 2},
	{0, 1, 0, 0}: {1, 2, 0},
	{0, 1, 0, 1}: {1, 0, 2},
	{0, 1, 1, 0}: {2, 1, 0},
	{1, 0, 0, 0}: {1, 2, 0},
	{1, 0, 0, 1}: {2, 1, 0},
	{1, 0, 1, 0}: {0, 1, 2},
}

// inversions is a map of inversion indicators to inversions.
var inversions = map[[4]byte][3]byte{
	{0, 0, 0, 0}: {1, 1, 0},
	{0, 0, 0, 1}: {0, 1, 0},
	{0, 0, 1, 0}: {0, 0, 1},
	{0, 1, 0, 0}: {1, 1, 1},
	{0, 1, 0, 1}: {1, 0, 1},
	{0, 1, 1, 0}: {0, 1, 1},
	{1, 0, 0, 0}: {1, 0, 0},
	{1, 0, 0, 1}: {0, 0, 0},
	{1, 0, 1, 0}: {1, 0, 1},
}

// syncHeader is the standard synchronization header applied to each burst.
var syncHeader []byte = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1}

// ManchesterEncode encodes a byte array of 0s and 1s to a byte array of 0s and
// 1s, but Manchester-encoded. Input is expected to be 0s and 1s.
// See https://en.wikipedia.org/wiki/Manchester_code
func ManchesterEncode(input []byte) ([]byte, error) {
	result := make([]byte, 0, len(input)*2)
	for i, b := range input {
		switch b {
		case 0:
			result = append(result, 1, 0)
		case 1:
			result = append(result, 0, 1)
		default:
			return nil, fmt.Errorf("expected only 0s and 1s; got %d at position %d", b, i)
		}
	}
	return result, nil
}

// ManchesterDecode decodes a Manchester-encoded byte array of 0s and 1s to a
// byte array of 0s and 1s.
// See https://en.wikipedia.org/wiki/Manchester_code
func ManchesterDecode(input []byte) ([]byte, error) {
	if len(input)%2 != 0 {
		return nil, fmt.Errorf("expected even length; got %d", len(input))
	}
	result := make([]byte, 0, len(input)/2)
	for i := 0; i < len(input); i += 2 {
		b1, b2 := input[i], input[i+1]
		if b1 > 1 {
			return nil, fmt.Errorf("expected only 0s and 1s; got %d at position %d", b1, i)
		}
		if b2 > 1 {
			return nil, fmt.Errorf("expected only 0s and 1s; got %d at position %d", b2, i+1)
		}
		if b1 == b2 {
			return nil, fmt.Errorf("pairs should be 01 or 10; got %d%d at position %d", b1, b2, i)
		}
		result = append(result, b2)
	}
	return result, nil
}

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

// EnvodeV2ToBursts encodes a Security+2.0 fixed and rolling code into two
// bitstreams, one for each half. The bitstreams have a standard prefix, and are
// then Manchester coded.
func EncodeV2ToBursts(fixedHigh uint8, fixedLow uint64, rolling uint32) ([2][]byte, error) {
	packets, err := EncodeV2(fixedHigh, fixedLow, rolling)
	if err != nil {
		return [2][]byte{}, err
	}

	var result [2][]byte

	for i, bits := range packets {
		bitsWithHeader := make([]byte, len(syncHeader)+2+len(bits))
		copy(bitsWithHeader, syncHeader)
		bitsWithHeader[len(syncHeader)+1] = byte(i)
		copy(bitsWithHeader[len(syncHeader)+2:], bits)
		result[i], err = ManchesterEncode(bitsWithHeader)
		if err != nil {
			return [2][]byte{}, err
		}
	}

	return result, nil
}

// EncodeV2 encodes a Security+2.0 fixed and rolling code into two packets, one
// for each half. It supports short (< 2**40) and long (< 2**72) fixed codes.
// Rolling code must be < 2**28.
func EncodeV2(fixedHigh uint8, fixedLow uint64, rolling uint32) ([2][]byte, error) {
	fixedHalves, long := getFixedHalves(fixedHigh, fixedLow)
	ternaryHalves, err := getRollingTernaryHalves(rolling)
	if err != nil {
		return [2][]byte{}, err
	}

	h1, err := encodeHalfV2(fixedHalves[0], ternaryHalves[0], long)
	if err != nil {
		return [2][]byte{}, err
	}
	h2, err := encodeHalfV2(fixedHalves[1], ternaryHalves[1], long)
	if err != nil {
		return [2][]byte{}, err
	}

	return [2][]byte{h1, h2}, nil
}

// encodeHalfV2 encodes half of a v2 code. The first half of the fixed code, and
// the first half of the ternary-"encrypted" rolling code are encoded completely
// separately from the second halves.
func encodeHalfV2(fixed []byte, rolling []byte, long bool) ([]byte, error) {
	result := make([]byte, 2, 64)
	partLength := 10
	if long {
		result[1] = 1
		partLength = 18
	}
	parts := [3][]byte{fixed[0:partLength], fixed[partLength : partLength*2], {}}

	result = append(result, rolling[:8]...)
	rest := make([]byte, partLength)
	copy(rest, rolling[8:])
	if long {
		copy(rest[10:], rolling[:8])
	}
	parts[2] = rest

	var orderIndicator [4]byte
	copy(orderIndicator[:], rolling[:4])
	var inversionIndicator [4]byte
	copy(inversionIndicator[:], rolling[4:8])
	order, ok := orders[orderIndicator]
	if !ok {
		return nil, fmt.Errorf("no order found for indicator %v", orderIndicator)
	}
	invert, ok := inversions[inversionIndicator]
	if !ok {
		return nil, fmt.Errorf("no inversion found for indicator %v", inversionIndicator)
	}

	parts[0], parts[1], parts[2] = parts[order[0]], parts[order[1]], parts[order[2]]

	for i := 0; i < partLength; i++ {
		for j := 0; j < 3; j++ {
			result = append(result, parts[j][i]^invert[j])
		}
	}

	return result, nil
}

// getFixedHalves returns the first and second halves of the fixed part, as
// []byte of 0s and 1s. If fixed < 2**40, it returns 20-bit halves, otherwise it
// returns 36-bit halves.
func getFixedHalves(fixedHigh uint8, fixedLow uint64) ([2][]byte, bool) {
	long := false
	bits := fmt.Sprintf("%064b", fixedLow)
	if fixedHigh > 0 || fixedLow >= 1<<40 {
		long = true
		bits = fmt.Sprintf("%08b", fixedHigh) + bits
	}

	if long {
		return [2][]byte{B(bits[0:36]), B(bits[36:72])}, long
	}
	return [2][]byte{B(bits[24:44]), B(bits[44:64])}, long
}

// getRollingTernaryHalves converts the rolling code into a binary-pair-coded
// ternary representation, split between two halves: a byte array of 0s and 1s.
func getRollingTernaryHalves(rolling uint32) ([2][]byte, error) {
	if rolling >= 1<<28 {
		return [2][]byte{}, fmt.Errorf("rolling code must be <= 2^28; got %d", rolling)
	}
	bitReversed := uint32(0)
	for i := 0; i < 28; i++ {
		bitReversed = bitReversed*2 + rolling&1
		rolling >>= 1
	}
	ternary := make([]byte, 0, 18)
	for i := 0; i < 18; i++ {
		ternary = append(ternary, byte(bitReversed%3))
		bitReversed /= 3
	}

	ternaryPieces := func(pieces [][]byte) []byte {
		var result []byte
		for _, piece := range pieces {
			for i := len(piece) - 1; i >= 0; i-- {
				trit := piece[i]
				result = append(result, trit>>1, trit&1)
			}
		}
		return result
	}
	part1 := ternaryPieces([][]byte{ternary[0:4], ternary[8:12], ternary[16:17]})
	part2 := ternaryPieces([][]byte{ternary[4:8], ternary[12:16], ternary[17:18]})
	return [2][]byte{part1, part2}, nil
}
