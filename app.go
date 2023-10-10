package main

import (
	samure "github.com/PucklaJ/samurai-render-go"
)

type App struct {
}

func CreateApp(argv []string) (*App, error) {
	a := &App{}
	return a, nil
}

func (a *App) OnEvent(ctx samure.Context, event interface{}) {

}

func (a *App) OnRender(ctx samure.Context, layerSurface samure.LayerSurface, o samure.Rect, deltaTime float64) {

}

func (a *App) OnUpdate(ctx samure.Context, deltaTime float64) {

}
