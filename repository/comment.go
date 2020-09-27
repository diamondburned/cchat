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

// Comment represents a raw comment string. Most use cases should use GoString()
// to get the comment's content.
type Comment struct {
	RawText string
}

// GoString formats the documentation string in 80 columns wide paragraphs.
func (c Comment) GoString() string {
	return c.WrapText(80)
}

// WrapText wraps the raw text in n columns wide paragraphs.
func (c Comment) WrapText(column int) string {
	var buf bytes.Buffer
	doc.ToText(&buf, c.Unindent(), "", "\t", column)
	return buf.String()
}

// Unindent removes the indentations that were there for the sake of syntax in
// RawText. It gets the lowest indentation level from each line and trim it.
func (c Comment) Unindent() string {
	// Trim new lines.
	txt := commentTrimSurrounding.ReplaceAllString(c.RawText, "")

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
