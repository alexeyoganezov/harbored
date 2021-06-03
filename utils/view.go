package utils

import "fyne.io/fyne"

type View struct {
	App    *fyne.App
	Window *fyne.Window
	El     *fyne.Container
	Renderable
}

type Renderable interface {
	Render() *fyne.Container
}
