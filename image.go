package ui

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"git.kirsle.net/go/render"
	"golang.org/x/image/bmp"
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
	Image   image.Image     // a Go image version
	texture render.Texturer // (SDL2) Texture, lazy inited on Present.
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

// ImageFromImage creates an Image from a Go standard library image.Image.
func ImageFromImage(im image.Image) (*Image, error) {
	return &Image{
		Type:  PNG,
		Image: im,
	}, nil
}

// ImageFromTexture creates an Image from a texture.
func ImageFromTexture(tex render.Texturer) *Image {
	return &Image{
		texture: tex,
		Image:   tex.Image(),
	}
}

// ImageFromFile creates an Image by opening a file from disk.
func ImageFromFile(filename string) (*Image, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	img, err := jpeg.Decode(fh)
	if err != nil {
		return nil, err
	}

	return &Image{
		Image: img,
	}, nil
}

// ReplaceFromImage replaces the image with a new image.
func (w *Image) ReplaceFromImage(im image.Image) error {
	// Free the old texture.
	if w.texture != nil {
		if err := w.texture.Free(); err != nil {
			return err
		}
		w.texture = nil
	}
	w.Image = im
	return nil
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

	// Open the file from disk.
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// Parse it.
	switch w.Type {
	case PNG:
		img, err := png.Decode(fh)
		if err != nil {
			return nil, err
		}
		w.Image = img
	case JPEG:
		img, err := jpeg.Decode(fh)
		if err != nil {
			return nil, err
		}
		w.Image = img
	case BMP:
		img, err := bmp.Decode(fh)
		if err != nil {
			return nil, err
		}
		w.Image = img
	}

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

// Size returns the dimensions of the image which is also the widget's size.
func (w *Image) Size() render.Rect {
	if w.Image != nil {
		bounds := w.Image.Bounds().Canon()
		return render.Rect{
			W: bounds.Max.X,
			H: bounds.Max.Y,
		}
	}
	return w.BaseWidget.Size()
}

// Counter for unique SDL2 texture names.
var __imageID int

// Present the widget. This should be called on your main thread
// if using SDL2 in case it needs to generate textures.
func (w *Image) Present(e render.Engine, p render.Point) {
	// Lazy load the (e.g. SDL2) texture from the stored bitmap.
	if w.texture == nil {
		if w.Image == nil {
			return
		}

		__imageID++
		tex, err := e.StoreTexture(fmt.Sprintf("ui.Image(%d).png", __imageID), w.Image)
		if err != nil {
			fmt.Printf("ui.Image.Present(): could not make texture: %s\n", err)
			return
		}
		w.texture = tex
	}

	size := w.texture.Size()
	dst := render.Rect{
		X: p.X,
		Y: p.Y,
		W: size.W,
		H: size.H,
	}
	e.Copy(w.texture, size, dst)

	// Call the BaseWidget Present in case we have subscribers.
	w.BaseWidget.Present(e, p)
}

// Destroy cleans up the image and releases textures.
func (w *Image) Destroy() {
	if w.texture != nil {
		w.texture.Free()
		w.texture = nil
	}
}
