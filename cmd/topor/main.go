package main

import (
	"fyne.io/fyne/v2/app"

	"github.com/Serge-Nook/topor/internal/gui"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(&gui.ToporTheme{})
	gui.NewToporApp(a)
}
