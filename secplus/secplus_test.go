package secplus_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/zellyn/openers/bits"
	"github.com/zellyn/openers/secplus"
)

func TestManchesterEncode(t *testing.T) {
	testcases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty",
			input: "",
			want:  "",
		},
		{
			name:  "wikipedia-example",
			input: "10100111001",
			want:  "0110011010010101101001",
		},
	}

	for i, tt := range testcases {
		t.Run(fmt.Sprintf("%d-%s", i, tt.name), func(t *testing.T) {
			got, err := secplus.ManchesterEncode(bits.B(tt.input))
			if err != nil {
				t.Error(err)
				return
			}
			if !bytes.Equal(got, bits.B(tt.want)) {
				t.Errorf("want ManchesterEncode(%v)==%v; got %v", tt.input, tt.want, got)
			}
		})
	}
}

func TestManchesterDecode(t *testing.T) {
	testcases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty",
			input: "",
			want:  "",
		},
		{
			name:  "wikipedia-example",
			input: "0110011010010101101001",
			want:  "10100111001",
		},
	}

	for i, tt := range testcases {
		t.Run(fmt.Sprintf("%d-%s", i, tt.name), func(t *testing.T) {
			got, err := secplus.ManchesterDecode(bits.B(tt.input))
			if err != nil {
				t.Error(err)
				return
			}
			if !bytes.Equal(got, bits.B(tt.want)) {
				t.Errorf("want ManchesterDecode(%v)==%v; got %v", tt.input, tt.want, got)
			}
		})
	}
}

