package main

import (
	"math"

	samure "github.com/PucklaJ/samurai-render-go"
)

func (a *App) pointerDown(ctx samure.Context, px, py float64) {
	switch a.state {
	case StateNone:
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
				a.grabberAnim = 0.0
				a.anchor[0], a.anchor[1] = px, py
				a.computeStartEnd(px, py, px, py)
				a.state = StateDragNormal
				ctx.SetRenderState(samure.RenderStateOnce)
			}
		}
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

func (a *App) pointerMove(ctx samure.Context, px, py, dx, dy float64) {
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
		a.handleOverlapAndAspectRatio()
		ctx.SetRenderState(samure.RenderStateOnce)
	case StateDragTop:
		h += y - poy
		y = poy
		a.start[1] = y
		a.end[1] = y + h
		a.handleOverlapAndAspectRatio()
		ctx.SetRenderState(samure.RenderStateOnce)
	case StateDragTopRight:
		w = pox - x
		h += y - poy
		y = poy
		a.start[1] = y
		a.end[0] = x + w
		a.end[1] = y + h
		a.handleOverlapAndAspectRatio()
		ctx.SetRenderState(samure.RenderStateOnce)
	case StateDragRight:
		w = pox - x
		a.end[0] = x + w
		a.handleOverlapAndAspectRatio()
		ctx.SetRenderState(samure.RenderStateOnce)
	case StateDragBottomRight:
		w = pox - x
		h = poy - y
		a.end[0] = x + w
		a.end[1] = y + h
		a.handleOverlapAndAspectRatio()
		ctx.SetRenderState(samure.RenderStateOnce)
	case StateDragBottom:
		h = poy - y
		a.end[1] = y + h
		a.handleOverlapAndAspectRatio()
		ctx.SetRenderState(samure.RenderStateOnce)
	case StateDragBottomLeft:
		w += x - pox
		x = pox
		h = poy - y
		a.start[0] = x
		a.end[0] = x + w
		a.end[1] = y + h
		a.handleOverlapAndAspectRatio()
		ctx.SetRenderState(samure.RenderStateOnce)
	case StateDragLeft:
		w += x - pox
		x = pox
		a.start[0] = x
		a.end[0] = x + w
		a.handleOverlapAndAspectRatio()
		ctx.SetRenderState(samure.RenderStateOnce)
	case StateDragMiddle:
		x += dx
		y += dy
		a.start[0] = x
		a.start[1] = y
		a.end[0] = x + w
		a.end[1] = y + h
		ctx.SetRenderState(samure.RenderStateOnce)
	}
}

func (a *App) OnEvent(ctx samure.Context, event interface{}) {
	switch e := event.(type) {
	case samure.EventPointerButton:
		if e.Button == samure.ButtonLeft {
			switch e.State {
			case samure.StatePressed:
				a.pointerDown(ctx, a.pointer[0], a.pointer[1])
			case samure.StateReleased:
				a.pointerUp(ctx)
			}
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
		a.pointerDown(ctx, a.pointer[0], a.pointer[1])
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
		a.pointerMove(ctx, px, py, dx, dy)
	case samure.EventTouchMotion:
		if a.touchID == nil || *a.touchID != e.TouchID {
			break
		}

		px := e.X + float64(a.touchFocus.Geo().X)
		py := e.Y + float64(a.touchFocus.Geo().Y)
		dx := px - a.pointer[0]
		dy := py - a.pointer[1]
		a.pointer[0], a.pointer[1] = px, py
		a.pointerMove(ctx, px, py, dx, dy)
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
