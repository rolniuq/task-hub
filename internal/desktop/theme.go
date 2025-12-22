package desktop

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type CustomTheme struct{}

func (t *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary:
		return color.RGBA{R: 99, G: 102, B: 241, A: 255} // Modern indigo
	case theme.ColorNameBackground:
		if variant == theme.VariantDark {
			return color.RGBA{R: 17, G: 24, B: 39, A: 255} // Dark slate
		}
		return color.RGBA{R: 248, G: 250, B: 252, A: 255} // Light slate
	case theme.ColorNameForeground:
		if variant == theme.VariantDark {
			return color.RGBA{R: 248, G: 250, B: 252, A: 255} // Light text
		}
		return color.RGBA{R: 15, G: 23, B: 42, A: 255} // Dark text
	case theme.ColorNameHover:
		return color.RGBA{R: 99, G: 102, B: 241, A: 50} // Indigo with transparency
	case theme.ColorNameFocus:
		return color.RGBA{R: 99, G: 102, B: 241, A: 100} // Indigo focus
	case theme.ColorNameButton:
		return color.RGBA{R: 99, G: 102, B: 241, A: 255} // Indigo button
	case theme.ColorNameDisabled:
		return color.RGBA{R: 148, G: 163, B: 184, A: 255} // Gray disabled
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 255, G: 255, B: 255, A: 255} // White input
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (t *CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 16
	case theme.SizeNameHeadingText:
		return 24
	case theme.SizeNameSubHeadingText:
		return 20
	case theme.SizeNameCaptionText:
		return 12
	case theme.SizeNamePadding:
		return 12
	case theme.SizeNameInnerPadding:
		return 8
	case theme.SizeNameScrollBar:
		return 12
	case theme.SizeNameScrollBarSmall:
		return 8
	case theme.SizeNameSeparatorThickness:
		return 1
	case theme.SizeNameInputBorder:
		return 2
	default:
		return theme.DefaultTheme().Size(name)
	}
}
