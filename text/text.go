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
	Link() (text, url string)
}

// Imager implies the segment should be replaced with a (possibly inlined)
// image.
type Imager interface {
	Image() (url string)
}

// Colorer is a text color format that a segment could implement. This is to be
// applied directly onto the text.
type Colorer interface {
	Color() uint16
}

// Attributor is a rich text markup format that a segment could implement. This
// is to be applied directly onto the text.
type Attributor interface {
	Attribute() Attribute
}

// Attribute is the type for basic rich text markup attributes.
type Attribute uint16

const (
	AttrBold Attribute = 1 << iota
	AttrItalics
	AttrUnderline
	AttrStrikethrough
	AttrSpoiler
	AttrMonospace
	AttrQuoted
)
