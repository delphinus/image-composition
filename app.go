package imco

import (
	"fmt"
	"os"

	"github.com/Songmu/prompter"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

// NewApp will return imco App
func NewApp() *cli.App {
	app := cli.NewApp()
	app.Usage = "Image composition sample"
	app.Version = Version
	app.Author = "delphinus"
	app.Email = "delphinus@remora.cx"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "input, i",
			Usage: "Input animation GIF (required)",
		},
		cli.StringFlag{
			Name:  "overlay, l",
			Usage: "Overlay image to cover the input animation GIF (required)",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "Output filename",
			Value: DefaultOutput,
		},
		cli.IntFlag{
			Name:  "width, w",
			Usage: "Width of the output file (default: same as input file)",
		},
		cli.IntFlag{
			Name:  "height, h",
			Usage: "Height of the output file (default: same as input file)",
		},
	}
	app.Action = action

	return app
}

func action(c *cli.Context) error {
	input := c.String("input")
	overlay := c.String("overlay")
	output := c.String("output")
	width := c.Int("width")
	height := c.Int("height")

	if input == "" || overlay == "" {
		cli.ShowAppHelp(c)
		return cli.NewExitError("input & overlay both needed", InitializationError)
	}

	if _, err := os.Stat(input); os.IsNotExist(err) {
		return cli.NewExitError(fmt.Sprintf("input file not found: %s", input), InitializationError)
	}

	if _, err := os.Stat(overlay); os.IsNotExist(err) {
		return cli.NewExitError(fmt.Sprintf("overlay file not found: %s", overlay), InitializationError)
	}

	if _, err := os.Stat(output); os.IsExist(err) {
		if ok := prompter.YN(fmt.Sprintf("'%s' already exists. can I overwrite this?", output), false); !ok {
			return cli.NewExitError("output file already exists", InitializationError)
		}
	}

	if err := OverlayImage(input, overlay, output, width, height); err != nil {
		return cli.NewExitError(errors.Wrap(err, "error in OverlayImange"), OverlayImageError)
	}

	return nil
}
