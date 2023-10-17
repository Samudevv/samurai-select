/***********************************************************************************
 *                         This file is part of samurai-select
 *                    https://github.com/PucklaJ/samurai-select
 ***********************************************************************************
 * Copyright (c) 2023 Jonas Pucher
 *
 * This software is provided ‘as-is’, without any express or implied
 * warranty. In no event will the authors be held liable for any damages
 * arising from the use of this software.
 *
 * Permission is granted to anyone to use this software for any purpose,
 * including commercial applications, and to alter it and redistribute it
 * freely, subject to the following restrictions:
 *
 * 1. The origin of this software must not be misrepresented; you must not
 * claim that you wrote the original software. If you use this software
 * in a product, an acknowledgment in the product documentation would be
 * appreciated but is not required.
 *
 * 2. Altered source versions must be plainly marked as such, and must not be
 * misrepresented as being the original software.
 *
 * 3. This notice may not be removed or altered from any source
 * distribution.
 ************************************************************************************/

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	flag "github.com/jessevdk/go-flags"
	css "github.com/mazznoer/csscolorparser"
)

var flags struct {
	BackgroundColor    string `long:"background-color" description:"Set the clear color that fills the screen" default:"#FFFFFF40"`
	SelectionColor     string `long:"selection-color" description:"Set the color that is used to draw the inside of the selection box" default:"#00000000"`
	BorderColor        string `long:"border-color" description:"Set the color that is used to draw the border around the selection box" default:"#000000FF"`
	TextColor          string `long:"text-color" description:"Set the color that is used for the text" default:"#000000FF"`
	GrabberColor       string `long:"grabber-color" description:"The fill color of the grabbers for altering the selection" default:"#101010FF"`
	GrabberBorderColor string `long:"grabber-border-color" description:"The border color of the grabbers for altering the selection" default:"#000000FF"`

	BorderWidth      float64 `long:"border-width" description:"The width of the border in pixels" default:"2.0"`
	Text             bool    `short:"t" long:"text" description:"Display the selection position and dimensions next to the selection box"`
	Font             string  `long:"font" description:"Set the font family of the text" default:"sans-serif"`
	ListFonts        bool    `long:"list-fonts" description:"List installed fonts that can be used"`
	FontSize         float64 `long:"font-size" description:"Set the font size of the text" default:"16"`
	TextPadding      float64 `long:"text-padding" description:"The distance between the selection box and each text" default:"10"`
	FreezeScreen     bool    `short:"z" long:"freeze" description:"Freeze the screen while performing the selection"`
	Screenshot       bool    `short:"s" long:"screenshot" description:"Use grim to perform a screenshot"`
	ScreenshotOutput string  `short:"o" long:"output" description:"File path where the screenshot will be stored. You can Use Explicit argument indexes (https://pkg.go.dev/fmt) where 1 is year, 2 is month, 3 is day 4 is hour 5 is minute and 6 is second" default:"screenshot-%[1]d.%02[2]d.%02[3]d-%02[4]d:%02[5]d:%02[6]d.png"`
	ScreenshotFlags  string  `long:"screenshot-flags" description:"These flags are passed to grim when performing the screenshot"`
	Command          string  `short:"c" long:"cmd" description:"Clear the screen and execute a command. This is useful to perform an action while the screen is frozen. Insert %geometry% where you want to put the resulting geometry."`
	Format           string  `short:"f" long:"format" description:"Set the format in which the geometry is output. Use Explicit argument indexes (https://pkg.go.dev/fmt) where 1 is x, 2 is y, 3 is width and 4 is height" default:"%[1]d,%[2]d %[3]dx%[4]d"`
	ForceAspectRatio string  `short:"a" long:"aspect-ratio" description:"Force an aspect ratio for the selection box in the format w:h"`
	AlterSelection   bool    `short:"A" long:"alter-selection" description:"This flag lets you change the selection box after releasing left click by dragging the box at the edges and corners"`
	GrabberRadius    float64 `long:"grabber-radius" description:"The radius of the grabbers for altering the selection" default:"7"`
	Debug            bool    `short:"d" long:"debug" description:"Show developer debug stuff"`
	NoAnimation      bool    `long:"no-anim" description:"Disable the bouncing animation of the grabbers if alter selection is enabled"`
	Regions          string  `short:"r" long:"regions" description:"Choose from predefined regions (e.g. windows) on the screen." optional:"1" optional-value:"auto" default:"none" choice:"none" choice:"auto" choice:"hyprland"`
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
	a.grabberColor = parseColor(flags.GrabberColor)
	a.grabberBorderColor = parseColor(flags.GrabberBorderColor)
	if flags.BorderWidth < 0.0 {
		fmt.Fprintf(os.Stderr, "--border-width values below zero are invalid\n")
		flags.BorderWidth = 0.0
	}
	a.padding = flags.TextPadding + flags.BorderWidth/2.0

	if flags.ForceAspectRatio != "" {
		// Parse aspect ratio
		words := strings.Split(flags.ForceAspectRatio, ":")
		for {
			if len(words) != 2 {
				fmt.Fprintln(os.Stderr, "Invalid aspect ratio")
				break
			}

			w, err := strconv.ParseInt(words[0], 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid aspect ratio: %v\n", err)
				break
			}
			h, err := strconv.ParseInt(words[1], 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid aspect ratio: %v\n", err)
				break
			}

			a.aspect = float64(w) / float64(h)
			break
		}
	}

	switch flags.Regions {
	case "none":
	case "auto":
		a.regionsObj = DetectRegions()
		if a.regionsObj == nil {
			fmt.Fprintf(os.Stderr, "Could not detect which compositor is running\n")
		}
	case "hyprland":
		a.regionsObj = &HyprlandRegions{}
	default:
		fmt.Fprintf(os.Stderr, "Invalid regions: \"%s\"\n", flags.Regions)
	}

	if a.regionsObj != nil {
		a.state = StateChooseRegion
		a.regions = a.regionsObj.OutputRegions()
		x, y, err := a.regionsObj.CursorPos()
		if err == nil {
			for i := range a.regions {
				if a.regions[i].PointInOutput(x, y) {
					a.chosenRegion = a.regions[i]
				}
			}

			if isRegionSet(a.chosenRegion) {
				a.currentRegionAnim[0] = float64(a.chosenRegion.X)
				a.currentRegionAnim[1] = float64(a.chosenRegion.Y)
				a.currentRegionAnim[2] = float64(a.chosenRegion.X + a.chosenRegion.W)
				a.currentRegionAnim[3] = float64(a.chosenRegion.Y + a.chosenRegion.H)
			}
		}
	}

	a.regionAnim = 1.0

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
