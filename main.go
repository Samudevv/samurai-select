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
	"os"
	"os/exec"
	"strings"
	"time"

	samure "github.com/PucklaJ/samurai-render-go"
	"github.com/PucklaJ/samurai-render-go/backends/cairo"
)

func run() int {
	a, err := CreateApp(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create samurai-select app: %v\n", err)
		return 1
	}

	b := &cairo.Backend{}

	cfg := samure.CreateContextConfig(a)
	cfg.PointerInteraction = true
	cfg.KeyboardInteraction = true

	ctx, err := samure.CreateContextWithBackend(cfg, b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create samurai-render context: %v\n", err)
		return 1
	}
	defer ctx.Destroy()

	if isRegionSet(a.chosenRegion) {
		ctx.SetPointerShape(samure.CursorShapePointer)
	}

	if a.state == StateChooseOutput {
		if a.regionsObj != nil {
			x, y, err := a.regionsObj.CursorPos()
			if err == nil {
				a.pointer[0] = float64(x)
				a.pointer[1] = float64(y)
				for i := 0; i < ctx.LenOutputs(); i++ {
					if ctx.Output(i).PointInOutput(int(a.pointer[0]), int(a.pointer[1])) {
						a.selectedOutput = ctx.Output(i)
					}
				}
			}
		}
		ctx.SetPointerShape(samure.CursorShapePointer)
	}

	if flags.FreezeScreen {
		for i := 0; i < ctx.LenOutputs(); i++ {
			o := ctx.Output(i)

			bg, err := samure.CreateLayerSurface(ctx, &o, samure.LayerTop, samure.AnchorFill, false, false, false)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not create surface to freeze screen for \"%s\": %v\n", o.Name(), err)
				continue
			}
			defer bg.Destroy(ctx)

			s, err := o.Screenshot(ctx, false)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Screenshot of \"%s\" failed: %v\n", o.Name(), err)
				continue
			}

			bg.DrawBuffer(s)
			s.Destroy()
		}
	}

	ctx.SetRenderState(samure.RenderStateOnce)
	ctx.Run()

	sel, err := a.GetSelection()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	outputName := "nil"
	if a.selectedOutput.Handle != nil {
		outputName = a.selectedOutput.Name()
	}

	geometry := fmt.Sprintf(flags.Format, sel.X, sel.Y, sel.W, sel.H, outputName)
	fmt.Println(geometry)

	if flags.Screenshot || flags.Command != "" {
		a.clearScreen = true
		for i := 0; i < ctx.LenOutputs(); i++ {
			ctx.RenderOutput(ctx.Output(i))
		}
		ctx.Flush()
	}

	if flags.Screenshot {
		now := time.Now()
		var screenshotFileName string
		if strings.Contains(flags.ScreenshotOutput, "%") {
			screenshotFileName = fmt.Sprintf(flags.ScreenshotOutput, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
		} else {
			screenshotFileName = flags.ScreenshotOutput
		}

		var screenshotFlags []string
		if flags.ScreenshotFlags != "" {
			screenshotFlags = append(screenshotFlags, strings.FieldsFunc(flags.ScreenshotFlags, func(c rune) bool {
				return c == ' '
			})...)
		}
		screenshotFlags = append(
			screenshotFlags,
			"-g",
			geometry,
			screenshotFileName,
		)
		grimPath, err := exec.LookPath("grim")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not find grim")
			return 1
		}

		grim := exec.Command(grimPath, screenshotFlags...)
		grim.Stderr = os.Stderr
		grim.Stdout = os.Stderr

		if err := grim.Run(); err != nil {
			return 1
		}
	}

	if flags.Command != "" {
		commandArgs := strings.FieldsFunc(flags.Command, func(c rune) bool {
			return c == ' '
		})
		for i := range commandArgs {
			commandArgs[i] = strings.ReplaceAll(commandArgs[i], "%geometry%", geometry)
		}
		cmd := exec.Command(commandArgs[0], commandArgs[1:]...)
		fmt.Println(cmd.Args)
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return 1
		}
	}

	return 0
}

func main() {
	os.Exit(run())
}
