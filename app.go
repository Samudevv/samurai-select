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
	StateDragMiddle      = iota
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

func (a *App) OnUpdate(ctx samure.Context, deltaTime float64) {
	if a.state < StateAlter || a.state > StateDragLeft {
		return
	}

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
}
