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
	"fmt"
	"math"

	samure "github.com/Samudevv/samurai-render-go"
	samure_cairo "github.com/Samudevv/samurai-render-go/backends/cairo"
	"github.com/gotk3/gotk3/cairo"
)

func (a *App) OnRender(ctx samure.Context, layerSurface samure.LayerSurface, o samure.Rect) {
	c := samure_cairo.Get(layerSurface)
	c.SetOperator(cairo.OPERATOR_SOURCE)

	if a.state == StateChooseOutput {
		if a.selectedOutput.Handle != nil && a.selectedOutput.Geo() == o {
			if a.clearScreen {
				c.SetSourceRGBA(0.0, 0.0, 0.0, 0.0)
			} else {
				c.SetSourceRGBA(
					a.selectionColor[0],
					a.selectionColor[1],
					a.selectionColor[2],
					a.selectionColor[3],
				)
			}
		} else {
			c.SetSourceRGBA(
				a.backgroundColor[0],
				a.backgroundColor[1],
				a.backgroundColor[2],
				a.backgroundColor[3],
			)
		}
		c.Paint()
		return
	}

	// Clear the screen with the background color
	c.SetSourceRGBA(
		a.backgroundColor[0],
		a.backgroundColor[1],
		a.backgroundColor[2],
		a.backgroundColor[3],
	)
	c.Paint()

	if (a.state == StateNone ||
		(a.state == StateChooseRegion && !isRegionAnimSet(a.currentRegionAnim))) &&
		!flags.Debug {
		return
	}

	var xGlobal, yGlobal, wGlobal, hGlobal float64
	var xLocal, yLocal, wLocal, hLocal float64

	if a.state == StateChooseRegion {
		xGlobal = a.currentRegionAnim[0]
		yGlobal = a.currentRegionAnim[1]
		wGlobal = a.currentRegionAnim[2] - a.currentRegionAnim[0]
		hGlobal = a.currentRegionAnim[3] - a.currentRegionAnim[1]
	} else {
		xGlobal = a.start[0]
		yGlobal = a.start[1]
		wGlobal = a.end[0] - a.start[0]
		hGlobal = a.end[1] - a.start[1]
	}

	xLocal = o.RelX(xGlobal) * layerSurface.Scale()
	yLocal = o.RelY(yGlobal) * layerSurface.Scale()
	wLocal = wGlobal * layerSurface.Scale()
	hLocal = hGlobal * layerSurface.Scale()

	if o.RectInOutput(int(xGlobal), int(yGlobal), int(wGlobal), int(hGlobal)) {
		borderWidthLocal := flags.BorderWidth * layerSurface.Scale()

		if a.clearScreen {
			c.SetSourceRGBA(0.0, 0.0, 0.0, 0.0)
			c.Rectangle(
				xLocal-borderWidthLocal/2.0,
				yLocal-borderWidthLocal/2.0,
				wLocal+borderWidthLocal,
				hLocal+borderWidthLocal,
			)
			c.Fill()
			return
		}

		// Render the selection
		c.SetSourceRGBA(
			a.selectionColor[0],
			a.selectionColor[1],
			a.selectionColor[2],
			a.selectionColor[3],
		)
		c.Rectangle(xLocal, yLocal, wLocal, hLocal)
		c.Fill()
		// Render the border of the selection
		c.SetSourceRGBA(
			a.borderColor[0],
			a.borderColor[1],
			a.borderColor[2],
			a.borderColor[3],
		)
		c.Rectangle(xLocal, yLocal, wLocal, hLocal)
		c.SetLineWidth(borderWidthLocal)
		c.Stroke()
	}

	a.renderGrabbers(c, o, xLocal, yLocal, wLocal, hLocal, layerSurface.Scale())

	if flags.Text && wLocal >= 1.0 && hLocal >= 1.0 {
		paddingLocal := a.padding * layerSurface.Scale()
		grabberRadiusLocal := a.grabberRadius * layerSurface.Scale()
		grabberBorderWidthLocal := a.grabberBorderWidth * layerSurface.Scale()

		widthStr := fmt.Sprintf("%d", int(wGlobal))
		heightStr := fmt.Sprintf("%d", int(hGlobal))
		xStr := fmt.Sprintf("X: %d", int(xGlobal))
		yStr := fmt.Sprintf("Y: %d", int(yGlobal))

		c.SelectFontFace(flags.Font, cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_NORMAL)
		c.SetFontSize(flags.FontSize * layerSurface.Scale())
		c.SetSourceRGBA(
			a.textColor[0],
			a.textColor[1],
			a.textColor[2],
			a.textColor[3],
		)
		widthExt := c.TextExtents(widthStr)
		heightExt := c.TextExtents(heightStr)
		xExt := c.TextExtents(xStr)
		yExt := c.TextExtents(yStr)

		widthTextPos := [2]float64{
			xLocal + wLocal/2.0 - widthExt.Width/2.0,
			yLocal + hLocal + widthExt.Height + paddingLocal,
		}
		heightTextPos := [2]float64{
			xLocal + wLocal + paddingLocal,
			yLocal + hLocal/2.0 + heightExt.Height/2.0,
		}
		xTextPos := [2]float64{
			xLocal,
			yLocal - paddingLocal,
		}
		yTextPos := [2]float64{
			xLocal - yExt.Width - paddingLocal,
			yLocal + yExt.Height,
		}

		if a.state >= StateAlter && a.state <= StateDragLeft {
			widthTextPos[1] += grabberRadiusLocal + grabberBorderWidthLocal/2.0
			heightTextPos[0] += grabberRadiusLocal + grabberBorderWidthLocal/2.0
			xTextPos[1] -= grabberRadiusLocal + grabberBorderWidthLocal/2.0
			yTextPos[0] -= grabberRadiusLocal + grabberBorderWidthLocal/2.0
		}

		widthTextPosGlobal := [2]float64{
			widthTextPos[0] / layerSurface.Scale(),
			widthTextPos[1] / layerSurface.Scale(),
		}
		heightTextPosGlobal := [2]float64{
			heightTextPos[0] / layerSurface.Scale(),
			heightTextPos[1] / layerSurface.Scale(),
		}
		widthExtGlobal := [2]float64{
			widthExt.Width / layerSurface.Scale(),
			widthExt.Height / layerSurface.Scale(),
		}
		heightExtGlobal := [2]float64{
			heightExt.Width / layerSurface.Scale(),
			heightExt.Height / layerSurface.Scale(),
		}
		xTextPosGlobal := [2]float64{
			xTextPos[0] / layerSurface.Scale(),
			xTextPos[1] / layerSurface.Scale(),
		}
		yTextPosGlobal := [2]float64{
			yTextPos[0] / layerSurface.Scale(),
			yTextPos[1] / layerSurface.Scale(),
		}
		xExtGlobal := [2]float64{
			xExt.Width / layerSurface.Scale(),
			xExt.Height / layerSurface.Scale(),
		}
		yExtGlobal := [2]float64{
			yExt.Width / layerSurface.Scale(),
			yExt.Height / layerSurface.Scale(),
		}

		// Only render text if it's inside the output
		if o.RectInOutput(int(widthTextPosGlobal[0])+o.X, int(widthTextPosGlobal[1])+o.Y-int(widthExtGlobal[1]), int(widthExtGlobal[0]), int(widthExtGlobal[1])) {
			c.MoveTo(widthTextPos[0], widthTextPos[1])
			c.ShowText(widthStr)
		}
		if o.RectInOutput(int(heightTextPosGlobal[0])+o.X, int(heightTextPosGlobal[1])+o.Y-int(heightExtGlobal[1]), int(heightExtGlobal[0]), int(heightExtGlobal[1])) {
			c.MoveTo(heightTextPos[0], heightTextPos[1])
			c.ShowText(heightStr)
		}
		if o.RectInOutput(int(xTextPosGlobal[0])+o.X, int(xTextPosGlobal[1])+o.Y, int(xExtGlobal[0]), int(xExtGlobal[1])) {
			c.MoveTo(xTextPos[0], xTextPos[1])
			c.ShowText(xStr)
		}
		if o.RectInOutput(int(yTextPosGlobal[0])+o.X, int(yTextPosGlobal[1])+o.Y, int(yExtGlobal[0]), int(yExtGlobal[1])) {
			c.MoveTo(yTextPos[0], yTextPos[1])
			c.ShowText(yStr)
		}
	}

	if flags.Debug {
		var stateStr string
		switch a.state {
		case StateNone:
			stateStr = "StateNone"
		case StateDragNormal:
			stateStr = "StateDragNormal"
		case StateAlter:
			stateStr = "StateAlter"
		case StateDragTopLeft:
			stateStr = "StateDragTopLeft"
		case StateDragTop:
			stateStr = "StateDragTop"
		case StateDragTopRight:
			stateStr = "StateDragTopRight"
		case StateDragRight:
			stateStr = "StateDragRight"
		case StateDragBottomRight:
			stateStr = "StateDragBottomRight"
		case StateDragBottom:
			stateStr = "StateDragBottom"
		case StateDragBottomLeft:
			stateStr = "StateDragBottomLeft"
		case StateDragLeft:
			stateStr = "StateDragLeft"
		case StateDragMiddle:
			stateStr = "StateDragMiddle"
		case StateChooseRegion:
			stateStr = "StateChooseRegion"
		default:
			stateStr = "Invalid State"
		}

		c.SelectFontFace("sans-serif", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_NORMAL)
		c.SetFontSize(30.0 * layerSurface.Scale())
		c.SetSourceRGBA(
			1.0,
			1.0,
			0.0,
			1.0,
		)

		debugStrings := []string{
			stateStr,
			fmt.Sprintf("xGlobal=%.1f", xGlobal),
			fmt.Sprintf("yGlobal=%.1f", yGlobal),
			fmt.Sprintf("wGlobal=%.1f", wGlobal),
			fmt.Sprintf("hGlobal=%.1f", hGlobal),
			fmt.Sprintf("xLocal=%.1f", xLocal),
			fmt.Sprintf("yLocal=%.1f", yLocal),
			fmt.Sprintf("wLocal=%.1f", wLocal),
			fmt.Sprintf("hLocal=%.1f", hLocal),
			fmt.Sprintf("App.regionAnim=%.1f", a.regionAnim),
		}

		var yPos float64 = 15.0

		for _, s := range debugStrings {
			ext := c.TextExtents(stateStr)
			yPos += ext.Height + 5.0
			c.MoveTo(20, yPos)
			c.ShowText(s)
		}
	}
}

