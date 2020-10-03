// DO NOT EDIT: THIS FILE IS GENERATED!

// Package text provides a rich text API for cchat interfaces to use.
//
//
// Asserting
//
// Although interfaces here contain asserter methods similarly to cchat, the
// backend should take care to not implement multiple interfaces that may seem
// conflicting. For example, if Avatarer is already implemented, then Imager
// shouldn't be.
package text

// Attribute is the type for basic rich text markup attributes.
type Attribute uint32

const (
	// Normal is a zero-value attribute.
	AttributeNormal Attribute = iota
	// Bold represents bold text.
	AttributeBold
	// Italics represents italicized text.
	AttributeItalics
	// Underline represents underlined text.
	AttributeUnderline
	// Strikethrough represents struckthrough text.
	AttributeStrikethrough
	// Spoiler represents spoiler text, which usually looks blacked out until
	// hovered or clicked on.
	AttributeSpoiler
	// Monospace represents monospaced text, typically for inline code.
	AttributeMonospace
	// Dimmed represents dimmed text, typically slightly less visible than other
	// text.
	AttributeDimmed
)

func (a Attribute) Is(is Attribute) bool {
	return a == is
}

// Rich is a normal text wrapped with optional format segments.
type Rich struct {
	Content  string
	Segments []Segment
}

// Attributor is a rich text markup format that a segment could implement. This
// is to be applied directly onto the text.
type Attributor interface {
	Attribute() Attribute
}

// Avatarer implies the segment should be replaced with a rounded-corners image.
// This works similarly to Imager.
//
// For segments that also implement mentioner, the image should be treated as a
// round avatar.
type Avatarer interface {
	// AvatarText returns the underlying text of the image. Frontends could use this
	// for hovering or displaying the text instead of the image.
	AvatarText() string
	// AvatarSize returns the requested dimension for the image. This function could
	// return (0, 0), which the frontend should use the avatar's dimensions.
	AvatarSize() (size int)
	// Avatar returns the URL for the image.
	Avatar() (url string)
}

// Codeblocker is a codeblock that supports optional syntax highlighting using
// the language given. Note that as this is a block, it will appear separately
// from the rest of the paragraph.
//
// This interface is equivalent to Markdown's codeblock syntax.
type Codeblocker interface {
	CodeblockLanguage() (language string)
}

// Colorer is a text color format that a segment could implement. This is to be
// applied directly onto the text.
//
// The Color method must return a valid 32-bit RGBA color. That is, if the text
// color is solid, then the alpha value must be 0xFF. Frontends that support
// 32-bit colors must render alpha accordingly without any edge cases.
type Colorer interface {
	// Color returns a 32-bit RGBA color.
	Color() uint32
}

// Imager implies the segment should be replaced with a (possibly inlined)
// image. Only the starting bound matters, as images cannot substitute texts.
//
// For segments that also implement mentioner, the image should be treated as a
// square avatar.
type Imager interface {
	// ImageText returns the underlying text of the image. Frontends could use this
	// for hovering or displaying the text instead of the image.
	ImageText() string
	// ImageSize returns the requested dimension for the image. This function could
	// return (0, 0), which the frontend should use the image's dimensions.
	ImageSize() (w int, h int)
	// Image returns the URL for the image.
	Image() (url string)
}

// Linker is a hyperlink format that a segment could implement. This implies
// that the segment should be replaced with a hyperlink, similarly to the anchor
// tag with href being the URL and the inner text being the text string.
type Linker interface {
	Link() (url string)
}

// Mentioner implies that the segment can be clickable, and when clicked it
// should open up a dialog containing information from MentionInfo().
//
// It is worth mentioning that frontends should assume whatever segment that
// Mentioner highlighted to be the display name of that user. This would allow
// frontends to flexibly layout the labels.
type Mentioner interface {
	// MentionInfo returns the popup information of the mentioned segment. This is
	// typically user information or something similar to that context.
	MentionInfo() Rich
}

// Quoteblocker represents a quoteblock that behaves similarly to the blockquote
// HTML tag. The quoteblock may be represented typically by an actaul quoteblock
// or with green arrows prepended to each line.
type Quoteblocker interface {
	// QuotePrefix returns the prefix that every line the segment covers have. This
	// is typically the greater-than sign ">" in Markdown. Frontends could use this
	// information to format the quote properly.
	QuotePrefix() (prefix string)
}

// Segment is the minimum requirement for a format segment. Frontends will use
// this to determine when the format starts and ends. They will also assert this
// interface to any other formatting interface, including Linker, Colorer and
// Attributor.
//
// Note that a segment may implement multiple interfaces. For example, a
// Mentioner may also implement Colorer.
type Segment interface {
	Bounds() (start int, end int)

	// Asserters.

	AsColorer() Colorer           // Optional
	AsLinker() Linker             // Optional
	AsImager() Imager             // Optional
	AsAvatarer() Avatarer         // Optional
	AsMentioner() Mentioner       // Optional
	AsAttributor() Attributor     // Optional
	AsCodeblocker() Codeblocker   // Optional
	AsQuoteblocker() Quoteblocker // Optional
}
