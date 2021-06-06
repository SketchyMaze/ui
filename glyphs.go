package ui

/*
Glyph images as Base64 encoded PNGs.
*/

import (
	"encoding/base64"
	"image"
	"image/png"
	"bytes"
)

// List of available glyphs.
const (
	// Downward pointed black arrow 9x9 pixels.
	GlyphDownArrow9x9 = `iVBORw0KGgoAAAANSUhEUgAAAAkAAAAJCAYAAADgkQYQAAAABmJLR0QA/wD/AP+gvaeTAAAACXBI
WXMAAC4jAAAuIwF4pT92AAAAKklEQVQY02NgoBZgZGBg+E+MIgYCChkZkTj/cRnCiCb4H4stWMF/
BpoBAGQSBQOpAugRAAAAAElFTkSuQmCC`
)

// GetGlyph loads a PNG image from a hard-coded glyph.
func GetGlyph(b64 string) (image.Image, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}

	scanner := bytes.NewReader(data)
	return png.Decode(scanner)
}