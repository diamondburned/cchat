// Package split provides a simple string splitting utility for use with
// CompleteMessage.
package split

import (
	"unicode"
	"unicode/utf8"
)

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

// SpaceIndexed returns a splitted string with the current index that
// CompleteMessage wants. The text is the entire input string and the offset is
// where the cursor currently is.
func SpaceIndexed(text string, offset int) ([]string, int) {
	// First count the fields.
	// This is an exact count if s is ASCII, otherwise it is an approximation.
	n := 0
	wasSpace := 1
	// setBits is used to track which bits are set in the bytes of s.
	setBits := uint8(0)
	for i := 0; i < len(text); i++ {
		r := text[i]
		setBits |= r
		isSpace := int(asciiSpace[r])
		n += wasSpace & ^isSpace
		wasSpace = isSpace
	}

	if setBits >= utf8.RuneSelf {
		// Some runes in the input string are not ASCII.
		return spaceIndexedRunes([]rune(text), offset)
	}

	// ASCII fast path
	a := make([]string, n)
	na := 0
	fieldStart := 0
	i := 0
	j := n - 1 // last by default

	// Skip spaces in the front of the input.
	for i < len(text) && asciiSpace[text[i]] != 0 {
		i++
	}

	fieldStart = i

	for i < len(text) {
		if asciiSpace[text[i]] == 0 {
			i++
			continue
		}

		a[na] = text[fieldStart:i]
		if fieldStart <= offset && offset <= i {
			j = na
		}

		na++
		i++

		// Skip spaces in between fields.
		for i < len(text) && asciiSpace[text[i]] != 0 {
			i++
		}
		fieldStart = i
	}
	if fieldStart < len(text) { // Last field might end at EOF.
		a[na] = text[fieldStart:]
	}

	return a, j
}

func spaceIndexedRunes(runes []rune, offset int) ([]string, int) {
	// A span is used to record a slice of s of the form s[start:end].
	// The start index is inclusive and the end index is exclusive.
	type span struct{ start, end int }

	spans := make([]span, 0, 16)

	// Find the field start and end indices.
	wasField := false
	fromIndex := 0
	for i, rune := range runes {
		if unicode.IsSpace(rune) {
			if wasField {
				spans = append(spans, span{start: fromIndex, end: i})
				wasField = false
			}
		} else {
			if !wasField {
				fromIndex = i
				wasField = true
			}
		}
	}

	// Last field might end at EOF.
	if wasField {
		spans = append(spans, span{fromIndex, len(runes)})
	}

	// Create strings from recorded field indices.
	a := make([]string, 0, len(spans))
	j := len(spans) - 1 // assume last

	for i, span := range spans {
		a = append(a, string(runes[span.start:span.end]))

		if span.start <= offset && offset <= span.end {
			j = i
		}
	}

	return a, j
}
