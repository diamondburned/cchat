package text

// Rich is a normal text wrapped with optional format segments.
type Rich struct {
	Content string
	// Segments are optional rich-text segment markers.
	Segments []Segment
}

func (r Rich) Empty() bool {
	return r.Content == ""
}

// Segment is the minimum requirement for a format segment. Frontends will use
// this to determine when the format starts and ends. They will also assert this
// interface to any other formatting interface, including Linker, Colorer and
// Attributor.
type Segment interface {
	Bounds() (start, end int)
}

// Linker is a hyperlink format that a segment could implement. This implies
// that the segment should be replaced with a hyperlink, similarly to the anchor
// tag with href being the URL and the inner text being the text string.
type Linker interface {
	Segment
	Link() (text, url string)
}

// Imager implies the segment should be replaced with a (possibly inlined)
// image.
type Imager interface {
	Segment
	Image() (url string)
}

// Colorer is a text color format that a segment could implement. This is to be
// applied directly onto the text.
type Colorer interface {
	Segment
	Color() uint32
}

// Attributor is a rich text markup format that a segment could implement. This
// is to be applied directly onto the text.
type Attributor interface {
	Segment
	Attribute() Attribute
}

// Attribute is the type for basic rich text markup attributes.
type Attribute uint16

// HasAttr returns whether or not "attr" has "this" attribute.
func (attr Attribute) Has(this Attribute) bool {
	return (attr & this) == this
}

const (
	// AttrBold represents bold text.
	AttrBold Attribute = 1 << iota
	// AttrItalics represents italicized text.
	AttrItalics
	// AttrUnderline represents underlined text.
	AttrUnderline
	// AttrStrikethrough represents strikethrough text.
	AttrStrikethrough
	// AttrSpoiler represents spoiler text, which usually looks blacked out
	// until hovered or clicked on.
	AttrSpoiler
	// AttrMonospace represents monospaced text, typically for inline code.
	AttrMonospace
	// AttrDimmed represents dimmed text, typically slightly less visible than
	// other text.
	AttrDimmed
)

// Codeblocker is a codeblock that supports optional syntax highlighting using
// the language given. Note that as this is a block, it will appear separately
// from the rest of the paragraph.
//
// This interface is equivalent to Markdown's codeblock syntax.
type Codeblocker interface {
	Segment
	CodeblockLanguage() string
}

// Quoteblocker represents a quoteblock that behaves similarly to the blockquote
// HTML tag. The quoteblock may be represented typically by an actaul quoteblock
// or with green arrows prepended to each line.
type Quoteblocker interface {
	Segment
	// Quote does nothing; it's only here to distinguish the interface.
	Quote()
}
