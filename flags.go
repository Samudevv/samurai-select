package main

import (
	"fmt"
	"os"

	flag "github.com/jessevdk/go-flags"
	css "github.com/mazznoer/csscolorparser"
)

var flags struct {
	BackgroundColor string `long:"background-color" description:"Set the clear color that fills the screen" default:"#FFFFFF40"`
	SelectionColor  string `long:"selection-color" description:"Set the color that is used to draw the inside of the selection box" default:"#00000000"`
	BorderColor     string `long:"border-color" description:"Set the color that is used to draw the border around the selection box" default:"#000000FF"`

	BorderWidth float64 `long:"border-width" description:"The width of the border in pixels" default:"2.0"`
}

func CreateApp(argv []string) (*App, error) {
	parser := flag.NewParser(&flags, flag.HelpFlag)
	argv, err := parser.ParseArgs(argv)
	if err != nil {
		if !flag.WroteHelp(err) {
			fmt.Fprintf(os.Stderr, "Arguments: %v\n", err)
			parser.WriteHelp(os.Stderr)
		} else {
			fmt.Fprint(os.Stderr, err)
		}
		os.Exit(1)
	}

	a := &App{}
	a.backgroundColor = parseColor(flags.BackgroundColor)
	a.selectionColor = parseColor(flags.SelectionColor)
	a.borderColor = parseColor(flags.BorderColor)
	if flags.BorderWidth < 0.0 {
		fmt.Fprintf(os.Stderr, "--border-width values below zero are invalid\n")
		a.borderWidth = 0.0
	} else {
		a.borderWidth = flags.BorderWidth
	}
	return a, nil
}

func parseColor(colorString string) [4]float64 {
	c, err := css.Parse(colorString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse color \"%s\": %v\n", colorString, err)
		return [4]float64{}
	}
	return [4]float64{
		c.R,
		c.G,
		c.B,
		c.A,
	}
}
