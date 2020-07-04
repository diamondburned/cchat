package text

// Rich is a normal text wrapped with optional format segments.
type Rich struct {
	Content string
	// Segments are optional rich-text segment markers.
	Segments []Segment
}

// Empty returns whether or not the rich text is considered empty.
func (r Rich) Empty() bool {
	return r.Content == ""
}

// String returns the content. This is used mainly for printing.
func (r Rich) String() string {
	return r.Content
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
	Link() (url string)
}

// Imager implies the segment should be replaced with a (possibly inlined)
// image. Only the starting bound matters, as images cannot substitute texts.
type Imager interface {
	Segment
	// Image returns the URL for the image.
	Image() (url string)
	// ImageSize returns the requested dimension for the image. This function
	// could return (0, 0), which the frontend should use the image's
	// dimensions.
	ImageSize() (w, h int)
	// ImageText returns the underlying text of the image. Frontends could use
	// this for hovering or displaying the text instead of the image.
	ImageText() string
}

// Avatarer implies the segment should be replaced with a rounded-corners
// image. This works similarly to Imager.
type Avatarer interface {
	Segment
	// Avatar returns the URL for the image.
	Avatar() (url string)
	// AvatarSize returns the requested dimension for the image. This function
	// could return (0, 0), which the frontend should use the avatar's
	// dimensions.
	AvatarSize() (size int)
	// AvatarText returns the underlying text of the image. Frontends could use
	// this for hovering or displaying the text instead of the image.
	AvatarText() string
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
