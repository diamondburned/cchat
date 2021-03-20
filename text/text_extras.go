package text

// Plain creates a new text.Rich without any formatting segments.
func Plain(text string) Rich {
	return Rich{Content: text}
}

// SolidColor takes in a 24-bit RGB color and overrides the alpha bits with
// 0xFF, making the color solid.
func SolidColor(rgb uint32) uint32 {
	return (rgb << 8) | 0xFF
}

// IsEmpty returns true if the given rich segment's content is empty. Note that
// a rich text is not necessarily empty if the content is empty, because there
// may be images within the segments.
func (r Rich) IsEmpty() bool {
	return r.Content == "" && len(r.Segments) == 0
}
