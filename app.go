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
	"errors"
	"math"

	samure "github.com/PucklaJ/samurai-render-go"
)

const (
	StateNone            = iota
	StateDragNormal      = iota
	StateAlter           = iota
	StateDragTopLeft     = iota
	StateDragTop         = iota
	StateDragTopRight    = iota
	StateDragRight       = iota
	StateDragBottomRight = iota
	StateDragBottom      = iota
	StateDragBottomLeft  = iota
	StateDragLeft        = iota

	GrabberAnimSpeed = 1.4
)

type App struct {
	start   [2]float64 // The top left corner of the selection box
	end     [2]float64 // The bottom right corner of the selection box
	pointer [2]float64 // The raw position of the pointer in global coordinates
	anchor  [2]float64 // The position where the pointer has been released
	offset  [2]float64

	state       int
	clearScreen bool
	touchID     *int
	touchFocus  samure.Output

	grabberAnim        float64
	grabberRadius      float64
	grabberBorderWidth float64

	backgroundColor    [4]float64
	selectionColor     [4]float64
	borderColor        [4]float64
	textColor          [4]float64
	grabberColor       [4]float64
	grabberBorderColor [4]float64
	padding            float64
	aspect             float64
}

func (a App) GetSelection() (samure.Rect, error) {
	if a.start[0] == 0.0 && a.start[1] == 0.0 && a.end[0] == 0.0 && a.end[1] == 0.0 {
		return samure.Rect{}, errors.New("selection cancelled")
	}

	return samure.Rect{
		X: int(a.start[0]),
		Y: int(a.start[1]),
		W: int(a.end[0] - a.start[0]),
		H: int(a.end[1] - a.start[1]),
	}, nil
}

func (a *App) OnEvent(ctx samure.Context, event interface{}) {
	switch e := event.(type) {
	case samure.EventPointerButton:
		switch a.state {
		case StateNone:
			if e.Button == samure.ButtonLeft && e.State == samure.StatePressed {
				a.anchor = a.pointer
				a.computeStartEnd()

				a.state = StateDragNormal
				ctx.SetRenderState(samure.RenderStateOnce)
			}
		case StateDragNormal:
			if e.Button == samure.ButtonLeft && e.State == samure.StateReleased {
				if flags.AlterSelection {
					a.state = StateAlter
					ctx.SetRenderState(samure.RenderStateOnce)
				} else {
					ctx.SetRunning(false)
				}
			}
		case StateAlter:
			if e.Button == samure.ButtonLeft && e.State == samure.StatePressed {
				px := a.pointer[0]
				py := a.pointer[1]
				x := a.start[0]
				y := a.start[1]
				w := a.end[0] - a.start[0]
				h := a.end[1] - a.start[1]

				if a.pointerInGrabber(px, py, x, y) {
					a.offset[0] = x - px
					a.offset[1] = y - py
					a.state = StateDragTopLeft
				} else if a.pointerInGrabber(px, py, x+w/2.0, y) {
					a.offset[0] = x + w/2.0 - px
					a.offset[1] = y - py
					a.state = StateDragTop
				} else if a.pointerInGrabber(px, py, x+w, y) {
					a.offset[0] = x + w - px
					a.offset[1] = y - py
					a.state = StateDragTopRight
				} else if a.pointerInGrabber(px, py, x+w, y+h/2.0) {
					a.offset[0] = x + w - px
					a.offset[1] = y + h/2.0 - py
					a.state = StateDragRight
				} else if a.pointerInGrabber(px, py, x+w, y+h) {
					a.offset[0] = x + w - px
					a.offset[1] = y + h - py
					a.state = StateDragBottomRight
				} else if a.pointerInGrabber(px, py, x+w/2.0, y+h) {
					a.offset[0] = x + w/2.0 - px
					a.offset[1] = y + h - py
					a.state = StateDragBottom
				} else if a.pointerInGrabber(px, py, x, y+h) {
					a.offset[0] = x - px
					a.offset[1] = y + h - py
					a.state = StateDragBottomLeft
				} else if a.pointerInGrabber(px, py, x, y+h/2.0) {
					a.offset[0] = x - px
					a.offset[1] = y + h/2.0 - py
					a.state = StateDragLeft
				} else {
					a.grabberAnim = 0.0
					a.anchor = a.pointer
					a.computeStartEnd()
					a.state = StateDragNormal
					ctx.SetRenderState(samure.RenderStateOnce)
				}
			}
		case StateDragTopLeft:
			fallthrough
		case StateDragTop:
			fallthrough
		case StateDragTopRight:
			fallthrough
		case StateDragRight:
			fallthrough
		case StateDragBottomRight:
			fallthrough
		case StateDragBottom:
			fallthrough
		case StateDragBottomLeft:
			fallthrough
		case StateDragLeft:
			if e.Button == samure.ButtonLeft && e.State == samure.StateReleased {
				a.state = StateAlter
			}
		}
	case samure.EventTouchDown:
		if a.touchID != nil && *a.touchID != e.TouchID {
			break
		}

		a.touchID = new(int)
		*a.touchID = e.TouchID
		a.touchFocus = e.Output

		a.anchor[0] = e.X + float64(e.Output.Geo().X)
		a.anchor[1] = e.Y + float64(e.Output.Geo().Y)

		a.pointer = a.anchor
		a.start = a.anchor
		a.end = a.start

		a.state = StateDragNormal
		ctx.SetRenderState(samure.RenderStateOnce)
	case samure.EventTouchUp:
		if a.touchID == nil || *a.touchID != e.TouchID {
			break
		}

		switch a.state {
		case StateDragNormal:
			a.touchID = nil
			ctx.SetRunning(false)
		}
	case samure.EventPointerMotion:
		a.pointer[0] = e.X + float64(e.Seat.PointerFocus().Output().Geo().X)
		a.pointer[1] = e.Y + float64(e.Seat.PointerFocus().Output().Geo().Y)

		px := a.pointer[0] + a.offset[0]
		py := a.pointer[1] + a.offset[1]
		x := a.start[0]
		y := a.start[1]
		w := a.end[0] - a.start[0]
		h := a.end[1] - a.start[1]

		switch a.state {
		case StateDragNormal:
			a.computeStartEnd()
			ctx.SetRenderState(samure.RenderStateOnce)
		case StateDragTopLeft:
			w += x - px
			h += y - py
			x = px
			y = py
			a.start[0] = x
			a.start[1] = y
			a.end[0] = x + w
			a.end[1] = y + h
			ctx.SetRenderState(samure.RenderStateOnce)
		case StateDragTop:
			h += y - py
			y = py
			a.start[1] = y
			a.end[1] = y + h
			ctx.SetRenderState(samure.RenderStateOnce)
		case StateDragTopRight:
			w = px - x
			h += y - py
			y = py
			a.start[1] = y
			a.end[0] = x + w
			a.end[1] = y + h
			ctx.SetRenderState(samure.RenderStateOnce)
		case StateDragRight:
			w = px - x
			a.end[0] = x + w
			ctx.SetRenderState(samure.RenderStateOnce)
		case StateDragBottomRight:
			w = px - x
			h = py - y
			a.end[0] = x + w
			a.end[1] = y + h
			ctx.SetRenderState(samure.RenderStateOnce)
		case StateDragBottom:
			h = py - y
			a.end[1] = y + h
			ctx.SetRenderState(samure.RenderStateOnce)
		case StateDragBottomLeft:
			w += x - px
			x = px
			h = py - y
			a.start[0] = x
			a.end[0] = x + w
			a.end[1] = y + h
			ctx.SetRenderState(samure.RenderStateOnce)
		case StateDragLeft:
			w += x - px
			x = px
			a.start[0] = x
			a.end[0] = x + w
			ctx.SetRenderState(samure.RenderStateOnce)
		}
	case samure.EventTouchMotion:
		if a.touchID == nil || *a.touchID != e.TouchID {
			break
		}

		a.pointer[0] = e.X + float64(a.touchFocus.Geo().X)
		a.pointer[1] = e.Y + float64(a.touchFocus.Geo().Y)

		switch a.state {
		case StateDragNormal:
			a.computeStartEnd()
			ctx.SetRenderState(samure.RenderStateOnce)
		}
	case samure.EventPointerEnter:
		e.Seat.SetPointerShape(samure.CursorShapeCrosshair)
	case samure.EventKeyboardKey:
		if e.Key == samure.KeyEsc && e.State == samure.StateReleased {
			a.start = [2]float64{0.0, 0.0}
			a.end = a.start
			ctx.SetRunning(false)
			break
		}

		switch a.state {
		case StateAlter:
			if e.Key == samure.KeyEnter && e.State == samure.StateReleased {
				ctx.SetRunning(false)
			}
		}

	}
}

