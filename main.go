package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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

	if flags.FreezeScreen {
		for i := 0; i < ctx.LenOutputs(); i++ {
			o := ctx.Output(i)

			bg, err := samure.CreateLayerSurface(ctx, &o, samure.LayerTop, samure.AnchorFill, false, false, false)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not create surface to freeze screen for \"%s\": %v\n", o.Name(), err)
				continue
			}
			defer bg.Destroy(ctx)

			s, err := o.Screenshot(ctx)
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

	geometry := fmt.Sprintf("%d,%d %dx%d", sel.X, sel.Y, sel.W, sel.H)
	fmt.Println(geometry)

	if flags.Screenshot || flags.Command != "" {
		a.clearScreen = true
		for i := 0; i < ctx.LenOutputs(); i++ {
			ctx.RenderOutput(ctx.Output(i), 0.0)
		}
		ctx.Flush()
	}

	if flags.Screenshot {
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
			flags.ScreenshotOutput,
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
