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
	StateNone       = iota
	StateDragNormal = iota
)

type App struct {
	start   [2]float64 // The top left corner of the selection box
	end     [2]float64 // The bottom right corner of the selection box
	pointer [2]float64 // The raw position of the pointer in global coordinates
	anchor  [2]float64 // The position where the pointer has been released

	state       int
	clearScreen bool

	backgroundColor [4]float64
	selectionColor  [4]float64
	borderColor     [4]float64
	textColor       [4]float64
	padding         float64
	aspect          float64
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
				ctx.SetRunning(false)
			}
		}
	case samure.EventPointerMotion:
		a.pointer[0] = e.X + float64(e.Seat.PointerFocus().Output().Geo().X)
		a.pointer[1] = e.Y + float64(e.Seat.PointerFocus().Output().Geo().Y)

		switch a.state {
		case StateDragNormal:
			a.computeStartEnd()
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
			}
		}
	}
}

func (a *App) OnUpdate(ctx samure.Context, deltaTime float64) {

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
