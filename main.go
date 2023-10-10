package main

import (
	"fmt"
	"os"

	samure "github.com/PucklaJ/samurai-render-go"
	"github.com/PucklaJ/samurai-render-go/backends/cairo"
)

func main() {
	a, err := CreateApp(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create samurai-select app: %v\n", err)
		os.Exit(1)
	}

	b := &cairo.Backend{}

	cfg := samure.CreateContextConfig(a)
	cfg.PointerInteraction = true
	cfg.KeyboardInteraction = true

	ctx, err := samure.CreateContextWithBackend(cfg, b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create samurai-render context: %v\n", err)
		os.Exit(1)
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
		os.Exit(1)
	}

	fmt.Printf("%d,%d %dx%d\n", sel.X, sel.Y, sel.W, sel.H)
}
