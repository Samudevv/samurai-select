package main

import (
	"fmt"
	"os"

	samure "github.com/PucklaJ/samurai-render-go"
	"github.com/PucklaJ/samurai-render-go/backends/cairo"
)

func main() {
	fmt.Println("Welcome to Samurai Select!")

	a, err := CreateApp(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create samurai-select app: %v\n", err)
		os.Exit(1)
	}

	b := &cairo.Backend{}

	cfg := samure.CreateContextConfig(a)

	ctx, err := samure.CreateContextWithBackend(cfg, b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create samurai-render context: %v\n", err)
		os.Exit(1)
	}
	defer ctx.Destroy()

	ctx.Run()

	fmt.Println("Good Bye! Please come again.")
}
