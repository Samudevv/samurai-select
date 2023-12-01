/***********************************************************************************
 *                         This file is part of samurai-select
 *                    https://github.com/Samudevv/samurai-select
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
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	samure "github.com/Samudevv/samurai-render-go"
)

const (
	StateNone            = iota
	StateDragNormal      = iota
	StateAlter           = iota
	StateDragMiddle      = iota
	StateDragTopLeft     = iota
	StateDragTop         = iota
	StateDragTopRight    = iota
	StateDragRight       = iota
	StateDragBottomRight = iota
	StateDragBottom      = iota
	StateDragBottomLeft  = iota
	StateDragLeft        = iota
	StateChooseRegion    = iota
	StateChooseOutput    = iota

	GrabberAnimSpeed = 1.4
	RegionAnimSpeed  = 2.5
)

type App struct {
	start          [2]float64 // The top left corner of the selection box
	end            [2]float64 // The bottom right corner of the selection box
	pointer        [2]float64 // The raw position of the pointer in global coordinates
	anchor         [2]float64 // The position where the pointer has been released
	offset         [2]float64
	selectedOutput samure.Output

	state       int
	clearScreen bool
	touchID     *int
	touchFocus  samure.Output
	cancelled   bool

	grabberAnim        float64
	grabberRadius      float64
	grabberBorderWidth float64

	selectedRegion    Region
	regionAnim        float64
	currentRegionAnim [4]float64
	startRegionAnim   [4]float64
	endRegionAnim     [4]float64

	backgroundColor    [4]float64
	selectionColor     [4]float64
	borderColor        [4]float64
	textColor          [4]float64
	grabberColor       [4]float64
	grabberBorderColor [4]float64
	padding            float64
	aspect             float64
	regionsObj         Regions
	regions            []Region
}

func (a App) GetSelection() (samure.Rect, error) {
	if a.cancelled {
		return samure.Rect{}, errors.New("selection cancelled")
	}

	switch a.state {
	case StateChooseRegion:
		return a.selectedRegion.Geo, nil
	case StateChooseOutput:
		return a.selectedOutput.Geo(), nil
	default:
		return samure.Rect{
			X: int(a.start[0]),
			Y: int(a.start[1]),
			W: int(a.end[0] - a.start[0]),
			H: int(a.end[1] - a.start[1]),
		}, nil
	}
}

func (a *App) OnUpdate(ctx samure.Context, deltaTime float64) {
	if a.state >= StateAlter && a.state <= StateDragLeft {
		if a.grabberAnim < 1.0 {
			if flags.NoAnimation {
				a.grabberAnim = 1.0
			} else {
				a.grabberAnim = math.Min(a.grabberAnim+GrabberAnimSpeed*deltaTime, 1.0)
			}

			a.grabberRadius = easeOutElastic(a.grabberAnim) * flags.GrabberRadius
			a.grabberBorderWidth = easeOutElastic(a.grabberAnim) * flags.BorderWidth
			ctx.SetRenderState(samure.RenderStateOnce)
		}
	} else if a.state == StateChooseRegion {
		if a.regionAnim < 1.0 {
			if flags.NoAnimation {
				a.regionAnim = 1.0
			} else {
				a.regionAnim = math.Min(a.regionAnim+RegionAnimSpeed*deltaTime, 1.0)
			}

			for i := 0; i < 4; i++ {
				a.currentRegionAnim[i] = a.startRegionAnim[i] + (a.endRegionAnim[i]-a.startRegionAnim[i])*easeOutQuint(a.regionAnim)
			}

			ctx.SetRenderState(samure.RenderStateOnce)
		}

		if !flags.FreezeScreen && a.regionsObj != nil {
			a.regions = a.regionsObj.OutputRegions()
			a.pointerMove(ctx, a.pointer[0], a.pointer[1], 0.0, 0.0, a.selectedOutput)
		}
	}
}

func (a App) createOutputString() (string, error) {
	// Retrieve data that will be output using the format
	sel, err := a.GetSelection()
	if err != nil {
		return "", err
	}

	outputName := "nil"
	if a.selectedOutput.Handle != nil {
		outputName = a.selectedOutput.Name()
	}

	regionName := "nil"
	if isRegionSet(a.selectedRegion.Geo) {
		regionName = a.selectedRegion.Name
	}

	var outputRelX, outputRelY, outputRelW, outputRelH int
	if a.selectedOutput.Handle != nil {
		// The corners of the selection relative to the selected output
		relX := a.selectedOutput.Geo().RelX(float64(sel.X))
		relY := a.selectedOutput.Geo().RelY(float64(sel.Y))
		relEndX := a.selectedOutput.Geo().RelX(float64(sel.X + sel.W))
		relEndY := a.selectedOutput.Geo().RelY(float64(sel.Y + sel.H))

		// Clamp the values above to the geometry of the output
		outputX := math.Min(math.Max(0.0, relX), float64(a.selectedOutput.Geo().W))
		outputY := math.Min(math.Max(0.0, relY), float64(a.selectedOutput.Geo().H))
		outputEndX := math.Min(math.Max(0.0, relEndX), float64(a.selectedOutput.Geo().W))
		outputEndY := math.Min(math.Max(0.0, relEndY), float64(a.selectedOutput.Geo().H))

		// Convert them to int and calculate width and height
		outputRelX = int(outputX)
		outputRelY = int(outputY)
		outputRelW = int(outputEndX - outputX)
		outputRelH = int(outputEndY - outputY)
	}

	var out strings.Builder
	var parseSpecifier bool
	for _, r := range flags.Format {
		if parseSpecifier {
			switch r {
			case 'x':
				out.WriteString(strconv.Itoa(sel.X))
			case 'y':
				out.WriteString(strconv.Itoa(sel.Y))
			case 'w':
				out.WriteString(strconv.Itoa(sel.W))
			case 'h':
				out.WriteString(strconv.Itoa(sel.H))
			case 'X':
				out.WriteString(strconv.Itoa(outputRelX))
			case 'Y':
				out.WriteString(strconv.Itoa(outputRelY))
			case 'W':
				out.WriteString(strconv.Itoa(outputRelW))
			case 'H':
				out.WriteString(strconv.Itoa(outputRelH))
			case 'r':
				out.WriteString(regionName)
			case 'o':
				out.WriteString(outputName)
			case '%':
				out.WriteRune(r)
			default:
				return "", fmt.Errorf("invalid format specifier: \"%s\"", string(r))
			}
			parseSpecifier = false
		} else {
			switch r {
			case '%':
				parseSpecifier = true
			default:
				out.WriteRune(r)
			}
		}
	}

	return out.String(), nil
}

func isRegionSet(r samure.Rect) bool {
	return r.X != 0 || r.Y != 0 || r.W != 0 || r.H != 0
}

func isRegionAnimSet(r [4]float64) bool {
	return r[0] != 0 || r[1] != 0 || r[2] != 0 || r[3] != 0
}

func unsetRegion(r *samure.Rect) {
	r.X = 0
	r.Y = 0
	r.W = 0
	r.H = 0
}

func createScreenshotFilename(t time.Time) (string, error) {
	var out strings.Builder
	var parseSpecifier bool
	for _, r := range flags.ScreenshotOutput {
		if parseSpecifier {
			switch r {
			case 'n':
				out.WriteString(strconv.Itoa(t.Nanosecond()))
			case 's':
				out.WriteString(fmt.Sprintf("%02d", t.Second()))
			case 'm':
				out.WriteString(fmt.Sprintf("%02d", t.Minute()))
			case 'h':
				out.WriteString(fmt.Sprintf("%02d", t.Hour()))
			case 'd':
				out.WriteString(fmt.Sprintf("%02d", t.Day()))
			case 'M':
				out.WriteString(fmt.Sprintf("%02d", t.Month()))
			case 'o':
				out.WriteString(t.Month().String())
			case 'y':
				out.WriteString(strconv.Itoa(t.Year()))
			case '%':
				out.WriteRune(r)
			default:
				return "", fmt.Errorf("invalid format specifier: \"%s\"", string(r))
			}
			parseSpecifier = false
		} else {
			switch r {
			case '%':
				parseSpecifier = true
			default:
				out.WriteRune(r)
			}
		}
	}

	return out.String(), nil
}
