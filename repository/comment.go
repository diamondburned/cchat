package repository

import (
	"bytes"
	"go/doc"
	"regexp"
	"strings"
)

var (
	// commentTrimSurrounding is a regex to trim surrounding new line and tabs.
	// This is needed to find the correct level of indentation.
	commentTrimSurrounding = regexp.MustCompile(`(^\n)|(\n\t+$)`)
)

// TabWidth is used to format comments.
var TabWidth = 4

// Comment represents a raw comment string. Most use cases should use GoString()
// to get the comment's content.
type Comment struct {
	Raw string
}

// IsEmpty returns true if the comment is empty.
func (c Comment) IsEmpty() bool {
	return c.Raw == ""
}

// GoString formats the documentation string in 80 columns wide paragraphs and
// prefix each line with "// ". The ident argument controls the nested level. If
// less than or equal to zero, then it is changed to 1, which is the top level.
func (c Comment) GoString(ident int) string {
	if c.Raw == "" {
		return ""
	}

	if ident < 1 {
		ident = 1
	}

	ident-- // 0th-indexed
	ident *= TabWidth

	var lines = strings.Split(c.WrapText(80-len("// ")-ident), "\n")
	for i, line := range lines {
		if line != "" {
			line = "// " + line
		} else {
			line = "//"
		}
		lines[i] = line
	}

	return strings.Join(lines, "\n")
}

// WrapText wraps the raw text in n columns wide paragraphs.
func (c Comment) WrapText(column int) string {
	var txt = c.Unindent()
	if txt == "" {
		return ""
	}

	buf := bytes.Buffer{}
	doc.ToText(&buf, txt, "", strings.Repeat(" ", TabWidth-1), column)

	text := strings.TrimRight(buf.String(), "\n")
	text = strings.Replace(text, "\t", strings.Repeat(" ", TabWidth), -1)

	return text
}

// Unindent removes the indentations that were there for the sake of syntax in
// RawText. It gets the lowest indentation level from each line and trim it.
func (c Comment) Unindent() string {
	if c.IsEmpty() {
		return ""
	}

	// Trim new lines.
	txt := commentTrimSurrounding.ReplaceAllString(c.Raw, "")

	// Split the lines and rejoin them later to trim the indentation.
	var lines = strings.Split(txt, "\n")
	var indent = 0

	// Get the minimum indentation count.
	for _, line := range lines {
		linedent := strings.Count(line, "\t")
		if linedent < 0 {
			continue
		}
		if linedent < indent || indent == 0 {
			indent = linedent
		}
	}

	// Trim the indentation.
	if indent > 0 {
		for i, line := range lines {
			if len(line) > 0 {
				lines[i] = line[indent-1:]
			}
		}
	}

	// Rejoin.
	txt = strings.Join(lines, "\n")

	return txt
}
