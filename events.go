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
	"math"

	samure "github.com/PucklaJ/samurai-render-go"
)

func (a *App) pointerDown(ctx samure.Context, px, py float64, focus samure.Output) {
	switch a.state {
	case StateNone:
		a.selectedOutput = focus
		a.anchor[0], a.anchor[1] = px, py
		a.computeStartEnd(px, py, px, py)

		a.state = StateDragNormal
		ctx.SetRenderState(samure.RenderStateOnce)
	case StateAlter:
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
			sel := samure.Rect{
				X: int(x),
				Y: int(y),
				W: int(w),
				H: int(h),
			}

			if sel.PointInOutput(int(px), int(py)) {
				a.state = StateDragMiddle
			} else {
				a.selectedOutput = focus
				a.grabberAnim = 0.0
				a.anchor[0], a.anchor[1] = px, py
				a.computeStartEnd(px, py, px, py)

				a.state = StateDragNormal
				ctx.SetRenderState(samure.RenderStateOnce)
			}
		}
	case StateChooseRegion:
		a.cancelled = !isRegionSet(a.chosenRegion)
		ctx.SetRunning(false)
	case StateChooseOutput:
		ctx.SetRunning(a.selectedOutput.Handle == nil)
	}
}

func (a *App) pointerUp(ctx samure.Context) {
	switch a.state {
	case StateDragNormal:
		if flags.AlterSelection {
			a.state = StateAlter
			ctx.SetRenderState(samure.RenderStateOnce)
		} else {
			ctx.SetRunning(false)
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
		fallthrough
	case StateDragMiddle:
		a.state = StateAlter
	}
}

func (a *App) pointerMove(ctx samure.Context, px, py, dx, dy float64, focus samure.Output) {
	pox := px + a.offset[0]
	poy := py + a.offset[1]
	x := a.start[0]
	y := a.start[1]
	w := a.end[0] - a.start[0]
	h := a.end[1] - a.start[1]

	switch a.state {
	case StateDragNormal:
		a.computeStartEnd(px, py, a.anchor[0], a.anchor[1])
		ctx.SetRenderState(samure.RenderStateOnce)
	case StateDragTopLeft:
		w += x - pox
		h += y - poy
		x = pox
		y = poy
		a.start[0] = x
		a.start[1] = y
		a.end[0] = x + w
		a.end[1] = y + h
		a.selectedOutput = focus
		a.handleOverlapAndAspectRatio()
		ctx.SetRenderState(samure.RenderStateOnce)
	case StateDragTop:
		h += y - poy
		y = poy
		a.start[1] = y
		a.end[1] = y + h
		a.selectedOutput = focus
		a.handleOverlapAndAspectRatio()
		ctx.SetRenderState(samure.RenderStateOnce)
	case StateDragTopRight:
		w = pox - x
		h += y - poy
		y = poy
		a.start[1] = y
		a.end[0] = x + w
		a.end[1] = y + h
		a.selectedOutput = focus
		a.handleOverlapAndAspectRatio()
		ctx.SetRenderState(samure.RenderStateOnce)
	case StateDragRight:
		w = pox - x
		a.end[0] = x + w
		a.selectedOutput = focus
		a.handleOverlapAndAspectRatio()
		ctx.SetRenderState(samure.RenderStateOnce)
	case StateDragBottomRight:
		w = pox - x
		h = poy - y
		a.end[0] = x + w
		a.end[1] = y + h
		a.selectedOutput = focus
		a.handleOverlapAndAspectRatio()
		ctx.SetRenderState(samure.RenderStateOnce)
	case StateDragBottom:
		h = poy - y
		a.end[1] = y + h
		a.selectedOutput = focus
		a.handleOverlapAndAspectRatio()
		ctx.SetRenderState(samure.RenderStateOnce)
	case StateDragBottomLeft:
		w += x - pox
		x = pox
		h = poy - y
		a.start[0] = x
		a.end[0] = x + w
		a.end[1] = y + h
		a.selectedOutput = focus
		a.handleOverlapAndAspectRatio()
		ctx.SetRenderState(samure.RenderStateOnce)
	case StateDragLeft:
		w += x - pox
		x = pox
		a.start[0] = x
		a.end[0] = x + w
		a.selectedOutput = focus
		a.handleOverlapAndAspectRatio()
		ctx.SetRenderState(samure.RenderStateOnce)
	case StateDragMiddle:
		x += dx
		y += dy
		a.start[0] = x
		a.start[1] = y
		a.end[0] = x + w
		a.end[1] = y + h
		a.selectedOutput = focus
		ctx.SetRenderState(samure.RenderStateOnce)
	case StateChooseRegion:
		a.selectedOutput = focus
		prevRegion := a.chosenRegion
		unsetRegion(&a.chosenRegion)

		for i := range a.regions {
			if a.regions[i].PointInOutput(int(px), int(py)) {
				a.chosenRegion = a.regions[i]
				break
			}
		}

		if a.chosenRegion != prevRegion {
			if a.regionAnim < 1.0 {
				a.startRegionAnim = a.currentRegionAnim
			} else {
				if isRegionSet(prevRegion) {
					a.startRegionAnim[0] = float64(prevRegion.X)
					a.startRegionAnim[1] = float64(prevRegion.Y)
					a.startRegionAnim[2] = float64(prevRegion.X + prevRegion.W)
					a.startRegionAnim[3] = float64(prevRegion.Y + prevRegion.H)
				} else {
					a.startRegionAnim[0] = float64(a.chosenRegion.X + a.chosenRegion.W/2)
					a.startRegionAnim[1] = float64(a.chosenRegion.Y + a.chosenRegion.H/2)
					a.startRegionAnim[2] = float64(a.chosenRegion.X + a.chosenRegion.W/2)
					a.startRegionAnim[3] = float64(a.chosenRegion.Y + a.chosenRegion.H/2)
				}
			}

			if isRegionSet(a.chosenRegion) {
				a.endRegionAnim[0] = float64(a.chosenRegion.X)
				a.endRegionAnim[1] = float64(a.chosenRegion.Y)
				a.endRegionAnim[2] = float64(a.chosenRegion.X + a.chosenRegion.W)
				a.endRegionAnim[3] = float64(a.chosenRegion.Y + a.chosenRegion.H)
			} else {
				a.endRegionAnim[0] = float64(prevRegion.X + prevRegion.W/2)
				a.endRegionAnim[1] = float64(prevRegion.Y + prevRegion.H/2)
				a.endRegionAnim[2] = float64(prevRegion.X + prevRegion.W/2)
				a.endRegionAnim[3] = float64(prevRegion.Y + prevRegion.H/2)
			}

			a.regionAnim = 0.0

			ctx.SetRenderState(samure.RenderStateOnce)
		}
	case StateChooseOutput:
		prevOutput := a.selectedOutput
		a.selectedOutput = samure.Output{Handle: nil}

		for i := 0; i < ctx.LenOutputs(); i++ {
			if ctx.Output(i).PointInOutput(int(px), int(py)) {
				a.selectedOutput = ctx.Output(i)
				break
			}
		}

		if a.selectedOutput != prevOutput {
			ctx.SetRenderState(samure.RenderStateOnce)
		}
	}
}

func (a *App) OnEvent(ctx samure.Context, event interface{}) {
	switch e := event.(type) {
	case samure.EventPointerButton:
		if e.Button == samure.ButtonLeft {
			switch e.State {
			case samure.StatePressed:
				a.pointerDown(ctx, a.pointer[0], a.pointer[1], e.Seat.PointerFocus().Output())
			case samure.StateReleased:
				a.pointerUp(ctx)
			}
			ctx.SetPointerShape(a.getCursorShape())
		}
	case samure.EventTouchDown:
		if a.touchID != nil && *a.touchID != e.TouchID {
			break
		}

		a.touchID = new(int)
		*a.touchID = e.TouchID
		a.touchFocus = e.Output

		a.pointer[0] = e.X + float64(e.Output.Geo().X)
		a.pointer[1] = e.Y + float64(e.Output.Geo().Y)
		a.pointerDown(ctx, a.pointer[0], a.pointer[1], e.Output)
	case samure.EventTouchUp:
		if a.touchID == nil || *a.touchID != e.TouchID {
			break
		}
		a.touchID = nil

		a.pointerUp(ctx)
	case samure.EventPointerMotion:
		px := e.X + float64(e.Seat.PointerFocus().Output().Geo().X)
		py := e.Y + float64(e.Seat.PointerFocus().Output().Geo().Y)
		dx := px - a.pointer[0]
		dy := py - a.pointer[1]
		a.pointer[0], a.pointer[1] = px, py

		ctx.SetPointerShape(a.getCursorShape())

		a.pointerMove(ctx, px, py, dx, dy, e.Seat.PointerFocus().Output())
	case samure.EventTouchMotion:
		if a.touchID == nil || *a.touchID != e.TouchID {
			break
		}

		px := e.X + float64(a.touchFocus.Geo().X)
		py := e.Y + float64(a.touchFocus.Geo().Y)
		dx := px - a.pointer[0]
		dy := py - a.pointer[1]
		a.pointer[0], a.pointer[1] = px, py
		a.pointerMove(ctx, px, py, dx, dy, e.Seat.TouchFocus().Output())
	case samure.EventPointerEnter:
		switch a.state {
		case StateNone:
			ctx.SetPointerShape(samure.CursorShapeCrosshair)
		}
	case samure.EventKeyboardKey:
		if e.Key == samure.KeyEsc && e.State == samure.StateReleased {
			a.cancelled = true
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

func (a *App) computeStartEnd(px, py, ax, ay float64) {
	width := math.Abs(px - ax)
	height := math.Abs(py - ay)

	if a.aspect != 0.0 {
		width = math.Max(width, height*a.aspect)
		height = math.Max(height, width/a.aspect)
	}

	if px < ax {
		a.start[0] = ax - width - 1
		a.end[0] = ax
	} else {
		a.start[0] = ax
		a.end[0] = ax + width + 1
	}
	if py < ay {
		a.start[1] = ay - height - 1
		a.end[1] = ay
	} else {
		a.start[1] = ay
		a.end[1] = ay + height + 1
	}
}

func (a App) pointerInGrabber(x, y, gx, gy float64) bool {
	dx := gx - x
	dy := gy - y
	r := a.grabberRadius + a.grabberBorderWidth/2.0
	return (dx*dx + dy*dy) < r*r
}

func (a *App) handleOverlapAndAspectRatio() {
	x := a.start[0]
	y := a.start[1]
	w := a.end[0] - a.start[0]
	h := a.end[1] - a.start[1]

	if w < 0 {
		x += w
		w = -w
		a.start[0] = x
		a.end[0] = x + w

		switch a.state {
		case StateDragTopLeft:
			a.state = StateDragTopRight
		case StateDragTopRight:
			a.state = StateDragTopLeft
		case StateDragBottomLeft:
			a.state = StateDragBottomRight
		case StateDragBottomRight:
			a.state = StateDragBottomLeft
		case StateDragLeft:
			a.state = StateDragRight
		case StateDragRight:
			a.state = StateDragLeft
		}

	}

	if h < 0 {
		y += h
		h = -h
		a.start[1] = y
		a.end[1] = y + h

		switch a.state {
		case StateDragTopLeft:
			a.state = StateDragBottomLeft
		case StateDragTopRight:
			a.state = StateDragBottomRight
		case StateDragBottomLeft:
			a.state = StateDragTopLeft
		case StateDragBottomRight:
			a.state = StateDragTopRight
		case StateDragTop:
			a.state = StateDragBottom
		case StateDragBottom:
			a.state = StateDragTop
		}
	}

	if a.aspect != 0.0 {
		x = a.start[0]
		y = a.start[1]
		w = a.end[0] - a.start[0]
		h = a.end[1] - a.start[1]

		width := math.Max(w, h*a.aspect)
		height := math.Max(h, w/a.aspect)

		switch a.state {
		case StateDragTopLeft:
			x -= width - w
			y -= height - h
		case StateDragTopRight:
			y -= height - h
		case StateDragBottomLeft:
			x -= width - w
		case StateDragTop:
			width = h * a.aspect
			height = h
		case StateDragBottom:
			width = h * a.aspect
			height = h
		case StateDragLeft:
			width = w
			height = w / a.aspect
		case StateDragRight:
			width = w
			height = w / a.aspect
		}

		a.start[0] = x
		a.start[1] = y
		a.end[0] = x + width
		a.end[1] = y + height
	}
}

func (a *App) getCursorShape() int {
	switch a.state {
	case StateNone, StateDragNormal:
		return samure.CursorShapeCrosshair
	case StateAlter:
		px := a.pointer[0]
		py := a.pointer[1]
		x := a.start[0]
		y := a.start[1]
		w := a.end[0] - a.start[0]
		h := a.end[1] - a.start[1]

		if a.pointerInGrabber(px, py, x, y) {
			return samure.CursorShapeNwResize
		} else if a.pointerInGrabber(px, py, x+w/2.0, y) {
			return samure.CursorShapeNResize
		} else if a.pointerInGrabber(px, py, x+w, y) {
			return samure.CursorShapeNeResize
		} else if a.pointerInGrabber(px, py, x+w, y+h/2.0) {
			return samure.CursorShapeEResize
		} else if a.pointerInGrabber(px, py, x+w, y+h) {
			return samure.CursorShapeSeResize
		} else if a.pointerInGrabber(px, py, x+w/2.0, y+h) {
			return samure.CursorShapeSResize
		} else if a.pointerInGrabber(px, py, x, y+h) {
			return samure.CursorShapeSwResize
		} else if a.pointerInGrabber(px, py, x, y+h/2.0) {
			return samure.CursorShapeWResize
		} else {
			sel := samure.Rect{
				X: int(x),
				Y: int(y),
				W: int(w),
				H: int(h),
			}

			if sel.PointInOutput(int(px), int(py)) {
				return samure.CursorShapeGrab
			} else {
				return samure.CursorShapeCrosshair
			}
		}
	case StateDragMiddle:
		return samure.CursorShapeGrabbing
	case StateDragTopLeft:
		return samure.CursorShapeNwResize
	case StateDragTop:
		return samure.CursorShapeNResize
	case StateDragTopRight:
		return samure.CursorShapeNeResize
	case StateDragRight:
		return samure.CursorShapeEResize
	case StateDragBottomRight:
		return samure.CursorShapeSeResize
	case StateDragBottom:
		return samure.CursorShapeSResize
	case StateDragBottomLeft:
		return samure.CursorShapeSwResize
	case StateDragLeft:
		return samure.CursorShapeWResize
	case StateChooseRegion:
		if isRegionSet(a.chosenRegion) {
			return samure.CursorShapePointer
		} else {
			return samure.CursorShapeDefault
		}
	case StateChooseOutput:
		return samure.CursorShapePointer
	default:
		return samure.CursorShapeDefault
	}
}