func (a App) renderGrabbers(c *cairo.Context, o samure.Rect, x, y, w, h, scale float64) {
	if a.state < StateAlter || a.state > StateDragLeft {
		return
	}

	a.renderGrabber(c, x, y, o, scale)         // Top Left
	a.renderGrabber(c, x+w/2.0, y, o, scale)   // Top
	a.renderGrabber(c, x+w, y, o, scale)       // Top Right
	a.renderGrabber(c, x+w, y+h/2.0, o, scale) // Right
	a.renderGrabber(c, x+w, y+h, o, scale)     // Bottom Right
	a.renderGrabber(c, x+w/2.0, y+h, o, scale) // Bottom
	a.renderGrabber(c, x, y+h, o, scale)       // Bottom Left
	a.renderGrabber(c, x, y+h/2.0, o, scale)   // Left
}

func (a App) renderGrabber(c *cairo.Context, x, y float64, o samure.Rect, scale float64) {
	grabberRadiusLocal := a.grabberRadius * scale
	grabberBorderWidthLocal := a.grabberBorderWidth * scale

	if !o.CircleInOutput(int(x/scale)+o.X, int(y/scale)+o.Y, int(flags.GrabberRadius+flags.BorderWidth/2.0)) {
		return
	}

	c.SetSourceRGBA(
		a.grabberColor[0],
		a.grabberColor[1],
		a.grabberColor[2],
		a.grabberColor[3],
	)
	c.Arc(x, y, grabberRadiusLocal, 0.0, math.Pi*2)
	c.Fill()
	c.SetSourceRGBA(
		a.grabberBorderColor[0],
		a.grabberBorderColor[1],
		a.grabberBorderColor[2],
		a.grabberBorderColor[3],
	)
	c.Arc(x, y, grabberRadiusLocal, 0.0, math.Pi*2)
	c.SetLineWidth(grabberBorderWidthLocal)
	c.Stroke()
}

func easeOutElastic(x float64) float64 {
	c4 := (2 * math.Pi) / 3

	if x == 0.0 || x == 1.0 {
		return x
	}

	return math.Pow(2.0, -10.0*x)*math.Sin((x*10.0-0.75)*c4) + 1.0
}

func easeOutQuint(x float64) float64 {
	return 1.0 - math.Pow(1.0-x, 5.0)
}
