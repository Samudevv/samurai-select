package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	flag "github.com/jessevdk/go-flags"
	css "github.com/mazznoer/csscolorparser"
)

var flags struct {
	BackgroundColor string `long:"background-color" description:"Set the clear color that fills the screen" default:"#FFFFFF40"`
	SelectionColor  string `long:"selection-color" description:"Set the color that is used to draw the inside of the selection box" default:"#00000000"`
	BorderColor     string `long:"border-color" description:"Set the color that is used to draw the border around the selection box" default:"#000000FF"`
	TextColor       string `long:"text-color" description:"Set the color that is used for the text" default:"#000000FF"`

	BorderWidth      float64 `long:"border-width" description:"The width of the border in pixels" default:"2.0"`
	Text             bool    `short:"t" long:"text" description:"Display the selection position and dimensions next to the selection box"`
	Font             string  `long:"font" description:"Set the font family of the text" default:"sans-serif"`
	ListFonts        bool    `long:"list-fonts" description:"List installed fonts that can be used"`
	FontSize         float64 `long:"font-size" description:"Set the font size of the text" default:"16"`
	TextPadding      float64 `long:"text-padding" description:"The distance between the selection box and each text" default:"10"`
	FreezeScreen     bool    `short:"z" long:"freeze" description:"Freeze the screen while performing the selection"`
	Screenshot       bool    `short:"s" long:"screenshot" description:"Use grim to perform a screenshot"`
	ScreenshotOutput string  `short:"o" long:"output" description:"File path where the screenshot will be stored" default:"screenshot.png"`
	ScreenshotFlags  string  `long:"screenshot-flags" description:"These flags are passed to grim when performing the screenshot"`
	Command          string  `short:"c" long:"cmd" description:"Clear the screen and execute a command. This is useful to perform an action while the screen is frozen. Insert %geometry% where you want to put the resulting geometry."`
	Format           string  `short:"f" long:"format" description:"Set the format in which the geometry is output. Use Explicit argument indexes (https://pkg.go.dev/fmt) where 1 is x, 2 is y, 3 is width and 4 is height" default:"%[1]d,%[2]d %[3]dx%[4]d"`
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

	if flags.ListFonts {
		fcList, err := exec.LookPath("fc-list")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can not list fonts: %v\n", err)
			os.Exit(1)
		}

		var fcListOut strings.Builder
		fcListCmd := exec.Command(fcList, "--brief")
		fcListCmd.Stdout = &fcListOut
		fcListCmd.Stderr = os.Stderr

		if err := fcListCmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to list fonts")
			os.Exit(1)
		}

		scanner := bufio.NewScanner(strings.NewReader(fcListOut.String()))

		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "family:") {
				words := strings.Split(line, "\"")
				if len(words) < 2 {
					continue
				}
				familyName := strings.ToLower(words[1])
				fmt.Print("\"", familyName, "\"\n")
			}
		}

		os.Exit(0)
	}

	a := &App{}
	a.backgroundColor = parseColor(flags.BackgroundColor)
	a.selectionColor = parseColor(flags.SelectionColor)
	a.borderColor = parseColor(flags.BorderColor)
	a.textColor = parseColor(flags.TextColor)
	if flags.BorderWidth < 0.0 {
		fmt.Fprintf(os.Stderr, "--border-width values below zero are invalid\n")
		flags.BorderWidth = 0.0
	}
	a.padding = flags.TextPadding + flags.BorderWidth/2.0

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
