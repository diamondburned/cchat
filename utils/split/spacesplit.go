package split

import (
	"unicode"
)

// These code to mostly be copied from package strings.

// SpaceIndexed returns a splitted string with the current index that
// CompleteMessage wants. The text is the entire input string and the offset is
// where the cursor currently is.
func SpaceIndexed(text string, offset int64) ([]string, int64) {
	// First count the fields.
	n, hasRunes := countSpace(text)
	if hasRunes {
		// Some runes in the input string are not ASCII.
		return spaceIndexedRunes([]rune(text), offset)
	}

	// ASCII fast path
	a := make([]string, n)
	na := int64(0)
	fieldStart := int64(0)
	i := int64(0)
	j := n - 1 // last by default

	// Skip spaces in the front of the input.
	for i < int64(len(text)) && asciiSpace[text[i]] != 0 {
		i++
	}

	fieldStart = i

	for i < int64(len(text)) {
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
		for i < int64(len(text)) && asciiSpace[text[i]] != 0 {
			i++
		}
		fieldStart = i
	}
	if fieldStart < int64(len(text)) { // Last field might end at EOF.
		a[na] = text[fieldStart:]
	}

	return a, j
}

func spaceIndexedRunes(runes []rune, offset int64) ([]string, int64) {
	// A span is used to record a slice of s of the form s[start:end].
	// The start index is inclusive and the end index is exclusive.
	type span struct{ start, end int64 }

	spans := make([]span, 0, 16)

	// Find the field start and end indices.
	wasField := false
	fromIndex := int64(0)
	for i, rune := range runes {
		if unicode.IsSpace(rune) {
			if wasField {
				spans = append(spans, span{start: fromIndex, end: int64(i)})
				wasField = false
			}
		} else {
			if !wasField {
				fromIndex = int64(i)
				wasField = true
			}
		}
	}

	// Last field might end at EOF.
	if wasField {
		spans = append(spans, span{fromIndex, int64(len(runes))})
	}

	// Create strings from recorded field indices.
	a := make([]string, 0, len(spans))
	j := int64(len(spans)) - 1 // assume last

	for i, span := range spans {
		a = append(a, string(runes[span.start:span.end]))

		if span.start <= offset && offset <= span.end {
			j = int64(i)
		}
	}

	return a, j
}
