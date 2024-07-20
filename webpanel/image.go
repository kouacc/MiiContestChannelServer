package webpanel

import (
	"bufio"
	"bytes"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"io"

	// Importing as a side effect allows for the image library to check for these formats
	_ "image/gif"
	_ "image/png"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

func resize(original io.Reader, x, y int) ([]byte, error) {
	decodedImage, _, err := image.Decode(original)
	if err != nil {
		return nil, err
	}

	newImage := image.NewRGBA(image.Rect(0, 0, x, y))
	draw.BiLinear.Scale(newImage, newImage.Bounds(), decodedImage, decodedImage.Bounds(), draw.Over, nil)

	var out bytes.Buffer
	err = jpeg.Encode(bufio.NewWriter(&out), newImage, nil)
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
