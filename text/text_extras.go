package text

// Plain creates a new text.Rich without any formatting segments.
func Plain(text string) Rich {
	return Rich{Content: text}
}

// SolidColor takes in a 24-bit RGB color and overrides the alpha bits with
// 0xFF, making the color solid.
func SolidColor(rgb uint32) uint32 {
	return rgb | (0xFF << 24)
}
