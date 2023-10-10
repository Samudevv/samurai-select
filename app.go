package main

import (
	"errors"

	samure "github.com/PucklaJ/samurai-render-go"
	samure_cairo "github.com/PucklaJ/samurai-render-go/backends/cairo"
	"github.com/gotk3/gotk3/cairo"
)

const (
	StateNone       = iota
	StateDragNormal = iota
)

type App struct {
	start [2]float64
	end   [2]float64
	hold  [2]float64

	state int

	backgroundColor [4]float64
	selectionColor  [4]float64
	borderColor     [4]float64
	borderWidth     float64
}

func CreateApp(argv []string) (*App, error) {
	a := &App{}
	a.backgroundColor = [4]float64{0.1, 0.1, 0.1, 0.5}
	a.selectionColor = [4]float64{0.0, 0.0, 0.0, 0.0}
	a.borderColor = [4]float64{0.0, 0.0, 0.0, 1.0}
	a.borderWidth = 2.0
	return a, nil
}

func (a App) GetSelection() (samure.Rect, error) {
	if a.start[0] == 0.0 && a.start[1] == 0.0 && a.end[0] == 0.0 && a.end[1] == 0.0 {
		return samure.Rect{}, errors.New("selection cancelled")
	}

	start := a.start
	end := a.end
	if start[0] > end[0] {
		start[0], end[0] = end[0], start[0]
	}
	if start[1] > end[1] {
		start[1], end[1] = end[1], start[1]
	}

	return samure.Rect{
		X: int(start[0]),
		Y: int(start[1]),
		W: int(end[0] - start[0]),
		H: int(end[1] - start[1]),
	}, nil
}

func (a *App) OnEvent(ctx samure.Context, event interface{}) {
	switch e := event.(type) {
	case samure.EventPointerButton:
		switch a.state {
		case StateNone:
			if e.Button == samure.ButtonLeft && e.State == samure.StatePressed {
				a.start = a.hold
				a.end = a.hold
				a.state = StateDragNormal
				ctx.SetRenderState(samure.RenderStateOnce)
			}
		case StateDragNormal:
			if e.Button == samure.ButtonLeft && e.State == samure.StateReleased {
				ctx.SetRunning(false)
			}
		}
	case samure.EventPointerMotion:
		switch a.state {
		case StateNone:
			a.hold[0] = e.X + float64(e.Seat.PointerFocus().Output().Geo().X)
			a.hold[1] = e.Y + float64(e.Seat.PointerFocus().Output().Geo().Y)
		case StateDragNormal:
			a.end[0] = e.X + float64(e.Seat.PointerFocus().Output().Geo().X)
			a.end[1] = e.Y + float64(e.Seat.PointerFocus().Output().Geo().Y)
			ctx.SetRenderState(samure.RenderStateOnce)
		}
	case samure.EventPointerEnter:
		e.Seat.SetPointerShape(samure.CursorShapeCrosshair)
	case samure.EventKeyboardKey:
		if e.State == samure.StateReleased {
			switch e.Key {
			case samure.KeyEsc:
				a.start = [2]float64{0.0, 0.0}
				a.end = a.start
				ctx.SetRunning(false)
			case samure.KeyEnter:
				ctx.SetRunning(false)
			}
		}
	}
}

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

	start := [2]float64{
		a.start[0],
		a.start[1],
	}
	end := [2]float64{
		a.end[0],
		a.end[1],
	}
	if start[0] > end[0] {
		start[0], end[0] = end[0], start[0]
	}
	if start[1] > end[1] {
		start[1], end[1] = end[1], start[1]
	}

	if o.RectInOutput(int(start[0]), int(start[1]), int(end[0]-start[0]), int(end[1]-start[1])) {
		start[0] = o.RelX(start[0])
		start[1] = o.RelY(start[1])
		end[0] = o.RelX(end[0])
		end[1] = o.RelY(end[1])

		// Render the selection
		c.SetSourceRGBA(
			a.selectionColor[0],
			a.selectionColor[1],
			a.selectionColor[2],
			a.selectionColor[3],
		)
		c.Rectangle(start[0], start[1], end[0]-start[0], end[1]-start[1])
		c.Fill()
		// Render the border of the selection
		c.SetSourceRGBA(
			a.borderColor[0],
			a.borderColor[1],
			a.borderColor[2],
			a.borderColor[3],
		)
		c.Rectangle(start[0], start[1], end[0]-start[0], end[1]-start[1])
		c.SetLineWidth(a.borderWidth)
		c.Stroke()
	}
}

func (a *App) OnUpdate(ctx samure.Context, deltaTime float64) {

}
