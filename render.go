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
	"fmt"

	samure "github.com/PucklaJ/samurai-render-go"
	samure_cairo "github.com/PucklaJ/samurai-render-go/backends/cairo"
	"github.com/gotk3/gotk3/cairo"
)

func (a *App) OnRender(ctx samure.Context, layerSurface samure.LayerSurface, o samure.Rect, deltaTime float64) {
	c := samure_cairo.Get(layerSurface)
	c.SetOperator(cairo.OPERATOR_SOURCE)

	// Clear the screen with the background color
	c.SetSourceRGBA(
		a.backgroundColor[0],
		a.backgroundColor[1],
		a.backgroundColor[2],
		a.backgroundColor[3],
	)
	c.Paint()

	if a.state == StateNone {
		return
	}

	xGlobal := a.start[0]
	yGlobal := a.start[1]
	w := a.end[0] - a.start[0]
	h := a.end[1] - a.start[1]

	if o.RectInOutput(int(xGlobal), int(yGlobal), int(w), int(h)) {
		x := o.RelX(xGlobal)
		y := o.RelY(yGlobal)

		if a.clearScreen {
			c.SetSourceRGBA(0.0, 0.0, 0.0, 0.0)
			c.Rectangle(
				x-flags.BorderWidth/2.0,
				y-flags.BorderWidth/2.0,
				w+flags.BorderWidth,
				h+flags.BorderWidth,
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
		c.Rectangle(x, y, w, h)
		c.Fill()
		// Render the border of the selection
		c.SetSourceRGBA(
			a.borderColor[0],
			a.borderColor[1],
			a.borderColor[2],
			a.borderColor[3],
		)
		c.Rectangle(x, y, w, h)
		c.SetLineWidth(flags.BorderWidth)
		c.Stroke()
	}

	if flags.Text {
		x := o.RelX(xGlobal)
		y := o.RelY(yGlobal)

		widthStr := fmt.Sprintf("%.0f", w)
		heightStr := fmt.Sprintf("%.0f", h)
		xStr := fmt.Sprintf("X: %.0f", x)
		yStr := fmt.Sprintf("Y: %.0f", y)

		c.SelectFontFace(flags.Font, cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_NORMAL)
		c.SetFontSize(flags.FontSize)
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
			x + w/2.0 - widthExt.Width/2.0,
			y + h + widthExt.Height + a.padding,
		}
		heightTextPos := [2]float64{
			x + w + a.padding,
			y + h/2.0 + heightExt.Height/2.0,
		}
		xTextPos := [2]float64{
			x,
			y - a.padding,
		}
		yTextPos := [2]float64{
			x - yExt.Width - a.padding,
			y + yExt.Height,
		}

		// Only render text if it's inside the output
		if o.RectInOutput(int(widthTextPos[0])+o.X, int(widthTextPos[1])+o.Y, int(widthExt.Width), int(widthExt.Height)) {
			c.MoveTo(widthTextPos[0], widthTextPos[1])
			c.ShowText(widthStr)
		}
		if o.RectInOutput(int(heightTextPos[0])+o.X, int(heightTextPos[1])+o.Y, int(heightExt.Width), int(heightExt.Height)) {
			c.MoveTo(heightTextPos[0], heightTextPos[1])
			c.ShowText(heightStr)
		}
		if o.RectInOutput(int(xTextPos[0])+o.X, int(xTextPos[1])+o.Y, int(xExt.Width), int(xExt.Height)) {
			c.MoveTo(xTextPos[0], xTextPos[1])
			c.ShowText(xStr)
		}
		if o.RectInOutput(int(yTextPos[0])+o.X, int(yTextPos[1])+o.Y, int(yExt.Width), int(yExt.Height)) {
			c.MoveTo(yTextPos[0], yTextPos[1])
			c.ShowText(yStr)
		}
	}
}