func TestEncodeV2(t *testing.T) {
	testcases := []struct {
		name       string
		fixedHigh  uint8
		fixedLow   uint64
		rolling    uint32
		want       []string
		wantBursts []string
	}{
		// Captures from my gate remotes.
		{
			name:      "capture1",
			fixedHigh: 4616223061045564932096 >> 64,
			fixedLow:  4616223061045564932096 & (1<<64 - 1),
			rolling:   240129675,
			want:      []string{"0100010000101101100001101010001111010111100100100000100110101101", "0100100001110010010110011010011110010011110110011110010010010011"},
			wantBursts: []string{
				"1010101010101010101010101010101001010101101010011010100110101010011001011001011010101001011001100110101001010101100110010101011010011010011010101010011010010110011001011001",
				"1010101010101010101010101010101001010101100110011010011010101001010110100110100110010110100101100110100101010110100110100101010110010110100101010110100110100110100110100101",
			},
		},
		{
			name:      "capture2",
			fixedHigh: 4616223061045564932096 >> 64,
			fixedLow:  4616223061045564932096 & (1<<64 - 1),
			rolling:   240129676,
			want:      []string{"0110011010110010100010110001110111101011000100100000000101010110", "0100001010001101101011111101101011101101001001101001111101111101"},
			wantBursts: []string{
				"1010101010101010101010101010101001010101101010010110100101100110010110100110011010100110010110101001010110010101011001100101101010011010011010101010101010011001100110010110",
				"1010101010101010101010101010101001010101100110011010101001100110101001011001011001100101010101011001011001100101011001011001101001101001011001101001010101011001010101011001",
			},
		},
		{
			name:      "capture3",
			fixedHigh: 4616223061045564932096 >> 64,
			fixedLow:  4616223061045564932096 & (1<<64 - 1),
			rolling:   240129677,
			want:      []string{"0110011001011011001011011100011010000110101001001101101000011111", "0100000000010110100010100110100010110100010010110010110110110110"},
			wantBursts: []string{
				"1010101010101010101010101010101001010101101010010110100101101001100101100101101001100101100101011010100101100110101010010110011001101001101001011001011001101010100101010101",
				"1010101010101010101010101010101001010101100110011010101010101010100110010110011010100110011010010110011010100110010110011010100110100110010110100110010110010110010110010110",
			},
		},
		{
			name:      "capture4",
			fixedHigh: 1222022221851718057984 >> 64,
			fixedLow:  1222022221851718057984 & (1<<64 - 1),
			rolling:   240124666,
			want:      []string{"0110001010011000111001001011100101111001110100100100110101010000", "0110010101100101001100001101101000101101001100101000101001101001"},
			wantBursts: []string{
				"1010101010101010101010101010101001010101101010010110101001100110100101101010010101101001101001100101011010011001010101101001010110011010011010011010010110011001100110101010",
				"1010101010101010101010101010101001010101100110010110100110011001011010011001101001011010101001011001011001101010011001011001101001011010011001101010011001101001011001101001",
			},
		},
		{
			name:      "capture5",
			fixedHigh: 1222022221851718057984 >> 64,
			fixedLow:  1222022221851718057984 & (1<<64 - 1),
			rolling:   240124667,
			want:      []string{"0110001001100111000110100110001000010100011001001001011000101111", "0110000110010001011010001011001010011001001010011010011001001011"},
			wantBursts: []string{
				"1010101010101010101010101010101001010101101010010110101001101001011010010101101010010110011010010110101001101010100110011010100101101001101001101001100101101010011001010101",
				"1010101010101010101010101010101001010101100110010110101010010110100110101001100101100110101001100101101001100110100101101001101001100110100101100110100101101001101001100101",
			},
		},
		{
			name:      "capture6",
			fixedHigh: 1222022221851718057984 >> 64,
			fixedLow:  1222022221851718057984 & (1<<64 - 1),
			rolling:   240124668,
			want:      []string{"0101011000001010100001000000110101100001110111110111111100010010", "0101100101100001101000001101101000101101101000001100101001101001"},
			wantBursts: []string{
				"1010101010101010101010101010101001010101101010011001100101101010101001100110011010101001101010101010010110011001011010101001010110010101010110010101010101011010100110100110",
				"1010101010101010101010101010101001010101100110011001011010011001011010101001011001101010101001011001011001101010011001011001011001101010101001011010011001101001011001101001",
			},
		},

		// Tests from the python secplus code at https://github.com/argilo/secplus/blob/dbe85c94/test_secplus.py#L73-L93
		{
			name:     "secplus-1",
			fixedLow: 70678577664,
			rolling:  240124710,
			want:     []string{"0001000100001011111000111111011011101110", "0010010110001110011110010011011011011011"},
		},
		{
			name:     "secplus-2",
			fixedLow: 70678577664,
			rolling:  240124711,
			want:     []string{"0001000010101111001100001001101101011000", "0010001000000111110101101100100110100110"},
		},
		{
			name:     "secplus-3",
			fixedLow: 70678577664,
			rolling:  240124712,
			want:     []string{"0010001001110100010101000010100110000001", "0000100101111000101000001100101100101101"},
		},
		{
			name:     "secplus-4",
			fixedLow: 70678577664,
			rolling:  240124713,
			want:     []string{"0010001000000010100011100110000010100101", "0000010110010101111001001111111011011111"},
		},
		{
			name:     "secplus-5",
			fixedLow: 62088643072,
			rolling:  240124714,
			want:     []string{"0000100010011011001010100101110011000101", "0010000110111010011010010001001011011011"},
		},
		{
			name:     "secplus-6",
			fixedLow: 62088643072,
			rolling:  240124715,
			want:     []string{"0000100001001000010000111110101000011110", "0001101000110101100001101000000100100000"},
		},
		{
			name:     "secplus-7",
			fixedLow: 62088643072,
			rolling:  240124716,
			want:     []string{"0000000000111111110111000010011111110000", "0001010110111001010001001010010011011011"},
		},
		{
			name:     "secplus-8",
			fixedLow: 62088643072,
			rolling:  240124717,
			want:     []string{"0010101010111110100111000001011111101000", "0001000110111000011010010001001011011001"},
		},
		{
			name:     "secplus-9",
			fixedLow: 66383610368,
			rolling:  240124718,
			want:     []string{"0001010101001000101001111111011010100111", "0010100110000111011110111010010011011011"},
		},
		{
			name:     "secplus-10",
			fixedLow: 66383610368,
			rolling:  240124719,
			want:     []string{"0001010100010011110011101101001000110101", "0010011000010101000101101000000100100000"},
		},
		{
			name:     "secplus-11",
			fixedLow: 66383610368,
			rolling:  240124720,
			want:     []string{"0010000010101101001111000010100100001000", "0010000101001100111100110101101111101101"},
		},
		{
			name:     "secplus-12",
			fixedLow: 66383610368,
			rolling:  240124721,
			want:     []string{"0010000001100110010110011001111111010011", "0001100110001110011010110011011111011111"},
		},
		{
			name:     "secplus-13",
			fixedLow: 74973544960,
			rolling:  240124722,
			want:     []string{"0000011000101001100111000100001111100010", "0000010110110001011101101011011111011011"},
		},
		{
			name:     "secplus-14",
			fixedLow: 74973544960,
			rolling:  240124723,
			want:     []string{"0000010110110010111000111011110000011101", "0000001000111000110000010100100110100110"},
		},
		{
			name:     "secplus-15",
			fixedLow: 74973544960,
			rolling:  240124724,
			want:     []string{"0010100101111111100011101100110011101000", "0010100101111001101001000101101100101101"},
		},
		{
			name:     "secplus-16",
			fixedLow: 74973544960,
			rolling:  240124725,
			want:     []string{"0010100100100101111000111110100001111010", "0010010110001010011110110011011111011111"},
		},
	}

	for i, tt := range testcases {
		t.Run(fmt.Sprintf("%d-%s", i, tt.name), func(t *testing.T) {
			got, err := secplus.EncodeV2(tt.fixedHigh, tt.fixedLow, tt.rolling)
			if err != nil {
				t.Error(err)
				return
			}
			if !bytes.Equal(got[0], bits.B(tt.want[0])) {
				t.Errorf("want EncodeV2(%d, %d, %d)[0]==%s; got %s", tt.fixedHigh, tt.fixedLow, tt.rolling, tt.want[0], bits.S(got[0]))
			}
			if !bytes.Equal(got[1], bits.B(tt.want[1])) {
				t.Errorf("want EncodeV2(%d, %d, %d)[1]==%s; got %s", tt.fixedHigh, tt.fixedLow, tt.rolling, tt.want[1], bits.S(got[1]))
			}

			if tt.wantBursts != nil {
				got, err = secplus.EncodeV2ToBursts(tt.fixedHigh, tt.fixedLow, tt.rolling)
				if err != nil {
					t.Error(err)
					return
				}
				if !bytes.Equal(got[0], bits.B(tt.wantBursts[0])) {
					t.Errorf("want EncodeV2ToBursts(%d, %d, %d)[0]==%s; got %s", tt.fixedHigh, tt.fixedLow, tt.rolling, tt.wantBursts[0], bits.S(got[0]))
				}
				if !bytes.Equal(got[1], bits.B(tt.wantBursts[1])) {
					t.Errorf("want EncodeV2ToBursts(%d, %d, %d)[1]==%s; got %s", tt.fixedHigh, tt.fixedLow, tt.rolling, tt.wantBursts[1], bits.S(got[1]))
				}
			}
		})
	}
}
