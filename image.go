package imco

import (
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"os"
	"path"

	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

const (
	targetX = 300
	targetY = 300
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
	outputGIF := gif.GIF{
		Image:           make([]*image.Paletted, len(inputImage.Image)),
		Delay:           inputImage.Delay,
		Disposal:        inputImage.Disposal,
		BackgroundIndex: inputImage.BackgroundIndex,
	}

	firstFrameBounds := inputImage.Image[0].Bounds()
	b := image.Rect(0, 0, firstFrameBounds.Dx(), firstFrameBounds.Dy())
	frameImage := image.NewRGBA(b)

	for i, frame := range inputImage.Image {
		if i == 0 {
			firstFrameBounds = frame.Bounds()
		}
		bounds := frame.Bounds()
		draw.Draw(frameImage, bounds, frame, bounds.Min, draw.Over)
		outputGIF.Image[i] = imageToPaletted(resizeImage(frameImage))
	}

	return &outputGIF, nil
}

func resizeImage(img image.Image) image.Image {
	return resize.Resize(targetX, targetY, img, resize.Bilinear)
}

func imageToPaletted(img image.Image) *image.Paletted {
	bounds := img.Bounds()
	paletted := image.NewPaletted(bounds, palette.Plan9)
	draw.FloydSteinberg.Draw(paletted, bounds, img, image.ZP)
	return paletted
}
