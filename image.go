package imco

import (
	"image"
	"image/gif"
	"image/png"
	"os"
	"path"

	"github.com/pkg/errors"
)

// OverlayImage makes an image overlayed by the other
func OverlayImage(input, overlay, output string) error {
	inputImage, err := loadAnimationGIF(input)
	if err != nil {
		return errors.Wrap(err, "error in loadImage")
	}

	overlayImage, err := loadImage(overlay)
	if err != nil {
		return errors.Wrap(err, "error in loadImage")
	}

	outputGIF, err := processOverlay(inputImage, overlayImage)
	if err != nil {
		return errors.Wrap(err, "error in processOverlay")
	}

	outputFile, err := os.Create(output)
	if err != nil {
		return errors.Wrapf(err, "error in Create '%s'", output)
	}

	if err := gif.EncodeAll(outputFile, outputGIF); err != nil {
		return errors.Wrap(err, "error in EncodeAll")
	}
	return nil
}

func loadImage(filename string) (*image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "error in Opening '%s'", filename)
	}

	var img image.Image
	switch path.Ext(filename) {
	case ".gif":
		img, err = gif.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	default:
		return nil, errors.Errorf("image must be GIF or PNG: %s", filename)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "error in Decoding '%s'", filename)
	}

	return &img, err
}

func loadAnimationGIF(filename string) (*gif.GIF, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "error in Opening '%s'", filename)
	}
	return gif.DecodeAll(file)
}

func processOverlay(inputImage *gif.GIF, overlayImage *image.Image) (*gif.GIF, error) {
	return nil, nil
}
