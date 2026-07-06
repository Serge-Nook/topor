package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// ToporTheme implements fyne.Theme with dark blue + cyan colors.
type ToporTheme struct{}

var _ fyne.Theme = (*ToporTheme)(nil)

func (t *ToporTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 10, G: 25, B: 41, A: 255} // #0a1929
	case theme.ColorNameButton:
		return color.NRGBA{R: 0, G: 188, B: 212, A: 255} // #00bcd4
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 26, G: 48, B: 80, A: 255} // #1a3050
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 58, G: 96, B: 112, A: 255} // #3a6070
	case theme.ColorNameError:
		return color.NRGBA{R: 239, G: 83, B: 80, A: 255} // #ef5350
	case theme.ColorNameForeground:
		return color.NRGBA{R: 128, G: 229, B: 208, A: 255} // #80e5d0
	case theme.ColorNameHover:
		return color.NRGBA{R: 0, G: 172, B: 193, A: 255} // #00acc1
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 15, G: 38, B: 64, A: 255} // #0f2640
	case theme.ColorNameInputBorder:
		return color.NRGBA{R: 26, G: 58, B: 92, A: 255} // #1a3a5c
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 92, G: 184, B: 165, A: 255} // #5cb8a5
	case theme.ColorNamePressed:
		return color.NRGBA{R: 0, G: 151, B: 167, A: 255} // #0097a7
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0, G: 188, B: 212, A: 255} // #00bcd4
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 30, G: 73, B: 118, A: 255} // #1e4976
	case theme.ColorNameSelection:
		return color.NRGBA{R: 0, G: 131, B: 143, A: 255} // #00838f
	case theme.ColorNameSeparator:
		return color.NRGBA{R: 26, G: 58, B: 92, A: 255} // #1a3a5c
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 100}
	case theme.ColorNameHeaderBackground:
		return color.NRGBA{R: 18, G: 42, B: 66, A: 255} // #122a42
	case theme.ColorNameMenuBackground:
		return color.NRGBA{R: 13, G: 33, B: 55, A: 255} // #0d2137
	case theme.ColorNameOverlayBackground:
		return color.NRGBA{R: 13, G: 33, B: 55, A: 255} // #0d2137
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (t *ToporTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *ToporTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *ToporTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 13
	case theme.SizeNamePadding:
		return 6
	case theme.SizeNameInnerPadding:
		return 4
	}
	return theme.DefaultTheme().Size(name)
}