func (a *App) OnUpdate(ctx samure.Context, deltaTime float64) {
	if a.state < StateAlter || a.state > StateDragLeft {
		return
	}

	if a.grabberAnim < 1.0 {
		a.grabberAnim = math.Min(a.grabberAnim+GrabberAnimSpeed*deltaTime, 1.0)
		a.grabberRadius = easeOutElastic(a.grabberAnim) * flags.GrabberRadius
		a.grabberBorderWidth = easeOutElastic(a.grabberAnim) * flags.BorderWidth
		ctx.SetRenderState(samure.RenderStateOnce)
	}
}

func (a *App) computeStartEnd() {
	width := math.Abs(a.pointer[0] - a.anchor[0])
	height := math.Abs(a.pointer[1] - a.anchor[1])

	if a.aspect != 0.0 {
		width = math.Max(width, height*a.aspect)
		height = math.Max(height, width/a.aspect)
	}

	if a.pointer[0] < a.anchor[0] {
		a.start[0] = a.anchor[0] - width - 1
		a.end[0] = a.anchor[0]
	} else {
		a.start[0] = a.anchor[0]
		a.end[0] = a.anchor[0] + width + 1
	}
	if a.pointer[1] < a.anchor[1] {
		a.start[1] = a.anchor[1] - height - 1
		a.end[1] = a.anchor[1]
	} else {
		a.start[1] = a.anchor[1]
		a.end[1] = a.anchor[1] + height + 1
	}
}

func (a App) pointerInGrabber(x, y, gx, gy float64) bool {
	dx := gx - x
	dy := gy - y
	r := a.grabberRadius + a.grabberBorderWidth/2.0
	return (dx*dx + dy*dy) < r*r
}
