package imco

import (
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
			Usage: "The input animation GIF (required)",
		},
		cli.StringFlag{
			Name:  "overlay, l",
			Usage: "The overlay image to cover the input animation GIF (required)",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "The output filename (required)",
		},
	}
	app.Action = action

	return app
}

func action(c *cli.Context) {
	input := c.String("input")
	overlay := c.String("overlay")
	output := c.String("output")

	if input == "" || overlay == "" || output == "" {
		cli.ShowAppHelp(c)
		return
	}
}
