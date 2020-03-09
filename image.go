package ui

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"git.kirsle.net/go/render"
)

// ImageType for supported image formats.
type ImageType string

// Supported image formats.
const (
	BMP  ImageType = "bmp"
	PNG            = "png"
	JPEG           = "jpg"
)

// Image is a widget that is backed by an image file.
type Image struct {
	BaseWidget

	// Configurable fields for the constructor.
	Type    ImageType
	Image   image.Image
	texture render.Texturer
}

// NewImage creates a new Image.
func NewImage(c Image) *Image {
	w := &Image{
		Type: c.Type,
	}
	if w.Type == "" {
		w.Type = BMP
	}

	w.IDFunc(func() string {
		return fmt.Sprintf(`Image<"%s">`, w.Type)
	})
	return w
}

// ImageFromTexture creates an Image from a texture.
func ImageFromTexture(tex render.Texturer) *Image {
	return &Image{
		texture: tex,
	}
}

// ImageFromFile creates an Image by opening a file from disk.
func ImageFromFile(e render.Engine, filename string) (*Image, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	img, err := jpeg.Decode(fh)
	if err != nil {
		return nil, err
	}

	tex, err := e.StoreTexture(filename, img)
	if err != nil {
		return nil, err
	}

	return &Image{
		Image:   img,
		texture: tex,
	}, nil
}

// OpenImage initializes an Image with a given file name.
//
// The file extension is important and should be a supported ImageType.
func OpenImage(e render.Engine, filename string) (*Image, error) {
	w := &Image{}
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".bmp":
		w.Type = BMP
	case ".png":
		w.Type = PNG
	case ".jpg":
		w.Type = JPEG
	case ".jpeg":
		w.Type = JPEG
	default:
		return nil, fmt.Errorf("OpenImage: %s: not a supported image type", filename)
	}

	tex, err := e.LoadTexture(filename)
	if err != nil {
		return nil, err
	}

	w.texture = tex
	return w, nil
}

// GetRGBA returns an image.RGBA from the image data.
func (w *Image) GetRGBA() *image.RGBA {
	var bounds = w.Image.Bounds()
	var rgba = image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			color := w.Image.At(x, y)
			rgba.Set(x, y, color)
		}
	}
	return rgba
}

// Compute the widget.
func (w *Image) Compute(e render.Engine) {
	w.Resize(w.texture.Size())
}

// Present the widget.
func (w *Image) Present(e render.Engine, p render.Point) {
	size := w.texture.Size()
	dst := render.Rect{
		X: p.X,
		Y: p.Y,
		W: size.W,
		H: size.H,
	}
	e.Copy(w.texture, size, dst)
}
