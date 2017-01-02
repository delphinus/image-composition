package imco

import (
	"image"
	"image/gif"
	"os"

	"github.com/pkg/errors"
)

// OverlayImage makes an image overlayed by the other
func OverlayImage(input, overlay, output string) error {
	inputImage, err := loadImage(input)
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

	image, err := gif.Decode(file)
	if err != nil {
		return nil, errors.Wrapf(err, "error in Decoding '%s'", filename)
	}

	return &image, err
}

func processOverlay(inputImage, overlayImage *image.Image) (*gif.GIF, error) {
	return nil, nil
}
