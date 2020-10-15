package split

import (
	"bytes"
)

// The original shellwords implementation belongs to mattn. This version alters
// some code along with some trivial optimizations.

// ArgsIndexed converts text into a shellwords-split list of strings. This
// function roughly follows shell syntax, in that it works with a single space
// character in bytes. This means that implementations could use bytes, as long
// as it only works with characters within the ASCII range, assuming the source
// text is in UTF-8 (which it should, per Go specifications).
func ArgsIndexed(text string, offset int64) (args []string, argIndex int64) {
	// Quickly loop over everything to roughly count spaces. It doesn't have to
	// be accurate. This isn't very useful, to be honest.
	args = make([]string, 0, approxArgLen(text))

	var escaped, doubleQuoted, singleQuoted bool
	argIndex = -1 // in case

	buf := bytes.Buffer{}
	buf.Grow(len(text))

	got := false
	cursor := 0

	for i, length := int64(0), int64(len(text)); i < length; i++ {
		r := text[i]
		if offset == i {
			argIndex = int64(len(args))
		}

		switch {
		case escaped:
			got = true
			escaped = false

			if doubleQuoted {
				switch r {
				case 'n':
					buf.WriteByte('\n')
					continue
				case 't':
					buf.WriteByte('\t')
					continue
				}
			}
			buf.WriteByte(r)
			continue

		case isSpace(r):
			switch {
			case singleQuoted, doubleQuoted:
				buf.WriteByte(r)
			case got:
				cursor += buf.Len()
				args = append(args, buf.String())
				buf.Reset()
				got = false
			}
			continue
		}

		switch r {
		case '\\':
			if singleQuoted {
				buf.WriteByte(r)
			} else {
				escaped = true
			}
			continue

		case '"':
			if !singleQuoted {
				if doubleQuoted {
					got = true
				}
				doubleQuoted = !doubleQuoted
				continue
			}

		case '\'':
			if !doubleQuoted {
				if singleQuoted {
					got = true
				}
				singleQuoted = !singleQuoted
				continue
			}
		}

		got = true
		buf.WriteByte(r)
	}

	if got || escaped || singleQuoted || doubleQuoted {
		if argIndex < 0 {
			argIndex = int64(len(args))
		}

		args = append(args, buf.String())
	}

	return
}

// this is completely optional.
func approxArgLen(text string) int {
	var arglen int
	var inside bool
	var escape bool
	for i := 0; i < len(text); i++ {
		switch b := text[i]; b {
		case '\\':
			escape = true
			continue
		case '\'', '"':
			if !escape {
				inside = !inside
			}
		default:
			if isSpace(b) && !inside {
				arglen++
			}
		}

		if escape {
			escape = false
		}
	}

	// Allocate 1 more just in case.
	return arglen + 1
}
