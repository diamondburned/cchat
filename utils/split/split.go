// Package split provides a simple string splitting utility for use with
// CompleteMessage.
package split

import "unicode/utf8"

// SplitFunc is the type that describes the two splitter functions, SpaceIndexed
// and ArgsIndexed.
type SplitFunc = func(text string, offset int64) ([]string, int64)

var (
	_ SplitFunc = SpaceIndexed
	_ SplitFunc = ArgsIndexed
)

// just helper functions here

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

func isSpace(b byte) bool { return asciiSpace[b] == 1 }

func countSpace(text string) (n int64, hasRunes bool) {
	// This implementation to also be mostly copy-pasted from package strings.

	// This is an exact count if s is ASCII, otherwise it is an approximation.
	// var n int64

	wasSpace := int64(1)
	// setBits is used to track which bits are set in the bytes of s.
	setBits := uint8(0)
	for i := 0; i < len(text); i++ {
		r := text[i]
		setBits |= r
		isSpace := int64(asciiSpace[r])
		n += wasSpace & ^isSpace
		wasSpace = isSpace
	}

	return n, setBits >= utf8.RuneSelf
}
