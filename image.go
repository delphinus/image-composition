package imco

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"os"
	"path"

	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"github.com/soniakeys/quant/median"
)

// OverlayImage makes an image overlayed by the other
func OverlayImage(input, overlay, output string, width, height uint) error {
	inputImage, err := loadAnimationGIF(input)
	if err != nil {
		return errors.Wrap(err, "error in loadImage")
	}

	overlayImage, err := loadImage(overlay)
	if err != nil {
		return errors.Wrap(err, "error in loadImage")
	}

	DebugTime("processOverlay start")
	outputGIF, err := processOverlay(inputImage, overlayImage, width, height)
	if err != nil {
		return errors.Wrap(err, "error in processOverlay")
	}
	DebugTime("processOverlay end")

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

func processOverlay(inputImage *gif.GIF, overlayImage *image.Image, width, height uint) (*gif.GIF, error) {
	outputGIF := gif.GIF{
		Image:           make([]*image.Paletted, len(inputImage.Image)),
		Delay:           inputImage.Delay,
		Disposal:        inputImage.Disposal,
		BackgroundIndex: inputImage.BackgroundIndex,
	}

	firstFrameBounds := inputImage.Image[0].Bounds()
	firstFrameRectangle := image.Rect(0, 0, firstFrameBounds.Dx(), firstFrameBounds.Dy())
	frameImage := image.NewRGBA(firstFrameRectangle)

	if width == 0 {
		width = uint(firstFrameBounds.Dx())
		if height == 0 {
			height = uint(firstFrameBounds.Dy())
		}
	}

	resizedOverlay := resizeImage(*overlayImage, width, height)

	for i, frame := range inputImage.Image {
		if i == 0 {
			firstFrameBounds = frame.Bounds()
		}
		bounds := frame.Bounds()
		draw.Draw(frameImage, bounds, frame, bounds.Min, draw.Over)

		sticker := image.NewRGBA(firstFrameRectangle)
		draw.Draw(sticker, bounds, frameImage, bounds.Min, draw.Src)
		draw.Draw(sticker, bounds, resizedOverlay, bounds.Min, draw.Over)

		var img image.Image
		if width != uint(firstFrameBounds.Dx()) {
			img = resizeImage(sticker, width, height)
		} else {
			img = sticker
		}
		outputGIF.Image[i] = imageToPaletted(img)
	}

	return &outputGIF, nil
}

func resizeImage(img image.Image, width, height uint) image.Image {
	return resize.Resize(width, height, img, resize.Bilinear)
}

func imageToPaletted(img image.Image) *image.Paletted {
	bounds := img.Bounds()
	q := median.Quantizer(256)
	palette := q.Quantize(make(color.Palette, 0, 256), img)
	paletted := image.NewPaletted(bounds, palette)
	draw.FloydSteinberg.Draw(paletted, bounds, img, image.ZP)
	return paletted
}
